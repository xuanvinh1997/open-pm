package api

import (
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/go-playground/validator/v10"
	"github.com/open-pm/open-pm/server/internal/config"
	"github.com/open-pm/open-pm/server/internal/service"
	"github.com/open-pm/open-pm/server/internal/email"
	"github.com/open-pm/open-pm/server/internal/storage"
)

type API struct {
	config    *config.Config
	queries   Queries
	validator *validator.Validate
	handler   http.Handler
	llm       service.LLMProvider
	storage   *storage.Storage
	mailer    *email.Mailer
}

// WithLLM sets the LLM provider on the API.
func WithLLM(llm service.LLMProvider) func(*API) {
	return func(a *API) {
		a.llm = llm
	}
}



func NewAPI(cfg *config.Config, queries Queries, opts ...func(*API)) *API {
	a := &API{
		config:    cfg,
		queries:   queries,
		validator: validator.New(),
	}
	for _, opt := range opts {
		opt(a)
	}


	r := newRouter()

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(a.loggingMiddleware)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.SiteURL, "http://localhost:3000", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Rate limiting
	if cfg.RateLimit.Enabled {
		r.Use(httprate.LimitByIP(100, time.Minute))
	}

	// Health check
	r.Get("/health", a.healthCheck)

	// Auth routes (public)
	r.Route("/auth", func(r *router) {
		r.Post("/signup", a.Signup)
		r.Post("/token", a.Token)
		r.Post("/recover", a.Recover)
		r.Get("/verify", a.Verify)

		// Authenticated auth routes
		r.Group(func(r *router) {
			r.Use(a.requireAuthentication)
			r.Get("/user", a.GetUser)
			r.Put("/user", a.UpdateUser)
			r.Post("/logout", a.Logout)
		})
	})

	// API v1 routes (all authenticated)
	r.Route("/api/v1", func(r *router) {
		r.Use(a.requireAuthentication)

		// Workspaces
		r.Get("/workspaces", a.ListWorkspaces)
		r.Post("/workspaces", a.CreateWorkspace)

		r.Route("/workspaces/{workspaceSlug}", func(r *router) {
			r.Use(a.requireWorkspaceAccess)
			r.Get("/", a.GetWorkspace)
			r.Put("/", a.UpdateWorkspace)
			r.Delete("/", a.DeleteWorkspace)

			// Members
			r.Get("/members", a.ListWorkspaceMembers)
			r.Put("/members/{memberID}", a.UpdateWorkspaceMember)
			r.Delete("/members/{memberID}", a.RemoveWorkspaceMember)

			// Invites
			r.Get("/invites", a.ListWorkspaceInvites)
			r.Post("/invites", a.CreateWorkspaceInvite)

			// Search
			r.Get("/search", a.GlobalSearch)

			// Assets (file attachments)
			r.Post("/assets", a.UploadAsset)
			r.Get("/assets", a.ListAssetsByEntity)
			r.Delete("/assets/{assetID}", a.DeleteAsset)

			// Notifications
			r.Get("/notifications", a.ListNotifications)
			r.Get("/notifications/unread-count", a.CountUnread)
			r.Post("/notifications/{notificationID}/read", a.MarkNotificationRead)
			r.Post("/notifications/read-all", a.MarkAllNotificationsRead)

			// Projects
			r.Get("/projects", a.ListProjects)
			r.Post("/projects", a.CreateProject)

			r.Route("/projects/{projectID}", func(r *router) {
				r.Use(a.requireProjectAccess)
				r.Get("/", a.GetProject)
				r.Put("/", a.UpdateProject)
				r.Delete("/", a.DeleteProject)

				// Project Members
				r.Get("/members", a.ListProjectMembers)
				r.Post("/members", a.AddProjectMember)
				r.Put("/members/{memberID}", a.UpdateProjectMember)
				r.Delete("/members/{memberID}", a.RemoveProjectMember)

				// States
				r.Get("/states", a.ListStates)
				r.Post("/states", a.CreateState)
				r.Put("/states/{stateID}", a.UpdateState)
				r.Delete("/states/{stateID}", a.DeleteState)

				// Labels
				r.Get("/labels", a.ListLabels)
				r.Post("/labels", a.CreateLabel)
				r.Put("/labels/{labelID}", a.UpdateLabel)
				r.Delete("/labels/{labelID}", a.DeleteLabel)

				// Estimates
				r.Get("/estimates", a.GetEstimateSystem)
				r.Put("/estimates", a.CreateOrUpdateEstimateSystem)

				// Issues
				r.Get("/issues", a.ListIssues)
				r.Post("/issues", a.CreateIssue)

				// Cycles
				r.Get("/cycles", a.ListCycles)
				r.Post("/cycles", a.CreateCycle)

				// Modules
				r.Get("/modules", a.ListModules)
				r.Post("/modules", a.CreateModule)

				r.Route("/modules/{moduleID}", func(r *router) {
					r.Get("/", a.GetModule)
					r.Put("/", a.UpdateModule)
					r.Delete("/", a.DeleteModule)

					r.Post("/issues", a.AddIssueToModule)
					r.Delete("/issues/{issueID}", a.RemoveIssueFromModule)
				})

				// Pages
				r.Get("/pages", a.ListPages)
				r.Post("/pages", a.CreatePage)

				r.Route("/pages/{pageID}", func(r *router) {
					r.Get("/", a.GetPage)
					r.Put("/", a.UpdatePage)
					r.Delete("/", a.DeletePage)
				})

				r.Route("/cycles/{cycleID}", func(r *router) {
					r.Get("/", a.GetCycle)
					r.Put("/", a.UpdateCycle)
					r.Delete("/", a.DeleteCycle)

					r.Post("/issues", a.AddIssueToCycle)
					r.Delete("/issues/{issueID}", a.RemoveIssueFromCycle)
				})

				r.Route("/issues/{issueID}", func(r *router) {
					r.Get("/", a.GetIssue)
					r.Put("/", a.UpdateIssue)
					r.Delete("/", a.DeleteIssue)

					r.Get("/comments", a.ListIssueComments)
					r.Post("/comments", a.CreateIssueComment)
					r.Put("/comments/{commentID}", a.UpdateIssueComment)
					r.Delete("/comments/{commentID}", a.DeleteIssueComment)

					r.Get("/activities", a.ListIssueActivities)

					// Work Logs
					r.Get("/work-logs", a.ListWorkLogs)
					r.Post("/work-logs", a.CreateWorkLog)
					r.Put("/work-logs/{workLogID}", a.UpdateWorkLog)
					r.Delete("/work-logs/{workLogID}", a.DeleteWorkLog)

					// Relations
					r.Get("/relations", a.ListIssueRelations)
					r.Post("/relations", a.AddIssueRelation)
					r.Delete("/relations/{relationID}", a.RemoveIssueRelation)

					// Links
					r.Get("/links", a.ListIssueLinks)
					r.Post("/links", a.CreateIssueLink)
					r.Put("/links/{linkID}", a.UpdateIssueLink)
					r.Delete("/links/{linkID}", a.DeleteIssueLink)
				})

				// Reports
				r.Get("/reports/issues", a.GetIssueAnalyticsReport)
				r.Get("/reports/work-logs", a.GetWorkLogReport)
				r.Get("/reports/cycles/{cycleID}", a.GetCycleReport)
				r.Get("/reports/velocity", a.GetVelocityReport)
				r.Get("/reports/health", a.GetProjectHealthReport)

				// Chat
				r.Get("/chat", a.ListChatHistory)
				r.Post("/chat", a.SendChatMessage)
				r.Delete("/chat", a.ClearChatHistory)
			})
		})

		// Invite accept (workspace-agnostic)
		r.Post("/invites/{token}/accept", a.AcceptWorkspaceInvite)
	})

	a.handler = r
	return a
}

// WithStorage sets the storage backend for file uploads.
func WithStorage(s *storage.Storage) func(*API) {
	return func(a *API) {
		a.storage = s
	}
}

// WithMailer sets the email mailer for sending notifications.
func WithMailer(m *email.Mailer) func(*API) {
	return func(a *API) {
		a.mailer = m
	}
}

func (a *API) Handler() http.Handler {
	return a.handler
}

func (a *API) healthCheck(w http.ResponseWriter, r *http.Request) error {
	return sendJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"version": "0.1.0",
	})
}
