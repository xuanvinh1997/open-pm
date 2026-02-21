-- name: CreateModule :one
INSERT INTO modules (project_id, workspace_id, name, description, start_date, target_date, status, lead_id, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetModuleByID :one
SELECT * FROM modules WHERE id = $1 AND deleted_at IS NULL;

-- name: ListModulesByProject :many
SELECT * FROM modules
WHERE project_id = $1 AND deleted_at IS NULL
ORDER BY sort_order ASC, created_at DESC;

-- name: UpdateModule :one
UPDATE modules SET
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    start_date = sqlc.narg('start_date'),
    target_date = sqlc.narg('target_date'),
    status = COALESCE(sqlc.narg('status'), status),
    lead_id = sqlc.narg('lead_id'),
    sort_order = COALESCE(sqlc.narg('sort_order'), sort_order)
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteModule :exec
UPDATE modules SET deleted_at = NOW() WHERE id = $1;

-- name: AddIssueToModule :exec
INSERT INTO module_issues (module_id, issue_id)
VALUES ($1, $2)
ON CONFLICT (module_id, issue_id) DO NOTHING;

-- name: RemoveIssueFromModule :exec
DELETE FROM module_issues WHERE module_id = $1 AND issue_id = $2;

-- name: ListModuleIssueIDs :many
SELECT issue_id FROM module_issues WHERE module_id = $1;
