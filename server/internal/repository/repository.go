package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/open-pm/open-pm/server/internal/api"
)

// Repository implements api.Queries using direct pgx queries.
// Once sqlc is set up and generated, this will wrap the generated code.
type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// scan helpers
func scanUser(row pgx.Row) (*api.User, error) {
	var u api.User
	var deletedAt *time.Time
	var appMeta, userMeta json.RawMessage
	err := row.Scan(
		&u.ID, &u.Email, &u.EncryptedPassword, &u.EmailConfirmedAt,
		nil, // invited_at
		nil, // confirmation_token
		nil, // confirmation_sent_at
		nil, // recovery_token
		nil, // recovery_sent_at
		nil, // last_sign_in_at
		&u.FirstName, &u.LastName, &u.AvatarURL, &u.DisplayName,
		&u.IsActive, &appMeta, &userMeta,
		&u.CreatedAt, &u.UpdatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}
	u.UserMetadata = userMeta
	return &u, nil
}

// --- Users ---

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*api.User, error) {
	return scanUser(r.pool.QueryRow(ctx,
		`SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL`, id))
}

func (r *Repository) GetUserByEmail(ctx context.Context, email *string) (*api.User, error) {
	return scanUser(r.pool.QueryRow(ctx,
		`SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL`, email))
}

func (r *Repository) CreateUser(ctx context.Context, email *string, encryptedPassword *string, firstName, lastName, displayName string, emailConfirmedAt *time.Time) (*api.User, error) {
	return scanUser(r.pool.QueryRow(ctx,
		`INSERT INTO users (email, encrypted_password, first_name, last_name, display_name, email_confirmed_at)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`,
		email, encryptedPassword, firstName, lastName, displayName, emailConfirmedAt))
}

func (r *Repository) UpdateUser(ctx context.Context, id uuid.UUID, firstName, lastName, displayName, avatarURL *string) (*api.User, error) {
	return scanUser(r.pool.QueryRow(ctx,
		`UPDATE users SET
			first_name = COALESCE($2, first_name),
			last_name = COALESCE($3, last_name),
			display_name = COALESCE($4, display_name),
			avatar_url = COALESCE($5, avatar_url)
		 WHERE id = $1 AND deleted_at IS NULL RETURNING *`,
		id, firstName, lastName, displayName, avatarURL))
}

func (r *Repository) SetUserLastSignIn(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET last_sign_in_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *Repository) SetUserEmailConfirmed(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET email_confirmed_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *Repository) SetUserConfirmationToken(ctx context.Context, id uuid.UUID, token string) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET confirmation_token = $2, confirmation_sent_at = NOW() WHERE id = $1`, id, token)
	return err
}

func (r *Repository) GetUserByConfirmationToken(ctx context.Context, token string) (*api.User, error) {
	return scanUser(r.pool.QueryRow(ctx,
		`SELECT * FROM users WHERE confirmation_token = $1 AND deleted_at IS NULL`, token))
}

func (r *Repository) SetUserRecoveryToken(ctx context.Context, id uuid.UUID, token string) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET recovery_token = $2, recovery_sent_at = NOW() WHERE id = $1`, id, token)
	return err
}

func (r *Repository) GetUserByRecoveryToken(ctx context.Context, token string) (*api.User, error) {
	return scanUser(r.pool.QueryRow(ctx,
		`SELECT * FROM users WHERE recovery_token = $1 AND deleted_at IS NULL`, token))
}

// --- Sessions ---

func (r *Repository) CreateSession(ctx context.Context, userID uuid.UUID, ip *string, userAgent *string) (*api.Session, error) {
	var s api.Session
	err := r.pool.QueryRow(ctx,
		`INSERT INTO sessions (user_id, ip, user_agent) VALUES ($1, $2::inet, $3) RETURNING id, user_id, ip::text, user_agent, created_at, updated_at, not_after`,
		userID, ip, userAgent,
	).Scan(&s.ID, &s.UserID, &s.IP, &s.UserAgent, &s.CreatedAt, nil, &s.NotAfter)
	return &s, err
}

