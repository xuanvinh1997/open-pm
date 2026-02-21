-- name: CreateWorkspace :one
INSERT INTO workspaces (name, slug, owner_id, organization_size)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetWorkspaceByID :one
SELECT * FROM workspaces WHERE id = $1 AND deleted_at IS NULL;

-- name: GetWorkspaceBySlug :one
SELECT * FROM workspaces WHERE slug = $1 AND deleted_at IS NULL;

-- name: ListWorkspacesByUser :many
SELECT w.* FROM workspaces w
JOIN workspace_members wm ON w.id = wm.workspace_id
WHERE wm.user_id = $1 AND wm.deleted_at IS NULL AND w.deleted_at IS NULL
ORDER BY w.created_at DESC;

-- name: UpdateWorkspace :one
UPDATE workspaces SET
    name = COALESCE(sqlc.narg('name'), name),
    slug = COALESCE(sqlc.narg('slug'), slug),
    logo_url = COALESCE(sqlc.narg('logo_url'), logo_url),
    organization_size = COALESCE(sqlc.narg('organization_size'), organization_size)
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteWorkspace :exec
UPDATE workspaces SET deleted_at = NOW() WHERE id = $1;

-- name: CreateWorkspaceMember :one
INSERT INTO workspace_members (workspace_id, user_id, role)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetWorkspaceMember :one
SELECT * FROM workspace_members
WHERE workspace_id = $1 AND user_id = $2 AND deleted_at IS NULL;

-- name: ListWorkspaceMembers :many
SELECT wm.*, u.email, u.first_name, u.last_name, u.display_name, u.avatar_url
FROM workspace_members wm
JOIN users u ON wm.user_id = u.id
WHERE wm.workspace_id = $1 AND wm.deleted_at IS NULL
ORDER BY wm.created_at ASC;

-- name: UpdateWorkspaceMemberRole :one
UPDATE workspace_members SET role = $3
WHERE workspace_id = $1 AND user_id = $2 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteWorkspaceMember :exec
UPDATE workspace_members SET deleted_at = NOW()
WHERE workspace_id = $1 AND user_id = $2;

-- name: CreateWorkspaceInvite :one
INSERT INTO workspace_member_invites (workspace_id, email, role, token, message)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetWorkspaceInviteByID :one
SELECT * FROM workspace_member_invites
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetWorkspaceInviteByToken :one
SELECT * FROM workspace_member_invites
WHERE token = $1 AND deleted_at IS NULL;

-- name: ListWorkspaceInvites :many
SELECT * FROM workspace_member_invites
WHERE workspace_id = $1 AND deleted_at IS NULL AND accepted = FALSE
ORDER BY created_at DESC;

-- name: AcceptWorkspaceInvite :exec
UPDATE workspace_member_invites SET accepted = TRUE, responded_at = NOW()
WHERE id = $1;

-- name: DeleteWorkspaceInvite :exec
UPDATE workspace_member_invites SET deleted_at = NOW()
WHERE id = $1;
