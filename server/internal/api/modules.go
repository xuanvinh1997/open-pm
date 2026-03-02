package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

var validModuleStatuses = map[string]bool{
	"backlog":     true,
	"planned":     true,
	"in-progress": true,
	"paused":      true,
	"completed":   true,
	"cancelled":   true,
}

type CreateModuleRequest struct {
	Name        string     `json:"name" validate:"required,min=1,max=255"`
	Description string     `json:"description,omitempty"`
	StartDate   string     `json:"start_date,omitempty"`
	TargetDate  string     `json:"target_date,omitempty"`
	Status      string     `json:"status,omitempty"`
	LeadID      *uuid.UUID `json:"lead_id,omitempty"`
}

type UpdateModuleRequest struct {
	Name        *string    `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string    `json:"description,omitempty"`
	StartDate   *string    `json:"start_date,omitempty"`
	TargetDate  *string    `json:"target_date,omitempty"`
	Status      *string    `json:"status,omitempty"`
	LeadID      *uuid.UUID `json:"lead_id,omitempty"`
	SortOrder   *float64   `json:"sort_order,omitempty"`
}

type ModuleIssueRequest struct {
	IssueID uuid.UUID `json:"issue_id" validate:"required"`
}

// ListModules handles GET .../projects/{projectID}/modules
func (a *API) ListModules(w http.ResponseWriter, r *http.Request) error {
	projectID := getProjectID(r.Context())

	modules, err := a.queries.ListModulesByProject(r.Context(), projectID)
	if err != nil {
		return internalServerError("failed to list modules")
	}

	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"results": modules,
	})
}

// CreateModule handles POST .../projects/{projectID}/modules
func (a *API) CreateModule(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)
	projectID := getProjectID(ctx)
	workspaceID := getWorkspaceID(ctx)

	var req CreateModuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	status := "backlog"
	if req.Status != "" {
		if !validModuleStatuses[req.Status] {
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

	module, err := a.queries.CreateModule(ctx, projectID, workspaceID, req.Name, req.Description, startDate, targetDate, status, req.LeadID, &userID)
	if err != nil {
		return internalServerError("failed to create module")
	}

	return sendJSON(w, http.StatusCreated, module)
}

// GetModule handles GET .../modules/{moduleID}
func (a *API) GetModule(w http.ResponseWriter, r *http.Request) error {
	moduleID, err := uuid.FromString(chi.URLParam(r, "moduleID"))
	if err != nil {
		return badRequestError("invalid module ID")
	}

	module, err := a.queries.GetModuleByID(r.Context(), moduleID)
	if err != nil {
		return notFoundError("module not found")
	}

	issues, err := a.queries.ListIssuesByModule(r.Context(), moduleID)
	if err != nil {
		return internalServerError("failed to list module issues")
	}

	totalIssues, err := a.queries.CountIssuesByModule(r.Context(), moduleID)
	if err != nil {
		totalIssues = int64(len(issues))
	}

	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"module":       module,
		"issues":       issues,
		"total_issues": totalIssues,
	})
}

// UpdateModule handles PUT .../modules/{moduleID}
func (a *API) UpdateModule(w http.ResponseWriter, r *http.Request) error {
	moduleID, err := uuid.FromString(chi.URLParam(r, "moduleID"))
	if err != nil {
		return badRequestError("invalid module ID")
	}

	var req UpdateModuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	if req.Status != nil && !validModuleStatuses[*req.Status] {
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

	module, err := a.queries.UpdateModule(r.Context(), moduleID, req.Name, req.Description, startDate, targetDate, req.Status, req.LeadID, req.SortOrder)
	if err != nil {
		return internalServerError("failed to update module")
	}

	return sendJSON(w, http.StatusOK, module)
}

// DeleteModule handles DELETE .../modules/{moduleID}
func (a *API) DeleteModule(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	moduleID, err := uuid.FromString(chi.URLParam(r, "moduleID"))
	if err != nil {
		return badRequestError("invalid module ID")
	}

	if getProjectRole(ctx) < RoleAdmin {
		return forbiddenError("only an admin can delete modules")
	}

	if err := a.queries.DeleteModule(ctx, moduleID); err != nil {
		return internalServerError("failed to delete module")
	}

	return sendEmpty(w, http.StatusNoContent)
}

// AddIssueToModule handles POST .../modules/{moduleID}/issues
func (a *API) AddIssueToModule(w http.ResponseWriter, r *http.Request) error {
	moduleID, err := uuid.FromString(chi.URLParam(r, "moduleID"))
	if err != nil {
		return badRequestError("invalid module ID")
	}

	var req ModuleIssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	if err := a.queries.AddIssueToModule(r.Context(), moduleID, req.IssueID); err != nil {
		return internalServerError("failed to add issue to module")
	}

	return sendEmpty(w, http.StatusCreated)
}

// RemoveIssueFromModule handles DELETE .../modules/{moduleID}/issues/{issueID}
func (a *API) RemoveIssueFromModule(w http.ResponseWriter, r *http.Request) error {
	moduleID, err := uuid.FromString(chi.URLParam(r, "moduleID"))
	if err != nil {
		return badRequestError("invalid module ID")
	}

	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}

	if err := a.queries.RemoveIssueFromModule(r.Context(), moduleID, issueID); err != nil {
		return internalServerError("failed to remove issue from module")
	}

	return sendEmpty(w, http.StatusNoContent)
}
