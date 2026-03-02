package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

var slugRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$`)

type CreateWorkspaceRequest struct {
	Name             string `json:"name" validate:"required,min=1,max=80"`
	Slug             string `json:"slug" validate:"required,min=2,max=48"`
	OrganizationSize string `json:"organization_size,omitempty"`
}

type UpdateWorkspaceRequest struct {
	Name             *string `json:"name,omitempty"`
	Slug             *string `json:"slug,omitempty"`
	LogoURL          *string `json:"logo_url,omitempty"`
	OrganizationSize *string `json:"organization_size,omitempty"`
}

// ListWorkspaces handles GET /api/v1/workspaces
func (a *API) ListWorkspaces(w http.ResponseWriter, r *http.Request) error {
	userID := getUserID(r.Context())
	workspaces, err := a.queries.ListWorkspacesByUser(r.Context(), userID)
	if err != nil {
		return internalServerError("failed to list workspaces")
	}
	return sendJSON(w, http.StatusOK, map[string]interface{}{"results": workspaces})
}

// CreateWorkspace handles POST /api/v1/workspaces
func (a *API) CreateWorkspace(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)

	var req CreateWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}

	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	if !slugRegex.MatchString(req.Slug) {
		return badRequestError("slug must contain only lowercase letters, numbers, and hyphens")
	}

	// Check if slug is taken
	_, err := a.queries.GetWorkspaceBySlug(ctx, req.Slug)
	if err == nil {
		return conflictError("workspace slug already taken")
	}

	workspace, err := a.queries.CreateWorkspace(ctx, req.Name, req.Slug, userID, &req.OrganizationSize)
	if err != nil {
		return internalServerError("failed to create workspace")
	}

	// Add creator as owner
	_, err = a.queries.CreateWorkspaceMember(ctx, workspace.ID, userID, RoleOwner)
	if err != nil {
		return internalServerError("failed to add workspace owner")
	}

	return sendJSON(w, http.StatusCreated, workspace)
}

// GetWorkspace handles GET /api/v1/workspaces/{workspaceSlug}
func (a *API) GetWorkspace(w http.ResponseWriter, r *http.Request) error {
	workspaceID := getWorkspaceID(r.Context())
	workspace, err := a.queries.GetWorkspaceByID(r.Context(), workspaceID)
	if err != nil {
		return notFoundError("workspace not found")
	}
	return sendJSON(w, http.StatusOK, workspace)
}

// UpdateWorkspace handles PUT /api/v1/workspaces/{workspaceSlug}
func (a *API) UpdateWorkspace(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	role := getWorkspaceRole(ctx)
	if role < RoleAdmin {
		return forbiddenError("only admins can update workspace settings")
	}

	workspaceID := getWorkspaceID(ctx)

	var req UpdateWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}

	if req.Slug != nil && !slugRegex.MatchString(*req.Slug) {
		return badRequestError("slug must contain only lowercase letters, numbers, and hyphens")
	}

	workspace, err := a.queries.UpdateWorkspace(ctx, workspaceID, req.Name, req.Slug, req.LogoURL, req.OrganizationSize)
	if err != nil {
		return internalServerError("failed to update workspace")
	}
	return sendJSON(w, http.StatusOK, workspace)
}

// DeleteWorkspace handles DELETE /api/v1/workspaces/{workspaceSlug}
func (a *API) DeleteWorkspace(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	role := getWorkspaceRole(ctx)
	if role < RoleOwner {
		return forbiddenError("only the owner can delete the workspace")
	}

	workspaceID := getWorkspaceID(ctx)
	if err := a.queries.DeleteWorkspace(ctx, workspaceID); err != nil {
		return internalServerError("failed to delete workspace")
	}
	return sendEmpty(w, http.StatusNoContent)
}

// ListWorkspaceMembers handles GET /api/v1/workspaces/{workspaceSlug}/members
func (a *API) ListWorkspaceMembers(w http.ResponseWriter, r *http.Request) error {
	workspaceID := getWorkspaceID(r.Context())
	members, err := a.queries.ListWorkspaceMembers(r.Context(), workspaceID)
	if err != nil {
		return internalServerError("failed to list workspace members")
	}
	return sendJSON(w, http.StatusOK, map[string]interface{}{"results": members})
}

// UpdateWorkspaceMember handles PUT /api/v1/workspaces/{workspaceSlug}/members/{memberID}
func (a *API) UpdateWorkspaceMember(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	callerRole := getWorkspaceRole(ctx)
	if callerRole < RoleAdmin {
		return forbiddenError("only admins can update member roles")
	}

	workspaceID := getWorkspaceID(ctx)
	userID := getUserID(ctx)
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

	// Validate role value
	if !isValidRole(body.Role) || body.Role == RoleOwner {
		return badRequestError("invalid role value")
	}

	// Cannot modify your own role
	if memberID == userID {
		return badRequestError("cannot modify your own role")
	}

	// Protect workspace owner from role changes
	workspace, err := a.queries.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		return internalServerError("failed to fetch workspace")
	}
	if memberID == workspace.OwnerID {
		return forbiddenError("cannot modify the workspace owner's role")
	}

	member, err := a.queries.UpdateWorkspaceMemberRole(ctx, workspaceID, memberID, body.Role)
	if err != nil {
		return internalServerError("failed to update member role")
	}
	return sendJSON(w, http.StatusOK, member)
}

// RemoveWorkspaceMember handles DELETE /api/v1/workspaces/{workspaceSlug}/members/{memberID}
func (a *API) RemoveWorkspaceMember(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	role := getWorkspaceRole(ctx)
	if role < RoleAdmin {
		return forbiddenError("only admins can remove members")
	}

	workspaceID := getWorkspaceID(ctx)
	userID := getUserID(ctx)
	memberID, err := uuid.FromString(chi.URLParam(r, "memberID"))
	if err != nil {
		return badRequestError("invalid member ID")
	}

	// Cannot remove yourself
	if memberID == userID {
		return badRequestError("cannot remove yourself from the workspace")
	}

	// Cannot remove the workspace owner
	workspace, err := a.queries.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		return internalServerError("failed to fetch workspace")
	}
	if memberID == workspace.OwnerID {
		return forbiddenError("cannot remove the workspace owner")
	}

	if err := a.queries.DeleteWorkspaceMember(ctx, workspaceID, memberID); err != nil {
		return internalServerError("failed to remove member")
	}
	return sendEmpty(w, http.StatusNoContent)
}

// CreateWorkspaceInvite handles POST /api/v1/workspaces/{workspaceSlug}/invites
func (a *API) CreateWorkspaceInvite(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	role := getWorkspaceRole(ctx)
	if role < RoleAdmin {
		return forbiddenError("only admins can invite members")
	}

	workspaceID := getWorkspaceID(ctx)

	var body struct {
		Email   string `json:"email" validate:"required,email"`
		Role    int16  `json:"role" validate:"required"`
		Message string `json:"message,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(body); err != nil {
		return validationError(err)
	}

	// Validate and cap invite role (cannot invite as owner)
	if !isValidRole(body.Role) || body.Role > RoleAdmin {
		return badRequestError("invite role must be guest, member, or admin")
	}

	token := generateToken()
	invite, err := a.queries.CreateWorkspaceInvite(ctx, workspaceID, body.Email, body.Role, token, &body.Message)
	if err != nil {
		return internalServerError("failed to create invite")
	}

	// Send invite email
	if a.mailer != nil {
		ws, _ := a.queries.GetWorkspaceByID(ctx, workspaceID)
		wsName := "a workspace"
		if ws != nil {
			wsName = ws.Name
		}
		acceptURL := fmt.Sprintf("%s/invites/%s", a.config.SiteURL, token)
		go func() {
			if err := a.mailer.Send(body.Email, fmt.Sprintf("You're invited to %s", wsName), "workspace_invite.html", map[string]string{
				"WorkspaceName": wsName,
				"Message":       body.Message,
				"AcceptURL":     acceptURL,
			}); err != nil {
				log.Warn().Err(err).Str("to", body.Email).Msg("failed to send invite email")
			}
		}()
	}

	return sendJSON(w, http.StatusCreated, invite)
}

