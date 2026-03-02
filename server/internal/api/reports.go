package api

import (
	"math"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

// GetIssueAnalyticsReport handles GET .../projects/{projectID}/reports/issues
func (a *API) GetIssueAnalyticsReport(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	projectID := getProjectID(ctx)

	byState, err := a.queries.CountIssuesByState(ctx, projectID)
	if err != nil {
		log.Error().Err(err).Msg("failed to count issues by state")
		return internalServerError("failed to load issue analytics")
	}

	byPriority, err := a.queries.CountIssuesByPriority(ctx, projectID)
	if err != nil {
		log.Error().Err(err).Msg("failed to count issues by priority")
		return internalServerError("failed to load issue analytics")
	}

	byType, err := a.queries.CountIssuesByType(ctx, projectID)
	if err != nil {
		log.Error().Err(err).Msg("failed to count issues by type")
		return internalServerError("failed to load issue analytics")
	}

	byAssignee, err := a.queries.CountIssuesByAssignee(ctx, projectID)
	if err != nil {
		log.Error().Err(err).Msg("failed to count issues by assignee")
		return internalServerError("failed to load issue analytics")
	}

	// Compute totals from state data
	var total, completed, overdue int
	for _, s := range byState {
		total += s.Count
		if s.StateGroup == "completed" {
			completed += s.Count
		}
	}
	// Count overdue separately
	overdueIssues, err := a.queries.ListOverdueIssues(ctx, projectID)
	if err == nil {
		overdue = len(overdueIssues)
	}

	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"by_state":    byState,
		"by_priority": byPriority,
		"by_type":     byType,
		"by_assignee": byAssignee,
		"total":       total,
		"completed":   completed,
		"overdue":     overdue,
	})
}

// GetWorkLogReport handles GET .../projects/{projectID}/reports/work-logs
func (a *API) GetWorkLogReport(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	projectID := getProjectID(ctx)

	// Parse date range from query params, default to last 30 days
	now := time.Now()
	startDate := now.AddDate(0, 0, -30)
	endDate := now

	if sd := r.URL.Query().Get("start_date"); sd != "" {
		if t, err := time.Parse("2006-01-02", sd); err == nil {
			startDate = t
		}
	}
	if ed := r.URL.Query().Get("end_date"); ed != "" {
		if t, err := time.Parse("2006-01-02", ed); err == nil {
			endDate = t
		}
	}

	byMember, err := a.queries.SumWorkLogsByMember(ctx, projectID, startDate, endDate)
	if err != nil {
		log.Error().Err(err).Msg("failed to sum work logs by member")
		return internalServerError("failed to load work log report")
	}

	byIssue, err := a.queries.SumWorkLogsByIssue(ctx, projectID, startDate, endDate)
	if err != nil {
		log.Error().Err(err).Msg("failed to sum work logs by issue")
		return internalServerError("failed to load work log report")
	}

	byDate, err := a.queries.SumWorkLogsByDate(ctx, projectID, startDate, endDate)
	if err != nil {
		log.Error().Err(err).Msg("failed to sum work logs by date")
		return internalServerError("failed to load work log report")
	}

	var totalMinutes int
	for _, m := range byMember {
		totalMinutes += m.TotalMinutes
	}

	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"by_member":     byMember,
		"by_issue":      byIssue,
		"by_date":       byDate,
		"total_minutes": totalMinutes,
		"start_date":    startDate.Format("2006-01-02"),
		"end_date":      endDate.Format("2006-01-02"),
	})
}

// GetSprintReport handles GET .../projects/{projectID}/reports/cycles/{sprintID}
func (a *API) GetSprintReport(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	sprintID, err := uuid.FromString(chi.URLParam(r, "sprintID"))
	if err != nil {
		return badRequestError("invalid cycle ID")
	}

	cycle, err := a.queries.GetSprintByID(ctx, sprintID)
	if err != nil {
		return notFoundError("cycle not found")
	}

	totalCount, err := a.queries.CountTotalSprintIssues(ctx, sprintID)
	if err != nil {
		return internalServerError("failed to count cycle issues")
	}

	completedCount, err := a.queries.CountCompletedSprintIssues(ctx, sprintID)
	if err != nil {
		return internalServerError("failed to count completed cycle issues")
	}

	byState, err := a.queries.CountSprintIssuesByState(ctx, sprintID)
	if err != nil {
		return internalServerError("failed to count cycle issues by state")
	}

	totalPoints, _ := a.queries.SumTotalSprintPoints(ctx, sprintID)
	completedPoints, _ := a.queries.SumCompletedSprintPoints(ctx, sprintID)

	// Build burndown data
	var burndown []BurndownPoint
	if cycle.StartDate != nil && cycle.EndDate != nil {
		burndown = buildBurndown(*cycle.StartDate, *cycle.EndDate, int(totalCount), int(completedCount), int(totalPoints), int(completedPoints))
	}

	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"sprint":            cycle,
		"total_issues":      totalCount,
		"completed_issues":  completedCount,
		"total_points":      totalPoints,
		"completed_points":  completedPoints,
		"burndown_points":   burndown,
		"issues_by_state":   byState,
	})
}

