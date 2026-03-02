package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

type CreateSprintRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description,omitempty"`
	StartDate   string `json:"start_date,omitempty"`
	EndDate     string `json:"end_date,omitempty"`
}

type UpdateSprintRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string  `json:"description,omitempty"`
	StartDate   *string  `json:"start_date,omitempty"`
	EndDate     *string  `json:"end_date,omitempty"`
	SortOrder   *float64 `json:"sort_order,omitempty"`
}

type SprintIssueRequest struct {
	IssueID uuid.UUID `json:"issue_id" validate:"required"`
}

// ListSprints handles GET .../projects/{projectID}/cycles
func (a *API) ListSprints(w http.ResponseWriter, r *http.Request) error {
	projectID := getProjectID(r.Context())

	cycles, err := a.queries.ListSprintsByProject(r.Context(), projectID)
	if err != nil {
		return internalServerError("failed to list sprints")
	}

	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"results": cycles,
	})
}

// CreateSprint handles POST .../projects/{projectID}/cycles
func (a *API) CreateSprint(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)
	projectID := getProjectID(ctx)
	workspaceID := getWorkspaceID(ctx)

	var req CreateSprintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	var startDate, endDate *time.Time
	if req.StartDate != "" {
		t, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return badRequestError("invalid start_date format, expected YYYY-MM-DD")
		}
		startDate = &t
	}
	if req.EndDate != "" {
		t, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return badRequestError("invalid end_date format, expected YYYY-MM-DD")
		}
		endDate = &t
	}
	if startDate != nil && endDate != nil && endDate.Before(*startDate) {
		return badRequestError("end_date must not be before start_date")
	}

	cycle, err := a.queries.CreateSprint(ctx, projectID, workspaceID, req.Name, req.Description, startDate, endDate, userID)
	if err != nil {
		return internalServerError("failed to create sprint")
	}

	return sendJSON(w, http.StatusCreated, cycle)
}

// GetSprint handles GET .../cycles/{sprintID}
func (a *API) GetSprint(w http.ResponseWriter, r *http.Request) error {
	sprintID, err := uuid.FromString(chi.URLParam(r, "sprintID"))
	if err != nil {
		return badRequestError("invalid sprint ID")
	}

	cycle, err := a.queries.GetSprintByID(r.Context(), sprintID)
	if err != nil {
		return notFoundError("sprint not found")
	}

	issues, err := a.queries.ListIssuesBySprint(r.Context(), sprintID)
	if err != nil {
		return internalServerError("failed to list cycle issues")
	}

	totalIssues, err := a.queries.CountIssuesBySprint(r.Context(), sprintID)
	if err != nil {
		totalIssues = int64(len(issues))
	}

	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"sprint":       cycle,
		"issues":      issues,
		"total_issues": totalIssues,
	})
}

// UpdateSprint handles PUT .../cycles/{sprintID}
func (a *API) UpdateSprint(w http.ResponseWriter, r *http.Request) error {
	sprintID, err := uuid.FromString(chi.URLParam(r, "sprintID"))
	if err != nil {
		return badRequestError("invalid sprint ID")
	}

	var req UpdateSprintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	var startDate, endDate *time.Time
	if req.StartDate != nil {
		t, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return badRequestError("invalid start_date format, expected YYYY-MM-DD")
		}
		startDate = &t
	}
	if req.EndDate != nil {
		t, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return badRequestError("invalid end_date format, expected YYYY-MM-DD")
		}
		endDate = &t
	}
	if startDate != nil && endDate != nil && endDate.Before(*startDate) {
		return badRequestError("end_date must not be before start_date")
	}

	cycle, err := a.queries.UpdateSprint(r.Context(), sprintID, req.Name, req.Description, startDate, endDate, req.SortOrder)
	if err != nil {
		return internalServerError("failed to update sprint")
	}

	return sendJSON(w, http.StatusOK, cycle)
}

// DeleteSprint handles DELETE .../cycles/{sprintID}
func (a *API) DeleteSprint(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	sprintID, err := uuid.FromString(chi.URLParam(r, "sprintID"))
	if err != nil {
		return badRequestError("invalid sprint ID")
	}

	// Check ownership: only admin or cycle owner can delete
	cycle, err := a.queries.GetSprintByID(ctx, sprintID)
	if err != nil {
		return notFoundError("sprint not found")
	}
	if getProjectRole(ctx) < RoleAdmin && cycle.OwnedBy != getUserID(ctx) {
		return forbiddenError("only the owner or an admin can delete this cycle")
	}

	if err := a.queries.DeleteSprint(ctx, sprintID); err != nil {
		return internalServerError("failed to delete sprint")
	}

	return sendEmpty(w, http.StatusNoContent)
}

// AddIssueToSprint handles POST .../cycles/{sprintID}/issues
func (a *API) AddIssueToSprint(w http.ResponseWriter, r *http.Request) error {
	sprintID, err := uuid.FromString(chi.URLParam(r, "sprintID"))
	if err != nil {
		return badRequestError("invalid sprint ID")
	}

	var req SprintIssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	if err := a.queries.AddIssueToSprint(r.Context(), sprintID, req.IssueID); err != nil {
		return internalServerError("failed to add issue to cycle")
	}

	return sendEmpty(w, http.StatusCreated)
}

// RemoveIssueFromSprint handles DELETE .../cycles/{sprintID}/issues/{issueID}
func (a *API) RemoveIssueFromSprint(w http.ResponseWriter, r *http.Request) error {
	sprintID, err := uuid.FromString(chi.URLParam(r, "sprintID"))
	if err != nil {
		return badRequestError("invalid sprint ID")
	}

	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}

	if err := a.queries.RemoveIssueFromSprint(r.Context(), sprintID, issueID); err != nil {
		return internalServerError("failed to remove issue from cycle")
	}

	return sendEmpty(w, http.StatusNoContent)
}
