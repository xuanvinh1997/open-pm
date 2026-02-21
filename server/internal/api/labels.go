package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

type CreateLabelRequest struct {
	Name        string     `json:"name" validate:"required,min=1"`
	Description string     `json:"description,omitempty"`
	Color       string     `json:"color" validate:"required"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	SortOrder   float64    `json:"sort_order,omitempty"`
}

// ListLabels handles GET .../projects/{projectID}/labels
func (a *API) ListLabels(w http.ResponseWriter, r *http.Request) error {
	projectID := getProjectID(r.Context())
	labels, err := a.queries.ListLabelsByProject(r.Context(), projectID)
	if err != nil {
		return internalServerError("failed to list labels")
	}
	return sendJSON(w, http.StatusOK, map[string]interface{}{"results": labels})
}

// CreateLabel handles POST .../projects/{projectID}/labels
func (a *API) CreateLabel(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	projectID := getProjectID(ctx)
	workspaceID := getWorkspaceID(ctx)

	var req CreateLabelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	if req.SortOrder == 0 {
		req.SortOrder = 65535
	}

	label, err := a.queries.CreateLabel(ctx, &projectID, workspaceID, req.ParentID, req.Name, req.Description, req.Color, req.SortOrder)
	if err != nil {
		return internalServerError("failed to create label")
	}
	return sendJSON(w, http.StatusCreated, label)
}

// UpdateLabel handles PUT .../labels/{labelID}
func (a *API) UpdateLabel(w http.ResponseWriter, r *http.Request) error {
	labelID, err := uuid.FromString(chi.URLParam(r, "labelID"))
	if err != nil {
		return badRequestError("invalid label ID")
	}

	var body struct {
		Name        *string    `json:"name,omitempty"`
		Description *string    `json:"description,omitempty"`
		Color       *string    `json:"color,omitempty"`
		ParentID    *uuid.UUID `json:"parent_id,omitempty"`
		SortOrder   *float64   `json:"sort_order,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return badRequestError("invalid request body")
	}

	label, err := a.queries.UpdateLabel(r.Context(), labelID, body.Name, body.Description, body.Color, body.ParentID, body.SortOrder)
	if err != nil {
		return internalServerError("failed to update label")
	}
	return sendJSON(w, http.StatusOK, label)
}

// DeleteLabel handles DELETE .../labels/{labelID}
func (a *API) DeleteLabel(w http.ResponseWriter, r *http.Request) error {
	labelID, err := uuid.FromString(chi.URLParam(r, "labelID"))
	if err != nil {
		return badRequestError("invalid label ID")
	}

	if err := a.queries.DeleteLabel(r.Context(), labelID); err != nil {
		return internalServerError("failed to delete label")
	}
	return sendEmpty(w, http.StatusNoContent)
}
