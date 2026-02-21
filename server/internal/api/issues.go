package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

var validPriorities = map[string]bool{
	"urgent": true, "high": true, "medium": true, "low": true, "none": true,
}

var validIssueTypes = map[string]bool{
	"story": true, "bug": true, "task": true, "epic": true,
}

type CreateIssueRequest struct {
	Name               string          `json:"name" validate:"required,min=1"`
	DescriptionHTML    string          `json:"description_html,omitempty"`
	DescriptionJSON    json.RawMessage `json:"description_json,omitempty"`
	Priority           string          `json:"priority,omitempty"`
	IssueType          string          `json:"issue_type,omitempty"`
	StateID            *uuid.UUID      `json:"state_id,omitempty"`
	ParentID           *uuid.UUID      `json:"parent_id,omitempty"`
	AssigneeIDs        []uuid.UUID     `json:"assignee_ids,omitempty"`
	LabelIDs           []uuid.UUID     `json:"label_ids,omitempty"`
	StartDate          *string         `json:"start_date,omitempty"`
	TargetDate         *string         `json:"target_date,omitempty"`
	IsDraft            bool            `json:"is_draft,omitempty"`
}

type UpdateIssueRequest struct {
	Name               *string         `json:"name,omitempty"`
	DescriptionHTML    *string         `json:"description_html,omitempty"`
	DescriptionJSON    json.RawMessage `json:"description_json,omitempty"`
	Priority           *string         `json:"priority,omitempty"`
	IssueType          *string         `json:"issue_type,omitempty"`
	StateID            *uuid.UUID      `json:"state_id,omitempty"`
	ParentID           *uuid.UUID      `json:"parent_id,omitempty"`
	SortOrder          *float64        `json:"sort_order,omitempty"`
	AssigneeIDs        []uuid.UUID     `json:"assignee_ids,omitempty"`
	LabelIDs           []uuid.UUID     `json:"label_ids,omitempty"`
	StartDate          *string         `json:"start_date,omitempty"`
	TargetDate         *string         `json:"target_date,omitempty"`
}

// ListIssues handles GET .../projects/{projectID}/issues
func (a *API) ListIssues(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	projectID := getProjectID(ctx)
	pg := parsePagination(r)

	issues, err := a.queries.ListIssuesByProject(ctx, projectID, pg.PerPage, pg.Offset())
	if err != nil {
		return internalServerError("failed to list issues")
	}

	count, err := a.queries.CountIssuesByProject(ctx, projectID)
	if err != nil {
		log.Warn().Err(err).Str("project_id", projectID.String()).Msg("failed to count issues")
	}

	return sendJSON(w, http.StatusOK, PaginatedResponse{
		Results:    issues,
		TotalCount: count,
	})
}

// CreateIssue handles POST .../projects/{projectID}/issues
func (a *API) CreateIssue(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)
	projectID := getProjectID(ctx)
	workspaceID := getWorkspaceID(ctx)

	var req CreateIssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	priority := req.Priority
	if priority == "" {
		priority = "none"
	}
	if !validPriorities[priority] {
		return badRequestError("invalid priority: %s", priority)
	}
	issueType := req.IssueType
	if issueType == "" {
		issueType = "task"
	}
	if !validIssueTypes[issueType] {
		return badRequestError("invalid issue type: %s", issueType)
	}

	var startDate, targetDate *time.Time
	if req.StartDate != nil {
		t, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return badRequestError("invalid start_date format, expected YYYY-MM-DD")
		}
		startDate = &t
	}
	if req.TargetDate != nil {
		t, err := time.Parse("2006-01-02", *req.TargetDate)
		if err != nil {
			return badRequestError("invalid target_date format, expected YYYY-MM-DD")
		}
		targetDate = &t
	}

	issue, err := a.queries.CreateIssue(ctx, CreateIssueParams{
		ProjectID:       projectID,
		WorkspaceID:     workspaceID,
		ParentID:        req.ParentID,
		StateID:         req.StateID,
		Name:            req.Name,
		DescriptionHTML: req.DescriptionHTML,
		DescriptionJSON: req.DescriptionJSON,
		Priority:        priority,
		IssueType:       issueType,
		StartDate:       startDate,
		TargetDate:      targetDate,
		IsDraft:         req.IsDraft,
		CreatedBy:       &userID,
	})
	if err != nil {
		return internalServerError("failed to create issue")
	}

	// Add assignees
	for _, assigneeID := range req.AssigneeIDs {
		if err := a.queries.AddIssueAssignee(ctx, issue.ID, assigneeID); err != nil {
			log.Warn().Err(err).Str("issue_id", issue.ID.String()).Str("assignee_id", assigneeID.String()).Msg("failed to add assignee")
		}
	}

	// Add labels
	for _, labelID := range req.LabelIDs {
		if err := a.queries.AddIssueLabel(ctx, issue.ID, labelID); err != nil {
			log.Warn().Err(err).Str("issue_id", issue.ID.String()).Str("label_id", labelID.String()).Msg("failed to add label")
		}
	}

	// Log activity
	if _, err := a.queries.CreateIssueActivity(ctx, CreateActivityParams{
		IssueID:     &issue.ID,
		ProjectID:   projectID,
		WorkspaceID: workspaceID,
		Verb:        "created",
		Comment:     "created the issue",
		ActorID:     &userID,
	}); err != nil {
		log.Warn().Err(err).Str("issue_id", issue.ID.String()).Msg("failed to log issue creation activity")
	}

	return sendJSON(w, http.StatusCreated, issue)
}