func (r *Repository) DeleteUserSessions(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE user_id = $1`, userID)
	return err
}

// --- Refresh Tokens ---

func (r *Repository) CreateRefreshToken(ctx context.Context, sessionID *uuid.UUID, userID uuid.UUID, token string, parent *string) (*api.RefreshToken, error) {
	var rt api.RefreshToken
	err := r.pool.QueryRow(ctx,
		`INSERT INTO refresh_tokens (session_id, user_id, token, parent) VALUES ($1, $2, $3, $4) RETURNING id, session_id, user_id, token, parent, revoked, created_at`,
		sessionID, userID, token, parent,
	).Scan(&rt.ID, &rt.SessionID, &rt.UserID, &rt.Token, &rt.Parent, &rt.Revoked, &rt.CreatedAt)
	return &rt, err
}

func (r *Repository) GetRefreshTokenByToken(ctx context.Context, token string) (*api.RefreshToken, error) {
	var rt api.RefreshToken
	err := r.pool.QueryRow(ctx,
		`SELECT id, session_id, user_id, token, parent, revoked, created_at FROM refresh_tokens WHERE token = $1`, token,
	).Scan(&rt.ID, &rt.SessionID, &rt.UserID, &rt.Token, &rt.Parent, &rt.Revoked, &rt.CreatedAt)
	return &rt, err
}

func (r *Repository) RevokeRefreshToken(ctx context.Context, token string) error {
	_, err := r.pool.Exec(ctx, `UPDATE refresh_tokens SET revoked = TRUE WHERE token = $1`, token)
	return err
}

func (r *Repository) RevokeRefreshTokenFamily(ctx context.Context, sessionID *uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE refresh_tokens SET revoked = TRUE WHERE session_id = $1`, sessionID)
	return err
}

// --- Workspaces ---

func (r *Repository) CreateWorkspace(ctx context.Context, name, slug string, ownerID uuid.UUID, organizationSize *string) (*api.Workspace, error) {
	var w api.Workspace
	err := r.pool.QueryRow(ctx,
		`INSERT INTO workspaces (name, slug, owner_id, organization_size) VALUES ($1, $2, $3, $4)
		 RETURNING id, name, slug, logo_url, owner_id, organization_size, created_at, updated_at`,
		name, slug, ownerID, organizationSize,
	).Scan(&w.ID, &w.Name, &w.Slug, &w.LogoURL, &w.OwnerID, &w.OrganizationSize, &w.CreatedAt, &w.UpdatedAt)
	return &w, err
}

func (r *Repository) GetWorkspaceByID(ctx context.Context, id uuid.UUID) (*api.Workspace, error) {
	var w api.Workspace
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, slug, logo_url, owner_id, organization_size, created_at, updated_at FROM workspaces WHERE id = $1 AND deleted_at IS NULL`, id,
	).Scan(&w.ID, &w.Name, &w.Slug, &w.LogoURL, &w.OwnerID, &w.OrganizationSize, &w.CreatedAt, &w.UpdatedAt)
	return &w, err
}

func (r *Repository) GetWorkspaceBySlug(ctx context.Context, slug string) (*api.Workspace, error) {
	var w api.Workspace
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, slug, logo_url, owner_id, organization_size, created_at, updated_at FROM workspaces WHERE slug = $1 AND deleted_at IS NULL`, slug,
	).Scan(&w.ID, &w.Name, &w.Slug, &w.LogoURL, &w.OwnerID, &w.OrganizationSize, &w.CreatedAt, &w.UpdatedAt)
	return &w, err
}

