package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

type CreateProjectRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Identifier  string  `json:"identifier" validate:"required,min=1,max=12"`
	Description string  `json:"description,omitempty"`
	Emoji       *string `json:"emoji,omitempty"`
}

type UpdateProjectRequest struct {
	Name            *string    `json:"name,omitempty"`
	Description     *string    `json:"description,omitempty"`
	Identifier      *string    `json:"identifier,omitempty"`
	Emoji           *string    `json:"emoji,omitempty"`
	CoverImageURL   *string    `json:"cover_image_url,omitempty"`
	DefaultAssignee *uuid.UUID `json:"default_assignee_id,omitempty"`
	ProjectLead     *uuid.UUID `json:"project_lead_id,omitempty"`
	Network         *int16     `json:"network,omitempty"`
	SprintView       *bool      `json:"sprint_view,omitempty"`
	EpicView      *bool      `json:"epic_view,omitempty"`
	PageView        *bool      `json:"page_view,omitempty"`
	InboxView       *bool      `json:"inbox_view,omitempty"`
}

// ListProjects handles GET /api/v1/workspaces/{workspaceSlug}/projects
func (a *API) ListProjects(w http.ResponseWriter, r *http.Request) error {
	workspaceID := getWorkspaceID(r.Context())
	projects, err := a.queries.ListProjectsByWorkspace(r.Context(), workspaceID)
	if err != nil {
		return internalServerError("failed to list projects")
	}
	return sendJSON(w, http.StatusOK, map[string]interface{}{"results": projects})
}

// CreateProject handles POST /api/v1/workspaces/{workspaceSlug}/projects
func (a *API) CreateProject(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)
	workspaceID := getWorkspaceID(ctx)

	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	project, err := a.queries.CreateProject(ctx, workspaceID, req.Name, req.Description, req.Identifier, &userID, &userID)
	if err != nil {
		return internalServerError("failed to create project")
	}

	// Add creator as project admin
	if _, err := a.queries.CreateProjectMember(ctx, project.ID, userID, RoleAdmin); err != nil {
		log.Warn().Err(err).Str("project_id", project.ID.String()).Str("user_id", userID.String()).Msg("failed to add creator as project admin")
	}

	// Create default states
	a.createDefaultStates(ctx, project.ID, workspaceID)

	return sendJSON(w, http.StatusCreated, project)
}

// GetProject handles GET /api/v1/workspaces/{workspaceSlug}/projects/{projectID}
func (a *API) GetProject(w http.ResponseWriter, r *http.Request) error {
	projectID := getProjectID(r.Context())
	project, err := a.queries.GetProjectByID(r.Context(), projectID)
	if err != nil {
		return notFoundError("project not found")
	}
	return sendJSON(w, http.StatusOK, project)
}

// UpdateProject handles PUT /api/v1/workspaces/{workspaceSlug}/projects/{projectID}
func (a *API) UpdateProject(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	if getProjectRole(ctx) < RoleAdmin {
		return forbiddenError("only project admins can update project settings")
	}
	projectID := getProjectID(ctx)

	var req UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}

	project, err := a.queries.UpdateProject(ctx, projectID,
		req.Name, req.Description, req.Identifier, req.Emoji,
		req.CoverImageURL, req.DefaultAssignee, req.ProjectLead,
		req.Network, req.SprintView, req.EpicView, req.PageView, req.InboxView,
	)
	if err != nil {
		return internalServerError("failed to update project")
	}
	return sendJSON(w, http.StatusOK, project)
}

// DeleteProject handles DELETE /api/v1/workspaces/{workspaceSlug}/projects/{projectID}
func (a *API) DeleteProject(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	role := getWorkspaceRole(ctx)
	if role < RoleAdmin {
		return forbiddenError("only admins can delete projects")
	}

	projectID := getProjectID(ctx)
	if err := a.queries.DeleteProject(ctx, projectID); err != nil {
		return internalServerError("failed to delete project")
	}
	return sendEmpty(w, http.StatusNoContent)
}

