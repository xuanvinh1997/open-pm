-- name: CreateNotification :one
INSERT INTO notifications (workspace_id, project_id, title, message, data, entity_type, entity_id, sender_id, receiver_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: ListNotificationsByUser :many
SELECT * FROM notifications
WHERE receiver_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountUnreadNotifications :one
SELECT COUNT(*) FROM notifications
WHERE receiver_id = $1 AND read_at IS NULL;

-- name: MarkNotificationRead :exec
UPDATE notifications SET read_at = NOW() WHERE id = $1 AND receiver_id = $2;

-- name: MarkAllNotificationsRead :exec
UPDATE notifications SET read_at = NOW() WHERE receiver_id = $1 AND read_at IS NULL;