// ListWorkspaceInvites handles GET /api/v1/workspaces/{workspaceSlug}/invites
func (a *API) ListWorkspaceInvites(w http.ResponseWriter, r *http.Request) error {
	if getWorkspaceRole(r.Context()) < RoleAdmin {
		return forbiddenError("only admins can view invites")
	}
	workspaceID := getWorkspaceID(r.Context())
	invites, err := a.queries.ListWorkspaceInvites(r.Context(), workspaceID)
	if err != nil {
		return internalServerError("failed to list invites")
	}
	return sendJSON(w, http.StatusOK, map[string]interface{}{"results": invites})
}

// AcceptWorkspaceInvite handles POST /api/v1/invites/{token}/accept
func (a *API) AcceptWorkspaceInvite(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)
	inviteToken := chi.URLParam(r, "token")

	invite, err := a.queries.GetWorkspaceInviteByToken(ctx, inviteToken)
	if err != nil {
		return notFoundError("invite not found")
	}

	if invite.Accepted {
		return badRequestError("invite already accepted")
	}

	// Accept invite
	_ = a.queries.AcceptWorkspaceInvite(ctx, invite.ID)

	// Add user as member
	_, err = a.queries.CreateWorkspaceMember(ctx, invite.WorkspaceID, userID, invite.Role)
	if err != nil {
		return internalServerError("failed to add member")
	}

	return sendJSON(w, http.StatusOK, map[string]string{"message": "invite accepted"})
}