func (r *Repository) ListWorkspacesByUser(ctx context.Context, userID uuid.UUID) ([]*api.Workspace, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT w.id, w.name, w.slug, w.logo_url, w.owner_id, w.organization_size, w.created_at, w.updated_at
		 FROM workspaces w JOIN workspace_members wm ON w.id = wm.workspace_id
		 WHERE wm.user_id = $1 AND wm.deleted_at IS NULL AND w.deleted_at IS NULL
		 ORDER BY w.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []*api.Workspace
	for rows.Next() {
		var w api.Workspace
		if err := rows.Scan(&w.ID, &w.Name, &w.Slug, &w.LogoURL, &w.OwnerID, &w.OrganizationSize, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		workspaces = append(workspaces, &w)
	}
	if workspaces == nil {
		workspaces = []*api.Workspace{}
	}
	return workspaces, nil
}

func (r *Repository) UpdateWorkspace(ctx context.Context, id uuid.UUID, name, slug, logoURL, organizationSize *string) (*api.Workspace, error) {
	var w api.Workspace
	err := r.pool.QueryRow(ctx,
		`UPDATE workspaces SET
			name = COALESCE($2, name), slug = COALESCE($3, slug),
			logo_url = COALESCE($4, logo_url), organization_size = COALESCE($5, organization_size)
		 WHERE id = $1 AND deleted_at IS NULL
		 RETURNING id, name, slug, logo_url, owner_id, organization_size, created_at, updated_at`,
		id, name, slug, logoURL, organizationSize,
	).Scan(&w.ID, &w.Name, &w.Slug, &w.LogoURL, &w.OwnerID, &w.OrganizationSize, &w.CreatedAt, &w.UpdatedAt)
	return &w, err
}

func (r *Repository) DeleteWorkspace(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE workspaces SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}

// --- Workspace Members ---

func (r *Repository) CreateWorkspaceMember(ctx context.Context, workspaceID, userID uuid.UUID, role int16) (*api.WorkspaceMember, error) {
	var m api.WorkspaceMember
	err := r.pool.QueryRow(ctx,
		`INSERT INTO workspace_members (workspace_id, user_id, role) VALUES ($1, $2, $3)
		 RETURNING id, workspace_id, user_id, role, is_active, created_at`,
		workspaceID, userID, role,
	).Scan(&m.ID, &m.WorkspaceID, &m.UserID, &m.Role, &m.IsActive, &m.CreatedAt)
	return &m, err
}

func (r *Repository) GetWorkspaceMember(ctx context.Context, workspaceID, userID uuid.UUID) (*api.WorkspaceMember, error) {
	var m api.WorkspaceMember
	err := r.pool.QueryRow(ctx,
		`SELECT id, workspace_id, user_id, role, is_active, created_at FROM workspace_members
		 WHERE workspace_id = $1 AND user_id = $2 AND deleted_at IS NULL`,
		workspaceID, userID,
	).Scan(&m.ID, &m.WorkspaceID, &m.UserID, &m.Role, &m.IsActive, &m.CreatedAt)
	return &m, err
}

func (r *Repository) ListWorkspaceMembers(ctx context.Context, workspaceID uuid.UUID) ([]*api.WorkspaceMemberWithUser, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT wm.id, wm.workspace_id, wm.user_id, wm.role, wm.is_active, wm.created_at,
		        u.email, u.first_name, u.last_name, u.display_name, u.avatar_url
		 FROM workspace_members wm JOIN users u ON wm.user_id = u.id
		 WHERE wm.workspace_id = $1 AND wm.deleted_at IS NULL
		 ORDER BY wm.created_at ASC`, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*api.WorkspaceMemberWithUser
	for rows.Next() {
		var m api.WorkspaceMemberWithUser
		if err := rows.Scan(
			&m.ID, &m.WorkspaceID, &m.UserID, &m.Role, &m.IsActive, &m.CreatedAt,
			&m.Email, &m.FirstName, &m.LastName, &m.DisplayName, &m.AvatarURL,
		); err != nil {
			return nil, err
		}
		members = append(members, &m)
	}
	if members == nil {
		members = []*api.WorkspaceMemberWithUser{}
	}
	return members, nil
}

func (r *Repository) UpdateWorkspaceMemberRole(ctx context.Context, workspaceID, userID uuid.UUID, role int16) (*api.WorkspaceMember, error) {
	var m api.WorkspaceMember
	err := r.pool.QueryRow(ctx,
		`UPDATE workspace_members SET role = $3
		 WHERE workspace_id = $1 AND user_id = $2 AND deleted_at IS NULL
		 RETURNING id, workspace_id, user_id, role, is_active, created_at`,
		workspaceID, userID, role,
	).Scan(&m.ID, &m.WorkspaceID, &m.UserID, &m.Role, &m.IsActive, &m.CreatedAt)
	return &m, err
}

func (r *Repository) DeleteWorkspaceMember(ctx context.Context, workspaceID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE workspace_members SET deleted_at = NOW() WHERE workspace_id = $1 AND user_id = $2`, workspaceID, userID)
	return err
}

// --- Workspace Invites ---

func (r *Repository) CreateWorkspaceInvite(ctx context.Context, workspaceID uuid.UUID, email string, role int16, token string, message *string) (*api.WorkspaceInvite, error) {
	var inv api.WorkspaceInvite
	err := r.pool.QueryRow(ctx,
		`INSERT INTO workspace_member_invites (workspace_id, email, role, token, message) VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, workspace_id, email, role, token, message, accepted, responded_at, created_at`,
		workspaceID, email, role, token, message,
	).Scan(&inv.ID, &inv.WorkspaceID, &inv.Email, &inv.Role, &inv.Token, &inv.Message, &inv.Accepted, &inv.RespondedAt, &inv.CreatedAt)
	return &inv, err
}

func (r *Repository) GetWorkspaceInviteByToken(ctx context.Context, token string) (*api.WorkspaceInvite, error) {
	var inv api.WorkspaceInvite
	err := r.pool.QueryRow(ctx,
		`SELECT id, workspace_id, email, role, token, message, accepted, responded_at, created_at
		 FROM workspace_member_invites WHERE token = $1 AND deleted_at IS NULL`, token,
	).Scan(&inv.ID, &inv.WorkspaceID, &inv.Email, &inv.Role, &inv.Token, &inv.Message, &inv.Accepted, &inv.RespondedAt, &inv.CreatedAt)
	return &inv, err
}