// GetIssue handles GET .../issues/{issueID}
func (a *API) GetIssue(w http.ResponseWriter, r *http.Request) error {
	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}

	issue, err := a.queries.GetIssueByID(r.Context(), issueID)
	if err != nil {
		return notFoundError("issue not found")
	}

	// Enrich with assignees and labels
	assignees, err := a.queries.ListIssueAssignees(r.Context(), issueID)
	if err != nil {
		log.Warn().Err(err).Str("issue_id", issueID.String()).Msg("failed to list issue assignees")
	}
	labels, err := a.queries.ListIssueLabels(r.Context(), issueID)
	if err != nil {
		log.Warn().Err(err).Str("issue_id", issueID.String()).Msg("failed to list issue labels")
	}
	subIssues, err := a.queries.ListIssuesByProject(r.Context(), issueID, 100, 0) // TODO: use ListSubIssues
	if err != nil {
		log.Warn().Err(err).Str("issue_id", issueID.String()).Msg("failed to list sub-issues")
	}

	type IssueDetail struct {
		*Issue
		Assignees []*UserSummary `json:"assignees"`
		Labels    []*Label       `json:"labels"`
		SubIssues []*Issue       `json:"sub_issues"`
	}

	return sendJSON(w, http.StatusOK, IssueDetail{
		Issue:     issue,
		Assignees: assignees,
		Labels:    labels,
		SubIssues: subIssues,
	})
}

// UpdateIssue handles PUT .../issues/{issueID}
func (a *API) UpdateIssue(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)
	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}

	var req UpdateIssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}

	if req.Priority != nil && *req.Priority != "" && !validPriorities[*req.Priority] {
		return badRequestError("invalid priority: %s", *req.Priority)
	}
	if req.IssueType != nil && *req.IssueType != "" && !validIssueTypes[*req.IssueType] {
		return badRequestError("invalid issue type: %s", *req.IssueType)
	}

	var updateStartDate, updateTargetDate *time.Time
	if req.StartDate != nil {
		if *req.StartDate == "" {
			updateStartDate = nil
		} else if t, err := time.Parse("2006-01-02", *req.StartDate); err == nil {
			updateStartDate = &t
		}
	}
	if req.TargetDate != nil {
		if *req.TargetDate == "" {
			updateTargetDate = nil
		} else if t, err := time.Parse("2006-01-02", *req.TargetDate); err == nil {
			updateTargetDate = &t
		}
	}

	issue, err := a.queries.UpdateIssue(ctx, issueID, UpdateIssueParams{
		StateID:         req.StateID,
		Name:            req.Name,
		DescriptionHTML: req.DescriptionHTML,
		DescriptionJSON: req.DescriptionJSON,
		Priority:        req.Priority,
		IssueType:       req.IssueType,
		StartDate:       updateStartDate,
		TargetDate:      updateTargetDate,
		ParentID:        req.ParentID,
		SortOrder:       req.SortOrder,
		UpdatedBy:       &userID,
	})
	if err != nil {
		return internalServerError("failed to update issue")
	}

	return sendJSON(w, http.StatusOK, issue)
}

