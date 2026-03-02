-- name: CreateSprint :one
INSERT INTO sprints (project_id, workspace_id, name, description, start_date, end_date, owned_by)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetSprintByID :one
SELECT * FROM sprints WHERE id = $1 AND deleted_at IS NULL;

-- name: ListSprintsByProject :many
SELECT * FROM sprints
WHERE project_id = $1 AND deleted_at IS NULL
ORDER BY start_date DESC NULLS LAST;

-- name: UpdateSprint :one
UPDATE sprints SET
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    start_date = sqlc.narg('start_date'),
    end_date = sqlc.narg('end_date'),
    sort_order = COALESCE(sqlc.narg('sort_order'), sort_order)
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteSprint :exec
UPDATE sprints SET deleted_at = NOW() WHERE id = $1;

-- name: AddIssueToSprint :exec
INSERT INTO sprint_issues (sprint_id, issue_id)
VALUES ($1, $2)
ON CONFLICT (sprint_id, issue_id) DO NOTHING;

-- name: RemoveIssueFromSprint :exec
DELETE FROM sprint_issues WHERE sprint_id = $1 AND issue_id = $2;

-- name: ListSprintIssueIDs :many
SELECT issue_id FROM sprint_issues WHERE sprint_id = $1;