// ListProjectMembers handles GET .../projects/{projectID}/members
func (a *API) ListProjectMembers(w http.ResponseWriter, r *http.Request) error {
	projectID := getProjectID(r.Context())
	members, err := a.queries.ListProjectMembers(r.Context(), projectID)
	if err != nil {
		return internalServerError("failed to list project members")
	}
	return sendJSON(w, http.StatusOK, map[string]interface{}{"results": members})
}

// AddProjectMember handles POST .../projects/{projectID}/members
func (a *API) AddProjectMember(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	callerRole := getProjectRole(ctx)
	if callerRole < RoleAdmin {
		return forbiddenError("only project admins can add members")
	}
	projectID := getProjectID(ctx)

	var body struct {
		UserID uuid.UUID `json:"user_id" validate:"required"`
		Role   int16     `json:"role" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return badRequestError("invalid request body")
	}

	if !isValidRole(body.Role) {
		return badRequestError("invalid role value")
	}

	// Cannot assign a role higher than your own
	if body.Role > callerRole {
		return forbiddenError("cannot assign a role higher than your own")
	}

	member, err := a.queries.CreateProjectMember(ctx, projectID, body.UserID, body.Role)
	if err != nil {
		return internalServerError("failed to add project member")
	}
	return sendJSON(w, http.StatusCreated, member)
}

// UpdateProjectMember handles PUT .../projects/{projectID}/members/{memberID}
func (a *API) UpdateProjectMember(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	callerRole := getProjectRole(ctx)
	if callerRole < RoleAdmin {
		return forbiddenError("only project admins can update member roles")
	}
	projectID := getProjectID(ctx)
	memberID, err := uuid.FromString(chi.URLParam(r, "memberID"))
	if err != nil {
		return badRequestError("invalid member ID")
	}

	var body struct {
		Role int16 `json:"role" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return badRequestError("invalid request body")
	}

	if !isValidRole(body.Role) {
		return badRequestError("invalid role value")
	}

	// Cannot assign a role higher than your own
	if body.Role > callerRole {
		return forbiddenError("cannot assign a role higher than your own")
	}

	member, err := a.queries.UpdateProjectMemberRole(ctx, projectID, memberID, body.Role)
	if err != nil {
		return internalServerError("failed to update member role")
	}
	return sendJSON(w, http.StatusOK, member)
}

// RemoveProjectMember handles DELETE .../projects/{projectID}/members/{memberID}
func (a *API) RemoveProjectMember(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	if getProjectRole(ctx) < RoleAdmin {
		return forbiddenError("only project admins can remove members")
	}
	projectID := getProjectID(ctx)
	memberID, err := uuid.FromString(chi.URLParam(r, "memberID"))
	if err != nil {
		return badRequestError("invalid member ID")
	}

	// Prevent removing yourself
	if memberID == getUserID(ctx) {
		return badRequestError("cannot remove yourself from a project")
	}

	if err := a.queries.DeleteProjectMember(ctx, projectID, memberID); err != nil {
		return internalServerError("failed to remove member")
	}
	return sendEmpty(w, http.StatusNoContent)
}

// createDefaultStates creates the standard workflow states for a new project.
func (a *API) createDefaultStates(ctx context.Context, projectID, workspaceID uuid.UUID) {
	type stateInit struct {
		Name     string
		Color    string
		Group    string
		Sequence float64
		Default  bool
	}

	defaults := []stateInit{
		{Name: "Backlog", Color: "#A3A3A3", Group: "backlog", Sequence: 1000, Default: true},
		{Name: "Todo", Color: "#3A3A3A", Group: "unstarted", Sequence: 2000},
		{Name: "In Progress", Color: "#F59E0B", Group: "started", Sequence: 3000},
		{Name: "Done", Color: "#22C55E", Group: "completed", Sequence: 4000},
		{Name: "Cancelled", Color: "#EF4444", Group: "cancelled", Sequence: 5000},
	}

	for _, s := range defaults {
		if _, err := a.queries.CreateState(ctx, projectID, workspaceID, s.Name, "", s.Color, s.Group, s.Sequence, s.Default); err != nil {
			log.Warn().Err(err).Str("project_id", projectID.String()).Str("state", s.Name).Msg("failed to create default state")
		}
	}
}