// DeleteIssue handles DELETE .../issues/{issueID}
func (a *API) DeleteIssue(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)
	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}

	// Allow admin or issue creator
	if getProjectRole(ctx) < RoleAdmin {
		issue, err := a.queries.GetIssueByID(ctx, issueID)
		if err != nil {
			return notFoundError("issue not found")
		}
		if issue.CreatedBy == nil || *issue.CreatedBy != userID {
			return forbiddenError("only admins or the issue creator can delete issues")
		}
	}

	if err := a.queries.DeleteIssue(ctx, issueID); err != nil {
		return internalServerError("failed to delete issue")
	}
	return sendEmpty(w, http.StatusNoContent)
}

// ListIssueComments handles GET .../issues/{issueID}/comments
func (a *API) ListIssueComments(w http.ResponseWriter, r *http.Request) error {
	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}

	comments, err := a.queries.ListIssueComments(r.Context(), issueID)
	if err != nil {
		return internalServerError("failed to list comments")
	}
	return sendJSON(w, http.StatusOK, map[string]interface{}{"results": comments})
}

// CreateIssueComment handles POST .../issues/{issueID}/comments
func (a *API) CreateIssueComment(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)
	workspaceID := getWorkspaceID(ctx)
	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}

	var body struct {
		CommentHTML     string          `json:"comment_html" validate:"required"`
		CommentJSON     json.RawMessage `json:"comment_json,omitempty"`
		CommentStripped string          `json:"comment_stripped,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return badRequestError("invalid request body")
	}

	comment, err := a.queries.CreateIssueComment(ctx, issueID, workspaceID, body.CommentHTML, body.CommentJSON, body.CommentStripped, &userID)
	if err != nil {
		return internalServerError("failed to create comment")
	}
	return sendJSON(w, http.StatusCreated, comment)
}

// UpdateIssueComment handles PUT .../comments/{commentID}
func (a *API) UpdateIssueComment(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)
	commentID, err := uuid.FromString(chi.URLParam(r, "commentID"))
	if err != nil {
		return badRequestError("invalid comment ID")
	}

	// Only comment author or admin can update
	if getProjectRole(ctx) < RoleAdmin {
		existing, err := a.queries.GetIssueCommentByID(ctx, commentID)
		if err != nil {
			return notFoundError("comment not found")
		}
		if existing.ActorID == nil || *existing.ActorID != userID {
			return forbiddenError("only the comment author or an admin can edit this comment")
		}
	}

	var body struct {
		CommentHTML     string          `json:"comment_html" validate:"required"`
		CommentJSON     json.RawMessage `json:"comment_json,omitempty"`
		CommentStripped string          `json:"comment_stripped,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return badRequestError("invalid request body")
	}

	comment, err := a.queries.UpdateIssueComment(ctx, commentID, body.CommentHTML, body.CommentJSON, body.CommentStripped)
	if err != nil {
		return internalServerError("failed to update comment")
	}
	return sendJSON(w, http.StatusOK, comment)
}

// DeleteIssueComment handles DELETE .../comments/{commentID}
func (a *API) DeleteIssueComment(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)
	commentID, err := uuid.FromString(chi.URLParam(r, "commentID"))
	if err != nil {
		return badRequestError("invalid comment ID")
	}

	// Only comment author or admin can delete
	if getProjectRole(ctx) < RoleAdmin {
		existing, err := a.queries.GetIssueCommentByID(ctx, commentID)
		if err != nil {
			return notFoundError("comment not found")
		}
		if existing.ActorID == nil || *existing.ActorID != userID {
			return forbiddenError("only the comment author or an admin can delete this comment")
		}
	}

	if err := a.queries.DeleteIssueComment(ctx, commentID); err != nil {
		return internalServerError("failed to delete comment")
	}
	return sendEmpty(w, http.StatusNoContent)
}

// ListIssueActivities handles GET .../issues/{issueID}/activities
func (a *API) ListIssueActivities(w http.ResponseWriter, r *http.Request) error {
	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}

	activities, err := a.queries.ListIssueActivities(r.Context(), issueID)
	if err != nil {
		return internalServerError("failed to list activities")
	}
	return sendJSON(w, http.StatusOK, map[string]interface{}{"results": activities})
}
