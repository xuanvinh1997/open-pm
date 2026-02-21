package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

type CreateIssueRequest struct {
	Name               string          `json:"name" validate:"required,min=1"`
	DescriptionHTML    string          `json:"description_html,omitempty"`
	DescriptionJSON    json.RawMessage `json:"description_json,omitempty"`
	Priority           string          `json:"priority,omitempty"`
	StateID            *uuid.UUID      `json:"state_id,omitempty"`
	ParentID           *uuid.UUID      `json:"parent_id,omitempty"`
	AssigneeIDs        []uuid.UUID     `json:"assignee_ids,omitempty"`
	LabelIDs           []uuid.UUID     `json:"label_ids,omitempty"`
	IsDraft            bool            `json:"is_draft,omitempty"`
}

type UpdateIssueRequest struct {
	Name               *string         `json:"name,omitempty"`
	DescriptionHTML    *string         `json:"description_html,omitempty"`
	DescriptionJSON    json.RawMessage `json:"description_json,omitempty"`
	Priority           *string         `json:"priority,omitempty"`
	StateID            *uuid.UUID      `json:"state_id,omitempty"`
	ParentID           *uuid.UUID      `json:"parent_id,omitempty"`
	SortOrder          *float64        `json:"sort_order,omitempty"`
	AssigneeIDs        []uuid.UUID     `json:"assignee_ids,omitempty"`
	LabelIDs           []uuid.UUID     `json:"label_ids,omitempty"`
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

	count, _ := a.queries.CountIssuesByProject(ctx, projectID)

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

	issue, err := a.queries.CreateIssue(ctx, CreateIssueParams{
		ProjectID:       projectID,
		WorkspaceID:     workspaceID,
		ParentID:        req.ParentID,
		StateID:         req.StateID,
		Name:            req.Name,
		DescriptionHTML: req.DescriptionHTML,
		DescriptionJSON: req.DescriptionJSON,
		Priority:        priority,
		IsDraft:         req.IsDraft,
		CreatedBy:       &userID,
	})
	if err != nil {
		return internalServerError("failed to create issue")
	}

	// Add assignees
	for _, assigneeID := range req.AssigneeIDs {
		_ = a.queries.AddIssueAssignee(ctx, issue.ID, assigneeID)
	}

	// Add labels
	for _, labelID := range req.LabelIDs {
		_ = a.queries.AddIssueLabel(ctx, issue.ID, labelID)
	}

	// Log activity
	_, _ = a.queries.CreateIssueActivity(ctx, CreateActivityParams{
		IssueID:     &issue.ID,
		ProjectID:   projectID,
		WorkspaceID: workspaceID,
		Verb:        "created",
		Comment:     "created the issue",
		ActorID:     &userID,
	})

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
	assignees, _ := a.queries.ListIssueAssignees(r.Context(), issueID)
	labels, _ := a.queries.ListIssueLabels(r.Context(), issueID)
	subIssues, _ := a.queries.ListIssuesByProject(r.Context(), issueID, 100, 0) // TODO: use ListSubIssues

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

	issue, err := a.queries.UpdateIssue(ctx, issueID, UpdateIssueParams{
		StateID:         req.StateID,
		Name:            req.Name,
		DescriptionHTML: req.DescriptionHTML,
		DescriptionJSON: req.DescriptionJSON,
		Priority:        req.Priority,
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
	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}

	if err := a.queries.DeleteIssue(r.Context(), issueID); err != nil {
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
	commentID, err := uuid.FromString(chi.URLParam(r, "commentID"))
	if err != nil {
		return badRequestError("invalid comment ID")
	}

	var body struct {
		CommentHTML     string          `json:"comment_html" validate:"required"`
		CommentJSON     json.RawMessage `json:"comment_json,omitempty"`
		CommentStripped string          `json:"comment_stripped,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return badRequestError("invalid request body")
	}

	comment, err := a.queries.UpdateIssueComment(r.Context(), commentID, body.CommentHTML, body.CommentJSON, body.CommentStripped)
	if err != nil {
		return internalServerError("failed to update comment")
	}
	return sendJSON(w, http.StatusOK, comment)
}

// DeleteIssueComment handles DELETE .../comments/{commentID}
func (a *API) DeleteIssueComment(w http.ResponseWriter, r *http.Request) error {
	commentID, err := uuid.FromString(chi.URLParam(r, "commentID"))
	if err != nil {
		return badRequestError("invalid comment ID")
	}

	if err := a.queries.DeleteIssueComment(r.Context(), commentID); err != nil {
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
