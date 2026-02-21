-- name: CreatePage :one
INSERT INTO pages (project_id, workspace_id, name, description_html, description_json, description_stripped, owned_by, parent_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPageByID :one
SELECT * FROM pages WHERE id = $1 AND deleted_at IS NULL;

-- name: ListPagesByProject :many
SELECT * FROM pages
WHERE project_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListPagesByWorkspace :many
SELECT * FROM pages
WHERE workspace_id = $1 AND project_id IS NULL AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdatePage :one
UPDATE pages SET
    name = COALESCE(sqlc.narg('name'), name),
    description_html = COALESCE(sqlc.narg('description_html'), description_html),
    description_json = COALESCE(sqlc.narg('description_json'), description_json),
    description_stripped = COALESCE(sqlc.narg('description_stripped'), description_stripped),
    color = COALESCE(sqlc.narg('color'), color),
    is_locked = COALESCE(sqlc.narg('is_locked'), is_locked),
    parent_id = sqlc.narg('parent_id')
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeletePage :exec
UPDATE pages SET deleted_at = NOW() WHERE id = $1;
