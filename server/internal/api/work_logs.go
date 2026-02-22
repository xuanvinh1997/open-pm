package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

type CreateWorkLogRequest struct {
	DurationMinutes int    `json:"duration_minutes" validate:"required,min=1"`
	Description     string `json:"description" validate:"max=1000"`
	StartDate       string `json:"start_date" validate:"required"`
	EndDate         string `json:"end_date" validate:"required"`
}

type UpdateWorkLogRequest struct {
	DurationMinutes *int    `json:"duration_minutes" validate:"omitempty,min=1"`
	Description     *string `json:"description" validate:"omitempty,max=1000"`
	StartDate       *string `json:"start_date"`
	EndDate         *string `json:"end_date"`
}

// ListWorkLogs handles GET .../issues/{issueID}/work-logs
func (a *API) ListWorkLogs(w http.ResponseWriter, r *http.Request) error {
	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}

	logs, err := a.queries.ListWorkLogsByIssue(r.Context(), issueID)
	if err != nil {
		return internalServerError("failed to list work logs")
	}

	totalMinutes, err := a.queries.SumWorkLogMinutesByIssue(r.Context(), issueID)
	if err != nil {
		log.Warn().Err(err).Str("issue_id", issueID.String()).Msg("failed to sum work log minutes")
	}

	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"results":       logs,
		"total_minutes": totalMinutes,
	})
}

// CreateWorkLog handles POST .../issues/{issueID}/work-logs
func (a *API) CreateWorkLog(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)
	projectID := getProjectID(ctx)
	workspaceID := getWorkspaceID(ctx)

	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}

	var req CreateWorkLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return badRequestError("invalid start_date format, expected YYYY-MM-DD")
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return badRequestError("invalid end_date format, expected YYYY-MM-DD")
	}
	if endDate.Before(startDate) {
		return badRequestError("end_date must not be before start_date")
	}

	workLog, err := a.queries.CreateWorkLog(ctx, issueID, projectID, workspaceID, req.DurationMinutes, req.Description, startDate, endDate, userID)
	if err != nil {
		return internalServerError("failed to create work log")
	}

	// Log activity
	hours := req.DurationMinutes / 60
	mins := req.DurationMinutes % 60
	var comment string
	if hours > 0 && mins > 0 {
		comment = fmt.Sprintf("logged %dh %dm of work", hours, mins)
	} else if hours > 0 {
		comment = fmt.Sprintf("logged %dh of work", hours)
	} else {
		comment = fmt.Sprintf("logged %dm of work", mins)
	}

	if _, err := a.queries.CreateIssueActivity(ctx, CreateActivityParams{
		IssueID:     &issueID,
		ProjectID:   projectID,
		WorkspaceID: workspaceID,
		Verb:        "created",
		Field:       strPtr("work_log"),
		NewValue:    strPtr(fmt.Sprintf("%d", req.DurationMinutes)),
		Comment:     comment,
		ActorID:     &userID,
	}); err != nil {
		log.Warn().Err(err).Str("issue_id", issueID.String()).Msg("failed to log work log creation activity")
	}

	return sendJSON(w, http.StatusCreated, workLog)
}

// UpdateWorkLog handles PUT .../issues/{issueID}/work-logs/{workLogID}
func (a *API) UpdateWorkLog(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)

	workLogID, err := uuid.FromString(chi.URLParam(r, "workLogID"))
	if err != nil {
		return badRequestError("invalid work log ID")
	}

	// Check ownership: only the logger or admin can update
	existing, err := a.queries.GetWorkLogByID(ctx, workLogID)
	if err != nil {
		return notFoundError("work log not found")
	}
	if getProjectRole(ctx) < RoleAdmin && existing.LoggedBy != userID {
		return forbiddenError("only the author or an admin can edit this work log")
	}

	var req UpdateWorkLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	var startDate, endDate *time.Time
	if req.StartDate != nil {
		t, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return badRequestError("invalid start_date format, expected YYYY-MM-DD")
		}
		startDate = &t
	}
	if req.EndDate != nil {
		t, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return badRequestError("invalid end_date format, expected YYYY-MM-DD")
		}
		endDate = &t
	}
	if startDate != nil && endDate != nil && endDate.Before(*startDate) {
		return badRequestError("end_date must not be before start_date")
	}

	workLog, err := a.queries.UpdateWorkLog(ctx, workLogID, req.DurationMinutes, req.Description, startDate, endDate)
	if err != nil {
		return internalServerError("failed to update work log")
	}

	return sendJSON(w, http.StatusOK, workLog)
}

// DeleteWorkLog handles DELETE .../issues/{issueID}/work-logs/{workLogID}
func (a *API) DeleteWorkLog(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)

	workLogID, err := uuid.FromString(chi.URLParam(r, "workLogID"))
	if err != nil {
		return badRequestError("invalid work log ID")
	}

	// Check ownership: only the logger or admin can delete
	existing, err := a.queries.GetWorkLogByID(ctx, workLogID)
	if err != nil {
		return notFoundError("work log not found")
	}
	if getProjectRole(ctx) < RoleAdmin && existing.LoggedBy != userID {
		return forbiddenError("only the author or an admin can delete this work log")
	}

	if err := a.queries.DeleteWorkLog(ctx, workLogID); err != nil {
		return internalServerError("failed to delete work log")
	}

	return sendEmpty(w, http.StatusNoContent)
}

func strPtr(s string) *string {
	return &s
}
