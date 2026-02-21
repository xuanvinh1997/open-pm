-- name: CreateState :one
INSERT INTO states (project_id, workspace_id, name, description, color, "group", sequence, is_default)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetStateByID :one
SELECT * FROM states WHERE id = $1 AND deleted_at IS NULL;

-- name: ListStatesByProject :many
SELECT * FROM states
WHERE project_id = $1 AND deleted_at IS NULL
ORDER BY sequence ASC;

-- name: UpdateState :one
UPDATE states SET
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    color = COALESCE(sqlc.narg('color'), color),
    "group" = COALESCE(sqlc.narg('group'), "group"),
    sequence = COALESCE(sqlc.narg('sequence'), sequence),
    is_default = COALESCE(sqlc.narg('is_default'), is_default)
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteState :exec
UPDATE states SET deleted_at = NOW() WHERE id = $1;

-- name: GetDefaultState :one
SELECT * FROM states
WHERE project_id = $1 AND is_default = TRUE AND deleted_at IS NULL
LIMIT 1;
