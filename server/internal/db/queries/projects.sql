-- name: CreateProject :one
INSERT INTO projects (workspace_id, name, description, identifier, created_by, project_lead_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetProjectByID :one
SELECT * FROM projects WHERE id = $1 AND deleted_at IS NULL;

-- name: ListProjectsByWorkspace :many
SELECT * FROM projects
WHERE workspace_id = $1 AND deleted_at IS NULL
ORDER BY sort_order ASC, created_at DESC;

-- name: UpdateProject :one
UPDATE projects SET
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    identifier = COALESCE(sqlc.narg('identifier'), identifier),
    emoji = COALESCE(sqlc.narg('emoji'), emoji),
    cover_image_url = COALESCE(sqlc.narg('cover_image_url'), cover_image_url),
    default_assignee_id = sqlc.narg('default_assignee_id'),
    project_lead_id = sqlc.narg('project_lead_id'),
    network = COALESCE(sqlc.narg('network'), network),
    sprint_view = COALESCE(sqlc.narg('sprint_view'), sprint_view),
    epic_view = COALESCE(sqlc.narg('epic_view'), epic_view),
    page_view = COALESCE(sqlc.narg('page_view'), page_view),
    inbox_view = COALESCE(sqlc.narg('inbox_view'), inbox_view)
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteProject :exec
UPDATE projects SET deleted_at = NOW() WHERE id = $1;

-- name: CreateProjectMember :one
INSERT INTO project_members (project_id, user_id, role)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetProjectMember :one
SELECT * FROM project_members
WHERE project_id = $1 AND user_id = $2 AND deleted_at IS NULL;

-- name: ListProjectMembers :many
SELECT pm.*, u.email, u.first_name, u.last_name, u.display_name, u.avatar_url
FROM project_members pm
JOIN users u ON pm.user_id = u.id
WHERE pm.project_id = $1 AND pm.deleted_at IS NULL
ORDER BY pm.created_at ASC;

-- name: UpdateProjectMemberRole :one
UPDATE project_members SET role = $3
WHERE project_id = $1 AND user_id = $2 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteProjectMember :exec
UPDATE project_members SET deleted_at = NOW()
WHERE project_id = $1 AND user_id = $2;
