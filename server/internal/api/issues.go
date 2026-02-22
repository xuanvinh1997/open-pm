package api

import (
	"encoding/json"
	"fmt"
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
	EstimatePoint      *int            `json:"estimate_point,omitempty"`
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
	EstimatePoint      *int            `json:"estimate_point,omitempty"`
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
	if startDate != nil && targetDate != nil && targetDate.Before(*startDate) {
		return badRequestError("target_date must not be before start_date")
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
		EstimatePoint:   req.EstimatePoint,
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

	// Auto-subscribe the creator and assignees
	_ = a.queries.AddIssueSubscriber(ctx, issue.ID, userID)
	for _, assigneeID := range req.AssigneeIDs {
		_ = a.queries.AddIssueSubscriber(ctx, issue.ID, assigneeID)
	}

	// Notify assigned users
	for _, assigneeID := range req.AssigneeIDs {
		a.notifyUsers(ctx, []uuid.UUID{assigneeID}, userID, notifyParams{
			WorkspaceID: workspaceID,
			ProjectID:   &projectID,
			Title:       fmt.Sprintf("You were assigned to \"%s\"", issue.Name),
			Message:     "You have been assigned to a new issue",
			EntityType:  strPtr("issue"),
			EntityID:    &issue.ID,
		})
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
	subIssues, err := a.queries.ListSubIssues(r.Context(), issueID)
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
	if updateStartDate != nil && updateTargetDate != nil && updateTargetDate.Before(*updateStartDate) {
		return badRequestError("target_date must not be before start_date")
	}

	// Fetch old issue for activity logging
	oldIssue, err := a.queries.GetIssueByID(ctx, issueID)
	if err != nil {
		return notFoundError("issue not found")
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
		EstimatePoint:   req.EstimatePoint,
		UpdatedBy:       &userID,
	})
	if err != nil {
		return internalServerError("failed to update issue")
	}

	// Sync assignees if provided
	if req.AssigneeIDs != nil {
		currentAssignees, _ := a.queries.ListIssueAssignees(ctx, issueID)
		currentSet := make(map[uuid.UUID]bool, len(currentAssignees))
		for _, assignee := range currentAssignees {
			currentSet[assignee.ID] = true
		}
		newSet := make(map[uuid.UUID]bool, len(req.AssigneeIDs))
		for _, id := range req.AssigneeIDs {
			newSet[id] = true
		}

		// Remove assignees no longer in the list
		for id := range currentSet {
			if !newSet[id] {
				if err := a.queries.RemoveIssueAssignee(ctx, issueID, id); err != nil {
					log.Warn().Err(err).Str("issue_id", issueID.String()).Str("assignee_id", id.String()).Msg("failed to remove assignee")
				}
			}
		}

		// Add new assignees
		for id := range newSet {
			if !currentSet[id] {
				if err := a.queries.AddIssueAssignee(ctx, issueID, id); err != nil {
					log.Warn().Err(err).Str("issue_id", issueID.String()).Str("assignee_id", id.String()).Msg("failed to add assignee")
				}
				// Auto-subscribe new assignees
				_ = a.queries.AddIssueSubscriber(ctx, issueID, id)

				// Notify newly assigned user
				a.notifyUsers(ctx, []uuid.UUID{id}, userID, notifyParams{
					WorkspaceID: issue.WorkspaceID,
					ProjectID:   &issue.ProjectID,
					Title:       fmt.Sprintf("You were assigned to \"%s\"", issue.Name),
					Message:     "You have been assigned to an issue",
					EntityType:  strPtr("issue"),
					EntityID:    &issueID,
				})
			}
		}
	}

	// Sync labels if provided
	if req.LabelIDs != nil {
		currentLabels, _ := a.queries.ListIssueLabels(ctx, issueID)
		currentSet := make(map[uuid.UUID]bool, len(currentLabels))
		for _, label := range currentLabels {
			currentSet[label.ID] = true
		}
		newSet := make(map[uuid.UUID]bool, len(req.LabelIDs))
		for _, id := range req.LabelIDs {
			newSet[id] = true
		}

		for id := range currentSet {
			if !newSet[id] {
				if err := a.queries.RemoveIssueLabel(ctx, issueID, id); err != nil {
					log.Warn().Err(err).Str("issue_id", issueID.String()).Str("label_id", id.String()).Msg("failed to remove label")
				}
			}
		}
		for id := range newSet {
			if !currentSet[id] {
				if err := a.queries.AddIssueLabel(ctx, issueID, id); err != nil {
					log.Warn().Err(err).Str("issue_id", issueID.String()).Str("label_id", id.String()).Msg("failed to add label")
				}
			}
		}
	}

	// Log activity for field changes
	logFieldChange := func(field, oldVal, newVal string) {
		if oldVal != newVal {
			if _, err := a.queries.CreateIssueActivity(ctx, CreateActivityParams{
				IssueID:     &issueID,
				ProjectID:   issue.ProjectID,
				WorkspaceID: issue.WorkspaceID,
				Verb:        "updated",
				Field:       strPtr(field),
				OldValue:    strPtr(oldVal),
				NewValue:    strPtr(newVal),
				Comment:     fmt.Sprintf("changed %s", field),
				ActorID:     &userID,
			}); err != nil {
				log.Warn().Err(err).Str("field", field).Msg("failed to log field change activity")
			}
		}
	}

	if req.Priority != nil {
		logFieldChange("priority", oldIssue.Priority, issue.Priority)
	}
	if req.IssueType != nil {
		logFieldChange("issue_type", oldIssue.IssueType, issue.IssueType)
	}
	if req.Name != nil {
		logFieldChange("name", oldIssue.Name, issue.Name)
	}
	if req.StateID != nil {
		oldStateStr := ""
		if oldIssue.StateID != nil {
			oldStateStr = oldIssue.StateID.String()
		}
		newStateStr := ""
		if issue.StateID != nil {
			newStateStr = issue.StateID.String()
		}
		logFieldChange("state", oldStateStr, newStateStr)
	}

	// Notify assignees and subscribers on state change
	if req.StateID != nil {
		assignees, _ := a.queries.ListIssueAssignees(ctx, issueID)
		subscribers, _ := a.queries.ListIssueSubscribers(ctx, issueID)
		receiverIDs := collectUniqueUserIDs(assignees, subscribers)

		a.notifyUsers(ctx, receiverIDs, userID, notifyParams{
			WorkspaceID: issue.WorkspaceID,
			ProjectID:   &issue.ProjectID,
			Title:       fmt.Sprintf("State changed on \"%s\"", issue.Name),
			Message:     "Issue state was updated",
			EntityType:  strPtr("issue"),
			EntityID:    &issueID,
		})
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

	// Notify assignees and subscribers about the new comment
	projectID := getProjectID(ctx)
	assignees, _ := a.queries.ListIssueAssignees(ctx, issueID)
	subscribers, _ := a.queries.ListIssueSubscribers(ctx, issueID)
	receiverIDs := collectUniqueUserIDs(assignees, subscribers)

	issue, _ := a.queries.GetIssueByID(ctx, issueID)
	title := "New comment on an issue you follow"
	if issue != nil {
		title = fmt.Sprintf("New comment on \"%s\"", issue.Name)
	}
	a.notifyUsers(ctx, receiverIDs, userID, notifyParams{
		WorkspaceID: workspaceID,
		ProjectID:   &projectID,
		Title:       title,
		Message:     truncateString(body.CommentStripped, 100),
		EntityType:  strPtr("issue_comment"),
		EntityID:    &issueID,
	})

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
