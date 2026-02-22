package api

import (
	"encoding/json"
	"net/http"
)

var validEstimateTypes = map[string]bool{
	"points": true, "categories": true, "custom": true,
}

type CreateEstimateSystemRequest struct {
	Type      string          `json:"type" validate:"required"`
	Estimates json.RawMessage `json:"estimates" validate:"required"`
}

// GetEstimateSystem handles GET .../projects/{projectID}/estimates
func (a *API) GetEstimateSystem(w http.ResponseWriter, r *http.Request) error {
	projectID := getProjectID(r.Context())

	es, err := a.queries.GetEstimateSystemByProject(r.Context(), projectID)
	if err != nil {
		// Return default if none configured
		return sendJSON(w, http.StatusOK, map[string]interface{}{
			"type":      "points",
			"estimates": json.RawMessage(`[{"key":"1","value":"1"},{"key":"2","value":"2"},{"key":"3","value":"3"},{"key":"5","value":"5"},{"key":"8","value":"8"},{"key":"13","value":"13"}]`),
		})
	}
	return sendJSON(w, http.StatusOK, es)
}

// CreateOrUpdateEstimateSystem handles PUT .../projects/{projectID}/estimates
func (a *API) CreateOrUpdateEstimateSystem(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	if getProjectRole(ctx) < RoleAdmin {
		return forbiddenError("only project admins can configure estimates")
	}
	projectID := getProjectID(ctx)
	workspaceID := getWorkspaceID(ctx)

	var req CreateEstimateSystemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if !validEstimateTypes[req.Type] {
		return badRequestError("invalid estimate type: %s", req.Type)
	}

	existing, err := a.queries.GetEstimateSystemByProject(ctx, projectID)
	if err != nil {
		// Create new
		es, err := a.queries.CreateEstimateSystem(ctx, projectID, workspaceID, req.Type, req.Estimates)
		if err != nil {
			return internalServerError("failed to create estimate system")
		}
		return sendJSON(w, http.StatusCreated, es)
	}

	// Update existing
	es, err := a.queries.UpdateEstimateSystem(ctx, existing.ID, &req.Type, req.Estimates)
	if err != nil {
		return internalServerError("failed to update estimate system")
	}
	return sendJSON(w, http.StatusOK, es)
}
