package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

// loggingMiddleware logs each request.
func (a *API) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Dur("duration", time.Since(start)).
			Msg("request")
	})
}

// requireAuthentication extracts and validates the JWT from the Authorization header.
func (a *API) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			handleError(w, r, unauthorizedError("missing authorization header"))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			handleError(w, r, unauthorizedError("invalid authorization format"))
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, unauthorizedError("unexpected signing method")
			}
			return []byte(a.config.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			handleError(w, r, unauthorizedError("invalid or expired token"))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			handleError(w, r, unauthorizedError("invalid token claims"))
			return
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			handleError(w, r, unauthorizedError("missing subject in token"))
			return
		}

		userID, err := uuid.FromString(sub)
		if err != nil {
			handleError(w, r, unauthorizedError("invalid user ID in token"))
			return
		}

		ctx := withUserID(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// requireWorkspaceAccess loads workspace from slug and verifies membership.
func (a *API) requireWorkspaceAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := getUserID(ctx)
		slug := chi.URLParam(r, "workspaceSlug")

		workspace, err := a.queries.GetWorkspaceBySlug(ctx, slug)
		if err != nil {
			handleError(w, r, notFoundError("workspace not found"))
			return
		}

		member, err := a.queries.GetWorkspaceMember(ctx, workspace.ID, userID)
		if err != nil {
			handleError(w, r, forbiddenError("not a member of this workspace"))
			return
		}

		ctx = withWorkspaceID(ctx, workspace.ID)
		ctx = withWorkspaceSlug(ctx, slug)
		ctx = withWorkspaceRole(ctx, member.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// requireProjectAccess loads project and verifies membership.
func (a *API) requireProjectAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := getUserID(ctx)
		projectIDStr := chi.URLParam(r, "projectID")

		projectID, err := uuid.FromString(projectIDStr)
		if err != nil {
			handleError(w, r, badRequestError("invalid project ID"))
			return
		}

		project, err := a.queries.GetProjectByID(ctx, projectID)
		if err != nil {
			handleError(w, r, notFoundError("project not found"))
			return
		}

		// Check project membership (or workspace admin gets access)
		wsRole := getWorkspaceRole(ctx)
		if wsRole < RoleAdmin {
			_, err := a.queries.GetProjectMember(ctx, project.ID, userID)
			if err != nil {
				handleError(w, r, forbiddenError("not a member of this project"))
				return
			}
		}

		ctx = withProjectID(ctx, project.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// requireWorkspaceAdmin ensures the user has admin or owner role.
func requireWorkspaceAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := getWorkspaceRole(r.Context())
		if role < RoleAdmin {
			handleError(w, r, forbiddenError("admin access required"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
