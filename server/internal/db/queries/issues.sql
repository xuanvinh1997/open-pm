-- name: CreateIssue :one
INSERT INTO issues (
    project_id, workspace_id, parent_id, state_id,
    name, description_html, description_json, description_stripped,
    priority, start_date, target_date, is_draft, created_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetIssueByID :one
SELECT * FROM issues WHERE id = $1 AND deleted_at IS NULL;

-- name: ListIssuesByProject :many
SELECT * FROM issues
WHERE project_id = $1 AND deleted_at IS NULL
ORDER BY sort_order ASC, created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountIssuesByProject :one
SELECT COUNT(*) FROM issues
WHERE project_id = $1 AND deleted_at IS NULL;

-- name: ListIssuesByState :many
SELECT * FROM issues
WHERE project_id = $1 AND state_id = $2 AND deleted_at IS NULL
ORDER BY sort_order ASC;

-- name: ListIssuesByCycle :many
SELECT i.* FROM issues i
JOIN cycle_issues ci ON i.id = ci.issue_id
WHERE ci.cycle_id = $1 AND i.deleted_at IS NULL
ORDER BY i.sort_order ASC;

-- name: ListIssuesByModule :many
SELECT i.* FROM issues i
JOIN module_issues mi ON i.id = mi.issue_id
WHERE mi.module_id = $1 AND i.deleted_at IS NULL
ORDER BY i.sort_order ASC;

-- name: ListSubIssues :many
SELECT * FROM issues
WHERE parent_id = $1 AND deleted_at IS NULL
ORDER BY sort_order ASC;

-- name: UpdateIssue :one
UPDATE issues SET
    state_id = COALESCE(sqlc.narg('state_id'), state_id),
    name = COALESCE(sqlc.narg('name'), name),
    description_html = COALESCE(sqlc.narg('description_html'), description_html),
    description_json = COALESCE(sqlc.narg('description_json'), description_json),
    description_stripped = COALESCE(sqlc.narg('description_stripped'), description_stripped),
    priority = COALESCE(sqlc.narg('priority'), priority),
    start_date = sqlc.narg('start_date'),
    target_date = sqlc.narg('target_date'),
    parent_id = sqlc.narg('parent_id'),
    sort_order = COALESCE(sqlc.narg('sort_order'), sort_order),
    is_draft = COALESCE(sqlc.narg('is_draft'), is_draft),
    completed_at = sqlc.narg('completed_at'),
    archived_at = sqlc.narg('archived_at'),
    updated_by = sqlc.narg('updated_by')
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteIssue :exec
UPDATE issues SET deleted_at = NOW() WHERE id = $1;

-- name: AddIssueAssignee :exec
INSERT INTO issue_assignees (issue_id, assignee_id)
VALUES ($1, $2)
ON CONFLICT (issue_id, assignee_id) DO NOTHING;

-- name: RemoveIssueAssignee :exec
DELETE FROM issue_assignees WHERE issue_id = $1 AND assignee_id = $2;

-- name: ListIssueAssignees :many
SELECT u.id, u.email, u.first_name, u.last_name, u.display_name, u.avatar_url
FROM issue_assignees ia
JOIN users u ON ia.assignee_id = u.id
WHERE ia.issue_id = $1;

-- name: AddIssueLabel :exec
INSERT INTO issue_labels (issue_id, label_id)
VALUES ($1, $2)
ON CONFLICT (issue_id, label_id) DO NOTHING;

-- name: RemoveIssueLabel :exec
DELETE FROM issue_labels WHERE issue_id = $1 AND label_id = $2;

-- name: ListIssueLabels :many
SELECT l.* FROM issue_labels il
JOIN labels l ON il.label_id = l.id
WHERE il.issue_id = $1 AND l.deleted_at IS NULL;

-- name: CreateIssueComment :one
INSERT INTO issue_comments (issue_id, workspace_id, comment_html, comment_json, comment_stripped, actor_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetIssueCommentByID :one
SELECT * FROM issue_comments WHERE id = $1 AND deleted_at IS NULL;

-- name: ListIssueComments :many
SELECT ic.*, u.email, u.first_name, u.last_name, u.display_name, u.avatar_url
FROM issue_comments ic
LEFT JOIN users u ON ic.actor_id = u.id
WHERE ic.issue_id = $1 AND ic.deleted_at IS NULL
ORDER BY ic.created_at ASC;

-- name: UpdateIssueComment :one
UPDATE issue_comments SET
    comment_html = $2,
    comment_json = $3,
    comment_stripped = $4
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteIssueComment :exec
UPDATE issue_comments SET deleted_at = NOW() WHERE id = $1;

-- name: CreateIssueActivity :one
INSERT INTO issue_activities (
    issue_id, project_id, workspace_id,
    verb, field, old_value, new_value,
    old_identifier, new_identifier, comment, actor_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: ListIssueActivities :many
SELECT ia.*, u.email, u.first_name, u.last_name, u.display_name, u.avatar_url
FROM issue_activities ia
LEFT JOIN users u ON ia.actor_id = u.id
WHERE ia.issue_id = $1
ORDER BY ia.created_at ASC;
