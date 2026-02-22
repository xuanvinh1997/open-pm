package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

type CreateCycleRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description,omitempty"`
	StartDate   string `json:"start_date,omitempty"`
	EndDate     string `json:"end_date,omitempty"`
}

type UpdateCycleRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string  `json:"description,omitempty"`
	StartDate   *string  `json:"start_date,omitempty"`
	EndDate     *string  `json:"end_date,omitempty"`
	SortOrder   *float64 `json:"sort_order,omitempty"`
}

type CycleIssueRequest struct {
	IssueID uuid.UUID `json:"issue_id" validate:"required"`
}

// ListCycles handles GET .../projects/{projectID}/cycles
func (a *API) ListCycles(w http.ResponseWriter, r *http.Request) error {
	projectID := getProjectID(r.Context())

	cycles, err := a.queries.ListCyclesByProject(r.Context(), projectID)
	if err != nil {
		return internalServerError("failed to list cycles")
	}

	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"results": cycles,
	})
}

// CreateCycle handles POST .../projects/{projectID}/cycles
func (a *API) CreateCycle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)
	projectID := getProjectID(ctx)
	workspaceID := getWorkspaceID(ctx)

	var req CreateCycleRequest
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

	cycle, err := a.queries.CreateCycle(ctx, projectID, workspaceID, req.Name, req.Description, startDate, endDate, userID)
	if err != nil {
		return internalServerError("failed to create cycle")
	}

	return sendJSON(w, http.StatusCreated, cycle)
}

// GetCycle handles GET .../cycles/{cycleID}
func (a *API) GetCycle(w http.ResponseWriter, r *http.Request) error {
	cycleID, err := uuid.FromString(chi.URLParam(r, "cycleID"))
	if err != nil {
		return badRequestError("invalid cycle ID")
	}

	cycle, err := a.queries.GetCycleByID(r.Context(), cycleID)
	if err != nil {
		return notFoundError("cycle not found")
	}

	issues, err := a.queries.ListIssuesByCycle(r.Context(), cycleID)
	if err != nil {
		return internalServerError("failed to list cycle issues")
	}

	totalIssues, err := a.queries.CountIssuesByCycle(r.Context(), cycleID)
	if err != nil {
		totalIssues = int64(len(issues))
	}

	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"cycle":       cycle,
		"issues":      issues,
		"total_issues": totalIssues,
	})
}

// UpdateCycle handles PUT .../cycles/{cycleID}
func (a *API) UpdateCycle(w http.ResponseWriter, r *http.Request) error {
	cycleID, err := uuid.FromString(chi.URLParam(r, "cycleID"))
	if err != nil {
		return badRequestError("invalid cycle ID")
	}

	var req UpdateCycleRequest
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

	cycle, err := a.queries.UpdateCycle(r.Context(), cycleID, req.Name, req.Description, startDate, endDate, req.SortOrder)
	if err != nil {
		return internalServerError("failed to update cycle")
	}

	return sendJSON(w, http.StatusOK, cycle)
}

// DeleteCycle handles DELETE .../cycles/{cycleID}
func (a *API) DeleteCycle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	cycleID, err := uuid.FromString(chi.URLParam(r, "cycleID"))
	if err != nil {
		return badRequestError("invalid cycle ID")
	}

	// Check ownership: only admin or cycle owner can delete
	cycle, err := a.queries.GetCycleByID(ctx, cycleID)
	if err != nil {
		return notFoundError("cycle not found")
	}
	if getProjectRole(ctx) < RoleAdmin && cycle.OwnedBy != getUserID(ctx) {
		return forbiddenError("only the owner or an admin can delete this cycle")
	}

	if err := a.queries.DeleteCycle(ctx, cycleID); err != nil {
		return internalServerError("failed to delete cycle")
	}

	return sendEmpty(w, http.StatusNoContent)
}

// AddIssueToCycle handles POST .../cycles/{cycleID}/issues
func (a *API) AddIssueToCycle(w http.ResponseWriter, r *http.Request) error {
	cycleID, err := uuid.FromString(chi.URLParam(r, "cycleID"))
	if err != nil {
		return badRequestError("invalid cycle ID")
	}

	var req CycleIssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	if err := a.queries.AddIssueToCycle(r.Context(), cycleID, req.IssueID); err != nil {
		return internalServerError("failed to add issue to cycle")
	}

	return sendEmpty(w, http.StatusCreated)
}

// RemoveIssueFromCycle handles DELETE .../cycles/{cycleID}/issues/{issueID}
func (a *API) RemoveIssueFromCycle(w http.ResponseWriter, r *http.Request) error {
	cycleID, err := uuid.FromString(chi.URLParam(r, "cycleID"))
	if err != nil {
		return badRequestError("invalid cycle ID")
	}

	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}

	if err := a.queries.RemoveIssueFromCycle(r.Context(), cycleID, issueID); err != nil {
		return internalServerError("failed to remove issue from cycle")
	}

	return sendEmpty(w, http.StatusNoContent)
}
