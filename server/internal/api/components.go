package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

type CreateComponentRequest struct {
	Name                string     `json:"name" validate:"required,min=1,max=255"`
	Description         string     `json:"description,omitempty"`
	LeadID              *uuid.UUID `json:"lead_id,omitempty"`
	DefaultAssigneeType string     `json:"default_assignee_type,omitempty" validate:"omitempty,oneof=project_default component_lead unassigned"`
	SortOrder           float64    `json:"sort_order,omitempty"`
}

type UpdateComponentRequest struct {
	Name                *string    `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description         *string    `json:"description,omitempty"`
	LeadID              *uuid.UUID `json:"lead_id,omitempty"`
	DefaultAssigneeType *string    `json:"default_assignee_type,omitempty" validate:"omitempty,oneof=project_default component_lead unassigned"`
	SortOrder           *float64   `json:"sort_order,omitempty"`
}

type ComponentIssueRequest struct {
	ComponentID uuid.UUID `json:"component_id" validate:"required"`
}

// ListComponents handles GET .../projects/{projectID}/components
func (a *API) ListComponents(w http.ResponseWriter, r *http.Request) error {
	projectID := getProjectID(r.Context())
	components, err := a.queries.ListComponentsByProject(r.Context(), projectID)
	if err != nil {
		return internalServerError("failed to list components")
	}
	return sendJSON(w, http.StatusOK, map[string]interface{}{"results": components})
}

// CreateComponent handles POST .../projects/{projectID}/components
func (a *API) CreateComponent(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	projectID := getProjectID(ctx)
	workspaceID := getWorkspaceID(ctx)

	var req CreateComponentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	defaultAssigneeType := req.DefaultAssigneeType
	if defaultAssigneeType == "" {
		defaultAssigneeType = "project_default"
	}
	sortOrder := req.SortOrder
	if sortOrder == 0 {
		sortOrder = 65535
	}

	component, err := a.queries.CreateComponent(ctx, projectID, workspaceID, req.Name, req.Description, req.LeadID, defaultAssigneeType, sortOrder)
	if err != nil {
		return internalServerError("failed to create component")
	}
	return sendJSON(w, http.StatusCreated, component)
}

// UpdateComponent handles PUT .../projects/{projectID}/components/{componentID}
func (a *API) UpdateComponent(w http.ResponseWriter, r *http.Request) error {
	componentID, err := uuid.FromString(chi.URLParam(r, "componentID"))
	if err != nil {
		return badRequestError("invalid component ID")
	}

	var req UpdateComponentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	component, err := a.queries.UpdateComponent(r.Context(), componentID, req.Name, req.Description, req.LeadID, req.DefaultAssigneeType, req.SortOrder)
	if err != nil {
		return internalServerError("failed to update component")
	}
	return sendJSON(w, http.StatusOK, component)
}

// DeleteComponent handles DELETE .../projects/{projectID}/components/{componentID}
func (a *API) DeleteComponent(w http.ResponseWriter, r *http.Request) error {
	componentID, err := uuid.FromString(chi.URLParam(r, "componentID"))
	if err != nil {
		return badRequestError("invalid component ID")
	}

	if err := a.queries.DeleteComponent(r.Context(), componentID); err != nil {
		return internalServerError("failed to delete component")
	}
	return sendEmpty(w, http.StatusNoContent)
}

// AddIssueComponent handles POST .../issues/{issueID}/components
func (a *API) AddIssueComponent(w http.ResponseWriter, r *http.Request) error {
	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}

	var req ComponentIssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	if err := a.queries.AddIssueComponent(r.Context(), issueID, req.ComponentID); err != nil {
		return internalServerError("failed to add component to issue")
	}
	return sendEmpty(w, http.StatusCreated)
}

// RemoveIssueComponent handles DELETE .../issues/{issueID}/components/{componentID}
func (a *API) RemoveIssueComponent(w http.ResponseWriter, r *http.Request) error {
	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}
	componentID, err := uuid.FromString(chi.URLParam(r, "componentID"))
	if err != nil {
		return badRequestError("invalid component ID")
	}

	if err := a.queries.RemoveIssueComponent(r.Context(), issueID, componentID); err != nil {
		return internalServerError("failed to remove component from issue")
	}
	return sendEmpty(w, http.StatusNoContent)
}