func (r *Repository) ListWorkspaceInvites(ctx context.Context, workspaceID uuid.UUID) ([]*api.WorkspaceInvite, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, workspace_id, email, role, token, message, accepted, responded_at, created_at
		 FROM workspace_member_invites WHERE workspace_id = $1 AND deleted_at IS NULL AND accepted = FALSE
		 ORDER BY created_at DESC`, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invites []*api.WorkspaceInvite
	for rows.Next() {
		var inv api.WorkspaceInvite
		if err := rows.Scan(&inv.ID, &inv.WorkspaceID, &inv.Email, &inv.Role, &inv.Token, &inv.Message, &inv.Accepted, &inv.RespondedAt, &inv.CreatedAt); err != nil {
			return nil, err
		}
		invites = append(invites, &inv)
	}
	if invites == nil {
		invites = []*api.WorkspaceInvite{}
	}
	return invites, nil
}

func (r *Repository) AcceptWorkspaceInvite(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE workspace_member_invites SET accepted = TRUE, responded_at = NOW() WHERE id = $1`, id)
	return err
}

// --- Projects ---

func (r *Repository) CreateProject(ctx context.Context, workspaceID uuid.UUID, name, description, identifier string, createdBy, projectLeadID *uuid.UUID) (*api.Project, error) {
	var p api.Project
	err := r.pool.QueryRow(ctx,
		`INSERT INTO projects (workspace_id, name, description, identifier, created_by, project_lead_id)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, workspace_id, name, description, identifier, network, emoji, cover_image_url,
		           default_assignee_id, project_lead_id, cycle_view, module_view, page_view, inbox_view,
		           sort_order, archived_at, created_by, created_at, updated_at`,
		workspaceID, name, description, identifier, createdBy, projectLeadID,
	).Scan(&p.ID, &p.WorkspaceID, &p.Name, &p.Description, &p.Identifier, &p.Network,
		&p.Emoji, &p.CoverImageURL, &p.DefaultAssigneeID, &p.ProjectLeadID,
		&p.CycleView, &p.ModuleView, &p.PageView, &p.InboxView,
		&p.SortOrder, &p.ArchivedAt, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	return &p, err
}

func (r *Repository) GetProjectByID(ctx context.Context, id uuid.UUID) (*api.Project, error) {
	var p api.Project
	err := r.pool.QueryRow(ctx,
		`SELECT id, workspace_id, name, description, identifier, network, emoji, cover_image_url,
		        default_assignee_id, project_lead_id, cycle_view, module_view, page_view, inbox_view,
		        sort_order, archived_at, created_by, created_at, updated_at
		 FROM projects WHERE id = $1 AND deleted_at IS NULL`, id,
	).Scan(&p.ID, &p.WorkspaceID, &p.Name, &p.Description, &p.Identifier, &p.Network,
		&p.Emoji, &p.CoverImageURL, &p.DefaultAssigneeID, &p.ProjectLeadID,
		&p.CycleView, &p.ModuleView, &p.PageView, &p.InboxView,
		&p.SortOrder, &p.ArchivedAt, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	return &p, err
}

func (r *Repository) ListProjectsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]*api.Project, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, workspace_id, name, description, identifier, network, emoji, cover_image_url,
		        default_assignee_id, project_lead_id, cycle_view, module_view, page_view, inbox_view,
		        sort_order, archived_at, created_by, created_at, updated_at
		 FROM projects WHERE workspace_id = $1 AND deleted_at IS NULL
		 ORDER BY sort_order ASC, created_at DESC`, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*api.Project
	for rows.Next() {
		var p api.Project
		if err := rows.Scan(&p.ID, &p.WorkspaceID, &p.Name, &p.Description, &p.Identifier, &p.Network,
			&p.Emoji, &p.CoverImageURL, &p.DefaultAssigneeID, &p.ProjectLeadID,
			&p.CycleView, &p.ModuleView, &p.PageView, &p.InboxView,
			&p.SortOrder, &p.ArchivedAt, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, &p)
	}
	if projects == nil {
		projects = []*api.Project{}
	}
	return projects, nil
}

func (r *Repository) UpdateProject(ctx context.Context, id uuid.UUID, name, description, identifier, emoji, coverImageURL *string, defaultAssigneeID, projectLeadID *uuid.UUID, network *int16, cycleView, moduleView, pageView, inboxView *bool) (*api.Project, error) {
	var p api.Project
	err := r.pool.QueryRow(ctx,
		`UPDATE projects SET
			name = COALESCE($2, name), description = COALESCE($3, description),
			identifier = COALESCE($4, identifier), emoji = COALESCE($5, emoji),
			cover_image_url = COALESCE($6, cover_image_url),
			default_assignee_id = $7, project_lead_id = $8,
			network = COALESCE($9, network), cycle_view = COALESCE($10, cycle_view),
			module_view = COALESCE($11, module_view), page_view = COALESCE($12, page_view),
			inbox_view = COALESCE($13, inbox_view)
		 WHERE id = $1 AND deleted_at IS NULL
		 RETURNING id, workspace_id, name, description, identifier, network, emoji, cover_image_url,
		           default_assignee_id, project_lead_id, cycle_view, module_view, page_view, inbox_view,
		           sort_order, archived_at, created_by, created_at, updated_at`,
		id, name, description, identifier, emoji, coverImageURL,
		defaultAssigneeID, projectLeadID, network, cycleView, moduleView, pageView, inboxView,
	).Scan(&p.ID, &p.WorkspaceID, &p.Name, &p.Description, &p.Identifier, &p.Network,
		&p.Emoji, &p.CoverImageURL, &p.DefaultAssigneeID, &p.ProjectLeadID,
		&p.CycleView, &p.ModuleView, &p.PageView, &p.InboxView,
		&p.SortOrder, &p.ArchivedAt, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	return &p, err
}

func (r *Repository) DeleteProject(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE projects SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}

// --- Project Members ---

func (r *Repository) CreateProjectMember(ctx context.Context, projectID, userID uuid.UUID, role int16) (*api.ProjectMember, error) {
	var m api.ProjectMember
	err := r.pool.QueryRow(ctx,
		`INSERT INTO project_members (project_id, user_id, role) VALUES ($1, $2, $3)
		 RETURNING id, project_id, user_id, role, is_active, created_at`,
		projectID, userID, role,
	).Scan(&m.ID, &m.ProjectID, &m.UserID, &m.Role, &m.IsActive, &m.CreatedAt)
	return &m, err
}

func (r *Repository) GetProjectMember(ctx context.Context, projectID, userID uuid.UUID) (*api.ProjectMember, error) {
	var m api.ProjectMember
	err := r.pool.QueryRow(ctx,
		`SELECT id, project_id, user_id, role, is_active, created_at FROM project_members
		 WHERE project_id = $1 AND user_id = $2 AND deleted_at IS NULL`,
		projectID, userID,
	).Scan(&m.ID, &m.ProjectID, &m.UserID, &m.Role, &m.IsActive, &m.CreatedAt)
	return &m, err
}

func (r *Repository) ListProjectMembers(ctx context.Context, projectID uuid.UUID) ([]*api.ProjectMemberWithUser, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT pm.id, pm.project_id, pm.user_id, pm.role, pm.is_active, pm.created_at,
		        u.email, u.first_name, u.last_name, u.display_name, u.avatar_url
		 FROM project_members pm JOIN users u ON pm.user_id = u.id
		 WHERE pm.project_id = $1 AND pm.deleted_at IS NULL
		 ORDER BY pm.created_at ASC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*api.ProjectMemberWithUser
	for rows.Next() {
		var m api.ProjectMemberWithUser
		if err := rows.Scan(&m.ID, &m.ProjectID, &m.UserID, &m.Role, &m.IsActive, &m.CreatedAt,
			&m.Email, &m.FirstName, &m.LastName, &m.DisplayName, &m.AvatarURL); err != nil {
			return nil, err
		}
		members = append(members, &m)
	}
	if members == nil {
		members = []*api.ProjectMemberWithUser{}
	}
	return members, nil
}

func (r *Repository) UpdateProjectMemberRole(ctx context.Context, projectID, userID uuid.UUID, role int16) (*api.ProjectMember, error) {
	var m api.ProjectMember
	err := r.pool.QueryRow(ctx,
		`UPDATE project_members SET role = $3 WHERE project_id = $1 AND user_id = $2 AND deleted_at IS NULL
		 RETURNING id, project_id, user_id, role, is_active, created_at`,
		projectID, userID, role,
	).Scan(&m.ID, &m.ProjectID, &m.UserID, &m.Role, &m.IsActive, &m.CreatedAt)
	return &m, err
}

func (r *Repository) DeleteProjectMember(ctx context.Context, projectID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE project_members SET deleted_at = NOW() WHERE project_id = $1 AND user_id = $2`, projectID, userID)
	return err
}
