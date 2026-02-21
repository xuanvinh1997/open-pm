package api

import (
	"context"

	"github.com/gofrs/uuid"
)

type contextKey string

const (
	ctxUserIDKey        contextKey = "user_id"
	ctxWorkspaceIDKey   contextKey = "workspace_id"
	ctxWorkspaceSlugKey contextKey = "workspace_slug"
	ctxWorkspaceRoleKey contextKey = "workspace_role"
	ctxProjectIDKey     contextKey = "project_id"
	ctxProjectRoleKey   contextKey = "project_role"
)

// Role constants
const (
	RoleOwner  int16 = 25
	RoleAdmin  int16 = 20
	RoleMember int16 = 15
	RoleGuest  int16 = 5
)

func withUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxUserIDKey, userID)
}

func getUserID(ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value(ctxUserIDKey).(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}

func withWorkspaceID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxWorkspaceIDKey, id)
}

func getWorkspaceID(ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value(ctxWorkspaceIDKey).(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}

func withWorkspaceSlug(ctx context.Context, slug string) context.Context {
	return context.WithValue(ctx, ctxWorkspaceSlugKey, slug)
}

func getWorkspaceSlug(ctx context.Context) string {
	if s, ok := ctx.Value(ctxWorkspaceSlugKey).(string); ok {
		return s
	}
	return ""
}

func withWorkspaceRole(ctx context.Context, role int16) context.Context {
	return context.WithValue(ctx, ctxWorkspaceRoleKey, role)
}

func getWorkspaceRole(ctx context.Context) int16 {
	if r, ok := ctx.Value(ctxWorkspaceRoleKey).(int16); ok {
		return r
	}
	return 0
}

func withProjectID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxProjectIDKey, id)
}

func getProjectID(ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value(ctxProjectIDKey).(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}
