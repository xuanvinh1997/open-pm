package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

type notifyParams struct {
	WorkspaceID uuid.UUID
	ProjectID   *uuid.UUID
	Title       string
	Message     string
	Data        []byte
	EntityType  *string
	EntityID    *uuid.UUID
}

// notifyUsers creates a notification for each receiver, skipping the sender.
// Errors are logged but not propagated.
func (a *API) notifyUsers(ctx context.Context, receivers []uuid.UUID, senderID uuid.UUID, params notifyParams) {
	for _, receiverID := range receivers {
		if receiverID == senderID {
			continue
		}
		if _, err := a.queries.CreateNotification(ctx,
			params.WorkspaceID, params.ProjectID,
			params.Title, params.Message,
			params.Data, params.EntityType, params.EntityID,
			&senderID, receiverID,
		); err != nil {
			log.Warn().Err(err).
				Str("receiver_id", receiverID.String()).
				Msg("failed to create notification")
		}
	}
}

// collectUniqueUserIDs merges assignee IDs and subscriber IDs into a deduplicated list.
func collectUniqueUserIDs(assignees []*UserSummary, subscriberIDs []uuid.UUID) []uuid.UUID {
	seen := make(map[uuid.UUID]bool)
	var result []uuid.UUID
	for _, a := range assignees {
		if !seen[a.ID] {
			seen[a.ID] = true
			result = append(result, a.ID)
		}
	}
	for _, id := range subscriberIDs {
		if !seen[id] {
			seen[id] = true
			result = append(result, id)
		}
	}
	return result
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// ListNotifications handles GET /api/v1/workspaces/{workspaceSlug}/notifications
func (a *API) ListNotifications(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)
	workspaceID := getWorkspaceID(ctx)
	pg := parsePagination(r)

	notifications, err := a.queries.ListNotificationsByUser(ctx, userID, workspaceID, pg.PerPage, pg.Offset())
	if err != nil {
		return internalServerError("failed to list notifications")
	}

	count, err := a.queries.CountUnreadNotifications(ctx, userID, workspaceID)
	if err != nil {
		log.Warn().Err(err).Msg("failed to count unread notifications")
	}

	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"results":      notifications,
		"unread_count": count,
	})
}

// CountUnread handles GET /api/v1/workspaces/{workspaceSlug}/notifications/unread-count
func (a *API) CountUnread(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)
	workspaceID := getWorkspaceID(ctx)

	count, err := a.queries.CountUnreadNotifications(ctx, userID, workspaceID)
	if err != nil {
		return internalServerError("failed to count unread notifications")
	}

	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"unread_count": count,
	})
}

// MarkNotificationRead handles POST /api/v1/workspaces/{workspaceSlug}/notifications/{notificationID}/read
func (a *API) MarkNotificationRead(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)

	notificationID, err := uuid.FromString(chi.URLParam(r, "notificationID"))
	if err != nil {
		return badRequestError("invalid notification ID")
	}

	if err := a.queries.MarkNotificationRead(ctx, notificationID, userID); err != nil {
		return internalServerError("failed to mark notification as read")
	}

	return sendEmpty(w, http.StatusNoContent)
}

// MarkAllNotificationsRead handles POST /api/v1/workspaces/{workspaceSlug}/notifications/read-all
func (a *API) MarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)
	workspaceID := getWorkspaceID(ctx)

	if err := a.queries.MarkAllNotificationsRead(ctx, userID, workspaceID); err != nil {
		return internalServerError("failed to mark all notifications as read")
	}

	return sendEmpty(w, http.StatusNoContent)
}
