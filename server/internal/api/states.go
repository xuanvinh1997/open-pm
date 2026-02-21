package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

type CreateStateRequest struct {
	Name        string  `json:"name" validate:"required,min=1"`
	Description string  `json:"description,omitempty"`
	Color       string  `json:"color" validate:"required"`
	Group       string  `json:"group" validate:"required,oneof=backlog unstarted started completed cancelled triage"`
	Sequence    float64 `json:"sequence,omitempty"`
	IsDefault   bool    `json:"is_default,omitempty"`
}

// ListStates handles GET .../projects/{projectID}/states
func (a *API) ListStates(w http.ResponseWriter, r *http.Request) error {
	projectID := getProjectID(r.Context())
	states, err := a.queries.ListStatesByProject(r.Context(), projectID)
	if err != nil {
		return internalServerError("failed to list states")
	}
	return sendJSON(w, http.StatusOK, map[string]interface{}{"results": states})
}

// CreateState handles POST .../projects/{projectID}/states
func (a *API) CreateState(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	projectID := getProjectID(ctx)
	workspaceID := getWorkspaceID(ctx)

	var req CreateStateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	if req.Sequence == 0 {
		req.Sequence = 65535
	}

	state, err := a.queries.CreateState(ctx, projectID, workspaceID, req.Name, req.Description, req.Color, req.Group, req.Sequence, req.IsDefault)
	if err != nil {
		return internalServerError("failed to create state")
	}
	return sendJSON(w, http.StatusCreated, state)
}

// UpdateState handles PUT .../states/{stateID}
func (a *API) UpdateState(w http.ResponseWriter, r *http.Request) error {
	stateID, err := uuid.FromString(chi.URLParam(r, "stateID"))
	if err != nil {
		return badRequestError("invalid state ID")
	}

	var body struct {
		Name        *string  `json:"name,omitempty"`
		Description *string  `json:"description,omitempty"`
		Color       *string  `json:"color,omitempty"`
		Group       *string  `json:"group,omitempty"`
		Sequence    *float64 `json:"sequence,omitempty"`
		IsDefault   *bool    `json:"is_default,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return badRequestError("invalid request body")
	}

	state, err := a.queries.UpdateState(r.Context(), stateID, body.Name, body.Description, body.Color, body.Group, body.Sequence, body.IsDefault)
	if err != nil {
		return internalServerError("failed to update state")
	}
	return sendJSON(w, http.StatusOK, state)
}

// DeleteState handles DELETE .../states/{stateID}
func (a *API) DeleteState(w http.ResponseWriter, r *http.Request) error {
	stateID, err := uuid.FromString(chi.URLParam(r, "stateID"))
	if err != nil {
		return badRequestError("invalid state ID")
	}

	if err := a.queries.DeleteState(r.Context(), stateID); err != nil {
		return internalServerError("failed to delete state")
	}
	return sendEmpty(w, http.StatusNoContent)
}