// GetVelocityReport handles GET .../projects/{projectID}/reports/velocity
func (a *API) GetVelocityReport(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	projectID := getProjectID(ctx)

	cycles, err := a.queries.ListSprintsByProject(ctx, projectID)
	if err != nil {
		return internalServerError("failed to list cycles")
	}

	var points []VelocityPoint
	for _, c := range cycles {
		total, err := a.queries.CountTotalSprintIssues(ctx, c.ID)
		if err != nil {
			continue
		}
		completed, err := a.queries.CountCompletedSprintIssues(ctx, c.ID)
		if err != nil {
			continue
		}
		totalPoints, _ := a.queries.SumTotalSprintPoints(ctx, c.ID)
		completedPoints, _ := a.queries.SumCompletedSprintPoints(ctx, c.ID)
		points = append(points, VelocityPoint{
			SprintID:        c.ID,
			SprintName:      c.Name,
			CompletedCount:  int(completed),
			TotalCount:      int(total),
			CompletedPoints: int(completedPoints),
			TotalPoints:     int(totalPoints),
		})
	}
	if points == nil {
		points = []VelocityPoint{}
	}

	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"results": points,
	})
}

// GetProjectHealthReport handles GET .../projects/{projectID}/reports/health
func (a *API) GetProjectHealthReport(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	projectID := getProjectID(ctx)

	overdueIssues, err := a.queries.ListOverdueIssues(ctx, projectID)
	if err != nil {
		return internalServerError("failed to list overdue issues")
	}

	blockedIDs, err := a.queries.ListBlockedIssueIDs(ctx, projectID)
	if err != nil {
		return internalServerError("failed to list blocked issues")
	}

	// Fetch full blocked issue details
	var blockedIssues []*Issue
	for _, id := range blockedIDs {
		issue, err := a.queries.GetIssueByID(ctx, id)
		if err == nil {
			blockedIssues = append(blockedIssues, issue)
		}
	}
	if blockedIssues == nil {
		blockedIssues = []*Issue{}
	}

	unestimatedIssues, err := a.queries.ListUnestimatedIssues(ctx, projectID)
	if err != nil {
		return internalServerError("failed to list unestimated issues")
	}

	return sendJSON(w, http.StatusOK, ProjectHealthReport{
		OverdueIssues:     overdueIssues,
		BlockedIssues:     blockedIssues,
		UnestimatedIssues: unestimatedIssues,
		OverdueCount:      len(overdueIssues),
		BlockedCount:      len(blockedIssues),
		UnestimatedCount:  len(unestimatedIssues),
	})
}

// buildBurndown creates an ideal burndown chart from start to end date.
// It supports both issue-count-based and story-point-based burndown curves.
func buildBurndown(startDate, endDate time.Time, totalIssues, completedIssues, totalPoints, completedPoints int) []BurndownPoint {
	start := startDate.Truncate(24 * time.Hour)
	end := endDate.Truncate(24 * time.Hour)
	now := time.Now().Truncate(24 * time.Hour)

	if end.Before(start) {
		return []BurndownPoint{}
	}

	totalDays := int(end.Sub(start).Hours()/24) + 1
	if totalDays < 1 {
		totalDays = 1
	}

	remainingCount := totalIssues
	remainingPoints := totalPoints
	var points []BurndownPoint

	for d := start; !d.After(end) && !d.After(now); d = d.AddDate(0, 0, 1) {
		dayIndex := int(d.Sub(start).Hours() / 24)

		// Ideal remaining (count-based)
		idealRemaining := float64(totalIssues) * (1.0 - float64(dayIndex)/float64(totalDays-1))
		idealRemaining = math.Max(0, idealRemaining)

		// Ideal remaining (point-based)
		idealRemainingPoints := float64(totalPoints) * (1.0 - float64(dayIndex)/float64(totalDays-1))
		idealRemainingPoints = math.Max(0, idealRemainingPoints)

		// Simple linear interpolation for completed counts
		if totalDays > 1 {
			progress := float64(dayIndex) / float64(totalDays-1)
			if completedIssues > 0 {
				done := int(math.Round(float64(completedIssues) * progress))
				remainingCount = totalIssues - done
			}
			if completedPoints > 0 {
				donePoints := int(math.Round(float64(completedPoints) * progress))
				remainingPoints = totalPoints - donePoints
			}
		}

		points = append(points, BurndownPoint{
			Date:                 d.Format("2006-01-02"),
			RemainingCount:       remainingCount,
			CompletedCount:       totalIssues - remainingCount,
			IdealRemaining:       math.Round(idealRemaining*100) / 100,
			RemainingPoints:      remainingPoints,
			CompletedPoints:      totalPoints - remainingPoints,
			IdealRemainingPoints: math.Round(idealRemainingPoints*100) / 100,
		})
	}

	if points == nil {
		points = []BurndownPoint{}
	}
	return points
}
