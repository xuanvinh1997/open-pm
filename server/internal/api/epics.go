package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

var validEpicStatuses = map[string]bool{
	"backlog":     true,
	"planned":     true,
	"in-progress": true,
	"paused":      true,
	"completed":   true,
	"cancelled":   true,
}

type CreateEpicRequest struct {
	Name        string     `json:"name" validate:"required,min=1,max=255"`
	Description string     `json:"description,omitempty"`
	StartDate   string     `json:"start_date,omitempty"`
	TargetDate  string     `json:"target_date,omitempty"`
	Status      string     `json:"status,omitempty"`
	LeadID      *uuid.UUID `json:"lead_id,omitempty"`
}

type UpdateEpicRequest struct {
	Name        *string    `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string    `json:"description,omitempty"`
	StartDate   *string    `json:"start_date,omitempty"`
	TargetDate  *string    `json:"target_date,omitempty"`
	Status      *string    `json:"status,omitempty"`
	LeadID      *uuid.UUID `json:"lead_id,omitempty"`
	SortOrder   *float64   `json:"sort_order,omitempty"`
}

type EpicIssueRequest struct {
	IssueID uuid.UUID `json:"issue_id" validate:"required"`
}

// ListEpics handles GET .../projects/{projectID}/modules
func (a *API) ListEpics(w http.ResponseWriter, r *http.Request) error {
	projectID := getProjectID(r.Context())

	modules, err := a.queries.ListEpicsByProject(r.Context(), projectID)
	if err != nil {
		return internalServerError("failed to list epics")
	}

	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"results": modules,
	})
}

// CreateEpic handles POST .../projects/{projectID}/modules
func (a *API) CreateEpic(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)
	projectID := getProjectID(ctx)
	workspaceID := getWorkspaceID(ctx)

	var req CreateEpicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	status := "backlog"
	if req.Status != "" {
		if !validEpicStatuses[req.Status] {
			return badRequestError("invalid status")
		}
		status = req.Status
	}

	var startDate, targetDate *time.Time
	if req.StartDate != "" {
		t, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return badRequestError("invalid start_date format, expected YYYY-MM-DD")
		}
		startDate = &t
	}
	if req.TargetDate != "" {
		t, err := time.Parse("2006-01-02", req.TargetDate)
		if err != nil {
			return badRequestError("invalid target_date format, expected YYYY-MM-DD")
		}
		targetDate = &t
	}
	if startDate != nil && targetDate != nil && targetDate.Before(*startDate) {
		return badRequestError("target_date must not be before start_date")
	}

	module, err := a.queries.CreateEpic(ctx, projectID, workspaceID, req.Name, req.Description, startDate, targetDate, status, req.LeadID, &userID)
	if err != nil {
		return internalServerError("failed to create epic")
	}

	return sendJSON(w, http.StatusCreated, module)
}

// GetEpic handles GET .../modules/{epicID}
func (a *API) GetEpic(w http.ResponseWriter, r *http.Request) error {
	epicID, err := uuid.FromString(chi.URLParam(r, "epicID"))
	if err != nil {
		return badRequestError("invalid epic ID")
	}

	module, err := a.queries.GetEpicByID(r.Context(), epicID)
	if err != nil {
		return notFoundError("epic not found")
	}

	issues, err := a.queries.ListIssuesByEpic(r.Context(), epicID)
	if err != nil {
		return internalServerError("failed to list epic issues")
	}

	totalIssues, err := a.queries.CountIssuesByEpic(r.Context(), epicID)
	if err != nil {
		totalIssues = int64(len(issues))
	}

	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"epic":       module,
		"issues":       issues,
		"total_issues": totalIssues,
	})
}

// UpdateEpic handles PUT .../modules/{epicID}
func (a *API) UpdateEpic(w http.ResponseWriter, r *http.Request) error {
	epicID, err := uuid.FromString(chi.URLParam(r, "epicID"))
	if err != nil {
		return badRequestError("invalid epic ID")
	}

	var req UpdateEpicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	if req.Status != nil && !validEpicStatuses[*req.Status] {
		return badRequestError("invalid status")
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

	module, err := a.queries.UpdateEpic(r.Context(), epicID, req.Name, req.Description, startDate, targetDate, req.Status, req.LeadID, req.SortOrder)
	if err != nil {
		return internalServerError("failed to update epic")
	}

	return sendJSON(w, http.StatusOK, module)
}

// DeleteEpic handles DELETE .../modules/{epicID}
func (a *API) DeleteEpic(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	epicID, err := uuid.FromString(chi.URLParam(r, "epicID"))
	if err != nil {
		return badRequestError("invalid epic ID")
	}

	if getProjectRole(ctx) < RoleAdmin {
		return forbiddenError("only an admin can delete modules")
	}

	if err := a.queries.DeleteEpic(ctx, epicID); err != nil {
		return internalServerError("failed to delete epic")
	}

	return sendEmpty(w, http.StatusNoContent)
}

// AddIssueToEpic handles POST .../modules/{epicID}/issues
func (a *API) AddIssueToEpic(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)

	epicID, err := uuid.FromString(chi.URLParam(r, "epicID"))
	if err != nil {
		return badRequestError("invalid epic ID")
	}

	var req EpicIssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	if err := a.queries.AddIssueToEpic(ctx, epicID, req.IssueID); err != nil {
		return internalServerError("failed to add issue to epic")
	}

	// Notify issue assignees & subscribers
	epic, _ := a.queries.GetEpicByID(ctx, epicID)
	issue, _ := a.queries.GetIssueByID(ctx, req.IssueID)
	if epic != nil && issue != nil {
		assignees, _ := a.queries.ListIssueAssignees(ctx, req.IssueID)
		subscribers, _ := a.queries.ListIssueSubscribers(ctx, req.IssueID)
		receiverIDs := collectUniqueUserIDs(assignees, subscribers)
		a.notifyUsers(ctx, receiverIDs, userID, notifyParams{
			WorkspaceID: issue.WorkspaceID,
			ProjectID:   &issue.ProjectID,
			Title:       fmt.Sprintf("\"%s\" added to epic \"%s\"", issue.Name, epic.Name),
			Message:     "An issue you follow was added to an epic",
			EntityType:  strPtr("issue"),
			EntityID:    &req.IssueID,
		})
	}

	return sendEmpty(w, http.StatusCreated)
}

// RemoveIssueFromEpic handles DELETE .../modules/{epicID}/issues/{issueID}
func (a *API) RemoveIssueFromEpic(w http.ResponseWriter, r *http.Request) error {
	epicID, err := uuid.FromString(chi.URLParam(r, "epicID"))
	if err != nil {
		return badRequestError("invalid epic ID")
	}

	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}

	if err := a.queries.RemoveIssueFromEpic(r.Context(), epicID, issueID); err != nil {
		return internalServerError("failed to remove issue from epic")
	}

	return sendEmpty(w, http.StatusNoContent)
}

type EpicMemberRequest struct {
	MemberID uuid.UUID `json:"member_id" validate:"required"`
	Role     string    `json:"role,omitempty"`
}

// ListEpicMembers handles GET .../epics/{epicID}/members
func (a *API) ListEpicMembers(w http.ResponseWriter, r *http.Request) error {
	epicID, err := uuid.FromString(chi.URLParam(r, "epicID"))
	if err != nil {
		return badRequestError("invalid epic ID")
	}

	members, err := a.queries.ListEpicMembers(r.Context(), epicID)
	if err != nil {
		return internalServerError("failed to list epic members")
	}

	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"results": members,
	})
}

// AddEpicMember handles POST .../epics/{epicID}/members
func (a *API) AddEpicMember(w http.ResponseWriter, r *http.Request) error {
	epicID, err := uuid.FromString(chi.URLParam(r, "epicID"))
	if err != nil {
		return badRequestError("invalid epic ID")
	}

	var req EpicMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	role := req.Role
	if role == "" {
		role = "member"
	}

	member, err := a.queries.AddEpicMember(r.Context(), epicID, req.MemberID, role)
	if err != nil {
		return internalServerError("failed to add epic member")
	}

	return sendJSON(w, http.StatusCreated, member)
}

// RemoveEpicMember handles DELETE .../epics/{epicID}/members/{memberID}
func (a *API) RemoveEpicMember(w http.ResponseWriter, r *http.Request) error {
	epicID, err := uuid.FromString(chi.URLParam(r, "epicID"))
	if err != nil {
		return badRequestError("invalid epic ID")
	}

	memberID, err := uuid.FromString(chi.URLParam(r, "memberID"))
	if err != nil {
		return badRequestError("invalid member ID")
	}

	if err := a.queries.RemoveEpicMember(r.Context(), epicID, memberID); err != nil {
		return internalServerError("failed to remove epic member")
	}

	return sendEmpty(w, http.StatusNoContent)
}
