-- name: CreateLabel :one
INSERT INTO labels (project_id, workspace_id, parent_id, name, description, color, sort_order)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetLabelByID :one
SELECT * FROM labels WHERE id = $1 AND deleted_at IS NULL;

-- name: ListLabelsByProject :many
SELECT * FROM labels
WHERE project_id = $1 AND deleted_at IS NULL
ORDER BY sort_order ASC;

-- name: ListLabelsByWorkspace :many
SELECT * FROM labels
WHERE workspace_id = $1 AND project_id IS NULL AND deleted_at IS NULL
ORDER BY sort_order ASC;

-- name: UpdateLabel :one
UPDATE labels SET
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    color = COALESCE(sqlc.narg('color'), color),
    parent_id = sqlc.narg('parent_id'),
    sort_order = COALESCE(sqlc.narg('sort_order'), sort_order)
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteLabel :exec
UPDATE labels SET deleted_at = NOW() WHERE id = $1;
