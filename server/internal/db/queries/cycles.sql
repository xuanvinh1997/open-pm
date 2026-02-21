-- name: CreateCycle :one
INSERT INTO cycles (project_id, workspace_id, name, description, start_date, end_date, owned_by)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetCycleByID :one
SELECT * FROM cycles WHERE id = $1 AND deleted_at IS NULL;

-- name: ListCyclesByProject :many
SELECT * FROM cycles
WHERE project_id = $1 AND deleted_at IS NULL
ORDER BY start_date DESC NULLS LAST;

-- name: UpdateCycle :one
UPDATE cycles SET
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    start_date = sqlc.narg('start_date'),
    end_date = sqlc.narg('end_date'),
    sort_order = COALESCE(sqlc.narg('sort_order'), sort_order)
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteCycle :exec
UPDATE cycles SET deleted_at = NOW() WHERE id = $1;

-- name: AddIssueToCycle :exec
INSERT INTO cycle_issues (cycle_id, issue_id)
VALUES ($1, $2)
ON CONFLICT (cycle_id, issue_id) DO NOTHING;

-- name: RemoveIssueFromCycle :exec
DELETE FROM cycle_issues WHERE cycle_id = $1 AND issue_id = $2;

-- name: ListCycleIssueIDs :many
SELECT issue_id FROM cycle_issues WHERE cycle_id = $1;
