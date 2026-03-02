package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

type CreatePageRequest struct {
	Name            string          `json:"name" validate:"required,min=1,max=255"`
	DescriptionHTML string          `json:"description_html,omitempty"`
	DescriptionJSON json.RawMessage `json:"description_json,omitempty"`
	ParentID        *uuid.UUID      `json:"parent_id,omitempty"`
	Color           string          `json:"color,omitempty"`
}

type UpdatePageRequest struct {
	Name               *string         `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	DescriptionHTML    *string         `json:"description_html,omitempty"`
	DescriptionJSON    json.RawMessage `json:"description_json,omitempty"`
	DescriptionStripped *string        `json:"description_stripped,omitempty"`
	Color              *string         `json:"color,omitempty"`
	IsLocked           *bool           `json:"is_locked,omitempty"`
	ParentID           *uuid.UUID      `json:"parent_id,omitempty"`
}

// ListPages handles GET .../projects/{projectID}/pages
func (a *API) ListPages(w http.ResponseWriter, r *http.Request) error {
	projectID := getProjectID(r.Context())

	pages, err := a.queries.ListPagesByProject(r.Context(), projectID)
	if err != nil {
		return internalServerError("failed to list pages")
	}

	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"results": pages,
	})
}

// CreatePage handles POST .../projects/{projectID}/pages
func (a *API) CreatePage(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)
	projectID := getProjectID(ctx)
	workspaceID := getWorkspaceID(ctx)

	var req CreatePageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	page, err := a.queries.CreatePage(ctx, &projectID, workspaceID, req.Name, req.DescriptionHTML, req.DescriptionJSON, "", userID, req.ParentID)
	if err != nil {
		return internalServerError("failed to create page")
	}

	return sendJSON(w, http.StatusCreated, page)
}

// GetPage handles GET .../pages/{pageID}
func (a *API) GetPage(w http.ResponseWriter, r *http.Request) error {
	pageID, err := uuid.FromString(chi.URLParam(r, "pageID"))
	if err != nil {
		return badRequestError("invalid page ID")
	}

	page, err := a.queries.GetPageByID(r.Context(), pageID)
	if err != nil {
		return notFoundError("page not found")
	}

	return sendJSON(w, http.StatusOK, page)
}

// UpdatePage handles PUT .../pages/{pageID}
func (a *API) UpdatePage(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	pageID, err := uuid.FromString(chi.URLParam(r, "pageID"))
	if err != nil {
		return badRequestError("invalid page ID")
	}

	// Check if page is locked
	existing, err := a.queries.GetPageByID(ctx, pageID)
	if err != nil {
		return notFoundError("page not found")
	}
	if existing.IsLocked && existing.OwnedBy != getUserID(ctx) && getProjectRole(ctx) < RoleAdmin {
		return forbiddenError("this page is locked and can only be edited by the owner or an admin")
	}

	var req UpdatePageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	page, err := a.queries.UpdatePage(ctx, pageID, req.Name, req.DescriptionHTML, req.DescriptionJSON, req.DescriptionStripped, req.Color, req.IsLocked, req.ParentID)
	if err != nil {
		return internalServerError("failed to update page")
	}

	return sendJSON(w, http.StatusOK, page)
}

// DeletePage handles DELETE .../pages/{pageID}
func (a *API) DeletePage(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	pageID, err := uuid.FromString(chi.URLParam(r, "pageID"))
	if err != nil {
		return badRequestError("invalid page ID")
	}

	// Check ownership: only admin or page owner can delete
	page, err := a.queries.GetPageByID(ctx, pageID)
	if err != nil {
		return notFoundError("page not found")
	}
	if getProjectRole(ctx) < RoleAdmin && page.OwnedBy != getUserID(ctx) {
		return forbiddenError("only the owner or an admin can delete this page")
	}

	if err := a.queries.DeletePage(ctx, pageID); err != nil {
		return internalServerError("failed to delete page")
	}

	return sendEmpty(w, http.StatusNoContent)
}
