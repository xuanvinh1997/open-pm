package repository

import (
	"context"
	"encoding/json"

	"github.com/gofrs/uuid"
	"github.com/open-pm/open-pm/server/internal/api"
)

func (r *Repository) CreateChatMessage(ctx context.Context, projectID, workspaceID, userID uuid.UUID, role, content string) (*api.ChatMessage, error) {
	var m api.ChatMessage
	err := r.pool.QueryRow(ctx,
		`INSERT INTO chat_messages (project_id, workspace_id, user_id, role, content)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, project_id, workspace_id, user_id, role, content, metadata, created_at`,
		projectID, workspaceID, userID, role, content,
	).Scan(&m.ID, &m.ProjectID, &m.WorkspaceID, &m.UserID, &m.Role, &m.Content, (*json.RawMessage)(&m.Metadata), &m.CreatedAt)
	return &m, err
}

func (r *Repository) ListChatMessages(ctx context.Context, projectID, userID uuid.UUID, limit int) ([]*api.ChatMessage, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, project_id, workspace_id, user_id, role, content, metadata, created_at
		 FROM chat_messages
		 WHERE project_id = $1 AND user_id = $2
		 ORDER BY created_at DESC
		 LIMIT $3`, projectID, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*api.ChatMessage
	for rows.Next() {
		var m api.ChatMessage
		if err := rows.Scan(&m.ID, &m.ProjectID, &m.WorkspaceID, &m.UserID, &m.Role, &m.Content, (*json.RawMessage)(&m.Metadata), &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, &m)
	}
	if messages == nil {
		messages = []*api.ChatMessage{}
	}

	// Reverse to get chronological order (oldest first)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *Repository) DeleteChatHistory(ctx context.Context, projectID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM chat_messages WHERE project_id = $1 AND user_id = $2`,
		projectID, userID)
	return err
}
