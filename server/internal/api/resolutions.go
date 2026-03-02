package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

type CreateIssueResolutionRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description string  `json:"description,omitempty"`
	IsDefault   bool    `json:"is_default,omitempty"`
	SortOrder   float64 `json:"sort_order,omitempty"`
}

type UpdateIssueResolutionRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string  `json:"description,omitempty"`
	IsDefault   *bool    `json:"is_default,omitempty"`
	SortOrder   *float64 `json:"sort_order,omitempty"`
}

// ListIssueResolutions handles GET .../projects/{projectID}/resolutions
func (a *API) ListIssueResolutions(w http.ResponseWriter, r *http.Request) error {
	projectID := getProjectID(r.Context())
	resolutions, err := a.queries.ListIssueResolutionsByProject(r.Context(), projectID)
	if err != nil {
		return internalServerError("failed to list resolutions")
	}
	return sendJSON(w, http.StatusOK, map[string]interface{}{"results": resolutions})
}

// CreateIssueResolution handles POST .../projects/{projectID}/resolutions
func (a *API) CreateIssueResolution(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	projectID := getProjectID(ctx)
	workspaceID := getWorkspaceID(ctx)

	var req CreateIssueResolutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	sortOrder := req.SortOrder
	if sortOrder == 0 {
		sortOrder = 65535
	}

	res, err := a.queries.CreateIssueResolution(ctx, projectID, workspaceID, req.Name, req.Description, req.IsDefault, sortOrder)
	if err != nil {
		return internalServerError("failed to create resolution")
	}
	return sendJSON(w, http.StatusCreated, res)
}

// UpdateIssueResolution handles PUT .../projects/{projectID}/resolutions/{resolutionID}
func (a *API) UpdateIssueResolution(w http.ResponseWriter, r *http.Request) error {
	resolutionID, err := uuid.FromString(chi.URLParam(r, "resolutionID"))
	if err != nil {
		return badRequestError("invalid resolution ID")
	}

	var req UpdateIssueResolutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	res, err := a.queries.UpdateIssueResolution(r.Context(), resolutionID, req.Name, req.Description, req.IsDefault, req.SortOrder)
	if err != nil {
		return internalServerError("failed to update resolution")
	}
	return sendJSON(w, http.StatusOK, res)
}

// DeleteIssueResolution handles DELETE .../projects/{projectID}/resolutions/{resolutionID}
func (a *API) DeleteIssueResolution(w http.ResponseWriter, r *http.Request) error {
	resolutionID, err := uuid.FromString(chi.URLParam(r, "resolutionID"))
	if err != nil {
		return badRequestError("invalid resolution ID")
	}

	if err := a.queries.DeleteIssueResolution(r.Context(), resolutionID); err != nil {
		return internalServerError("failed to delete resolution")
	}
	return sendEmpty(w, http.StatusNoContent)
}
