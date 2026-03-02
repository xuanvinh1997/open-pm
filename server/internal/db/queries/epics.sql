-- name: CreateEpic :one
INSERT INTO epics (project_id, workspace_id, name, description, start_date, target_date, status, lead_id, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetEpicByID :one
SELECT * FROM epics WHERE id = $1 AND deleted_at IS NULL;

-- name: ListEpicsByProject :many
SELECT * FROM epics
WHERE project_id = $1 AND deleted_at IS NULL
ORDER BY sort_order ASC, created_at DESC;

-- name: UpdateEpic :one
UPDATE epics SET
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    start_date = sqlc.narg('start_date'),
    target_date = sqlc.narg('target_date'),
    status = COALESCE(sqlc.narg('status'), status),
    lead_id = sqlc.narg('lead_id'),
    sort_order = COALESCE(sqlc.narg('sort_order'), sort_order)
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteEpic :exec
UPDATE epics SET deleted_at = NOW() WHERE id = $1;

-- name: AddIssueToEpic :exec
INSERT INTO epic_issues (epic_id, issue_id)
VALUES ($1, $2)
ON CONFLICT (epic_id, issue_id) DO NOTHING;

-- name: RemoveIssueFromEpic :exec
DELETE FROM epic_issues WHERE epic_id = $1 AND issue_id = $2;

-- name: ListEpicIssueIDs :many
SELECT issue_id FROM epic_issues WHERE epic_id = $1;
