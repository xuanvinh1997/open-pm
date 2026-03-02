package repository

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/open-pm/open-pm/server/internal/api"
)

// --- Reports: Issue Analytics ---

func (r *Repository) CountIssuesByState(ctx context.Context, projectID uuid.UUID) ([]*api.IssuesByStateReport, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT s.id, s.name, s."group", s.color, COUNT(i.id)::int
		 FROM states s
		 LEFT JOIN issues i ON i.state_id = s.id AND i.deleted_at IS NULL
		 WHERE s.project_id = $1 AND s.deleted_at IS NULL
		 GROUP BY s.id, s.name, s."group", s.color
		 ORDER BY s.sequence ASC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*api.IssuesByStateReport
	for rows.Next() {
		var r api.IssuesByStateReport
		if err := rows.Scan(&r.StateID, &r.StateName, &r.StateGroup, &r.Color, &r.Count); err != nil {
			return nil, err
		}
		results = append(results, &r)
	}
	if results == nil {
		results = []*api.IssuesByStateReport{}
	}
	return results, nil
}

func (r *Repository) CountIssuesByPriority(ctx context.Context, projectID uuid.UUID) ([]*api.IssuesByPriorityReport, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT priority, COUNT(*)::int
		 FROM issues
		 WHERE project_id = $1 AND deleted_at IS NULL
		 GROUP BY priority
		 ORDER BY CASE priority
		   WHEN 'urgent' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3
		   WHEN 'low' THEN 4 WHEN 'none' THEN 5 ELSE 6 END`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*api.IssuesByPriorityReport
	for rows.Next() {
		var r api.IssuesByPriorityReport
		if err := rows.Scan(&r.Priority, &r.Count); err != nil {
			return nil, err
		}
		results = append(results, &r)
	}
	if results == nil {
		results = []*api.IssuesByPriorityReport{}
	}
	return results, nil
}

func (r *Repository) CountIssuesByType(ctx context.Context, projectID uuid.UUID) ([]*api.IssuesByTypeReport, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT issue_type, COUNT(*)::int
		 FROM issues
		 WHERE project_id = $1 AND deleted_at IS NULL
		 GROUP BY issue_type
		 ORDER BY issue_type`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*api.IssuesByTypeReport
	for rows.Next() {
		var r api.IssuesByTypeReport
		if err := rows.Scan(&r.IssueType, &r.Count); err != nil {
			return nil, err
		}
		results = append(results, &r)
	}
	if results == nil {
		results = []*api.IssuesByTypeReport{}
	}
	return results, nil
}

func (r *Repository) CountIssuesByAssignee(ctx context.Context, projectID uuid.UUID) ([]*api.IssuesByAssigneeReport, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT u.id, u.display_name, u.avatar_url, COUNT(i.id)::int
		 FROM issue_assignees ia
		 JOIN users u ON ia.assignee_id = u.id
		 JOIN issues i ON ia.issue_id = i.id
		 WHERE i.project_id = $1 AND i.deleted_at IS NULL
		 GROUP BY u.id, u.display_name, u.avatar_url
		 ORDER BY COUNT(i.id) DESC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*api.IssuesByAssigneeReport
	for rows.Next() {
		var r api.IssuesByAssigneeReport
		if err := rows.Scan(&r.UserID, &r.DisplayName, &r.AvatarURL, &r.Count); err != nil {
			return nil, err
		}
		results = append(results, &r)
	}
	if results == nil {
		results = []*api.IssuesByAssigneeReport{}
	}
	return results, nil
}

// --- Reports: Work Logs ---

func (r *Repository) SumWorkLogsByMember(ctx context.Context, projectID uuid.UUID, startDate, endDate time.Time) ([]*api.WorkLogSummaryByMember, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT u.id, u.display_name, u.avatar_url, COALESCE(SUM(w.duration_minutes), 0)::int, COUNT(w.id)::int
		 FROM work_logs w
		 JOIN users u ON w.logged_by = u.id
		 WHERE w.project_id = $1 AND w.deleted_at IS NULL
		   AND w.start_date >= $2 AND w.end_date <= $3
		 GROUP BY u.id, u.display_name, u.avatar_url
		 ORDER BY SUM(w.duration_minutes) DESC`, projectID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*api.WorkLogSummaryByMember
	for rows.Next() {
		var r api.WorkLogSummaryByMember
		if err := rows.Scan(&r.UserID, &r.DisplayName, &r.AvatarURL, &r.TotalMinutes, &r.LogCount); err != nil {
			return nil, err
		}
		results = append(results, &r)
	}
	if results == nil {
		results = []*api.WorkLogSummaryByMember{}
	}
	return results, nil
}

func (r *Repository) SumWorkLogsByIssue(ctx context.Context, projectID uuid.UUID, startDate, endDate time.Time) ([]*api.WorkLogSummaryByIssue, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT i.id, i.name, i.sequence_id, COALESCE(SUM(w.duration_minutes), 0)::int
		 FROM work_logs w
		 JOIN issues i ON w.issue_id = i.id
		 WHERE w.project_id = $1 AND w.deleted_at IS NULL
		   AND w.start_date >= $2 AND w.end_date <= $3
		 GROUP BY i.id, i.name, i.sequence_id
		 ORDER BY SUM(w.duration_minutes) DESC`, projectID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*api.WorkLogSummaryByIssue
	for rows.Next() {
		var r api.WorkLogSummaryByIssue
		if err := rows.Scan(&r.IssueID, &r.IssueName, &r.SequenceID, &r.TotalMinutes); err != nil {
			return nil, err
		}
		results = append(results, &r)
	}
	if results == nil {
		results = []*api.WorkLogSummaryByIssue{}
	}
	return results, nil
}

func (r *Repository) SumWorkLogsByDate(ctx context.Context, projectID uuid.UUID, startDate, endDate time.Time) ([]*api.WorkLogSummaryByDate, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT w.start_date::text, COALESCE(SUM(w.duration_minutes), 0)::int
		 FROM work_logs w
		 WHERE w.project_id = $1 AND w.deleted_at IS NULL
		   AND w.start_date >= $2 AND w.end_date <= $3
		 GROUP BY w.start_date
		 ORDER BY w.start_date ASC`, projectID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*api.WorkLogSummaryByDate
	for rows.Next() {
		var r api.WorkLogSummaryByDate
		if err := rows.Scan(&r.Date, &r.TotalMinutes); err != nil {
			return nil, err
		}
		results = append(results, &r)
	}
	if results == nil {
		results = []*api.WorkLogSummaryByDate{}
	}
	return results, nil
}

// --- Reports: Project Health ---

const issueColumns = `id, project_id, workspace_id, parent_id, state_id, name, description_html, description_json, description_stripped,
		        priority, issue_type, start_date, target_date, sequence_id, sort_order, estimate_point, completed_at, archived_at, is_draft, created_by, updated_by, created_at, updated_at,
		        reporter_id, resolution_id, resolved_at`

func scanIssueRow(rows interface{ Scan(dest ...interface{}) error }) (*api.Issue, error) {
	var i api.Issue
	err := rows.Scan(&i.ID, &i.ProjectID, &i.WorkspaceID, &i.ParentID, &i.StateID, &i.Name,
		&i.DescriptionHTML, &i.DescriptionJSON, &i.DescriptionStripped,
		&i.Priority, &i.IssueType, &i.StartDate, &i.TargetDate, &i.SequenceID, &i.SortOrder,
		&i.EstimatePoint, &i.CompletedAt, &i.ArchivedAt, &i.IsDraft, &i.CreatedBy, &i.UpdatedBy, &i.CreatedAt, &i.UpdatedAt,
		&i.ReporterID, &i.ResolutionID, &i.ResolvedAt)
	return &i, err
}

func (r *Repository) ListOverdueIssues(ctx context.Context, projectID uuid.UUID) ([]*api.Issue, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+issueColumns+`
		 FROM issues
		 WHERE project_id = $1 AND deleted_at IS NULL
		   AND target_date < CURRENT_DATE AND completed_at IS NULL
		 ORDER BY target_date ASC
		 LIMIT 50`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []*api.Issue
	for rows.Next() {
		i, err := scanIssueRow(rows)
		if err != nil {
			return nil, err
		}
		issues = append(issues, i)
	}
	if issues == nil {
		issues = []*api.Issue{}
	}
	return issues, nil
}

func (r *Repository) ListBlockedIssueIDs(ctx context.Context, projectID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT DISTINCT ir.issue_id
		 FROM issue_relations ir
		 JOIN issues i ON ir.issue_id = i.id
		 WHERE i.project_id = $1 AND i.deleted_at IS NULL AND i.completed_at IS NULL
		   AND ir.relation_type = 'blocked_by'`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if ids == nil {
		ids = []uuid.UUID{}
	}
	return ids, nil
}

func (r *Repository) ListUnestimatedIssues(ctx context.Context, projectID uuid.UUID) ([]*api.Issue, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+issueColumns+`
		 FROM issues
		 WHERE project_id = $1 AND deleted_at IS NULL
		   AND estimate_point IS NULL AND completed_at IS NULL
		 ORDER BY created_at DESC
		 LIMIT 50`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []*api.Issue
	for rows.Next() {
		i, err := scanIssueRow(rows)
		if err != nil {
			return nil, err
		}
		issues = append(issues, i)
	}
	if issues == nil {
		issues = []*api.Issue{}
	}
	return issues, nil
}

// --- Reports: Sprints ---

func (r *Repository) ListSprintIssueIDs(ctx context.Context, sprintID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT ci.issue_id
		 FROM sprint_issues ci
		 JOIN issues i ON ci.issue_id = i.id
		 WHERE ci.sprint_id = $1 AND i.deleted_at IS NULL`, sprintID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if ids == nil {
		ids = []uuid.UUID{}
	}
	return ids, nil
}

func (r *Repository) CountSprintIssuesByState(ctx context.Context, sprintID uuid.UUID) ([]*api.IssuesByStateReport, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT s.id, s.name, s."group", s.color, COUNT(i.id)::int
		 FROM sprint_issues ci
		 JOIN issues i ON ci.issue_id = i.id
		 JOIN states s ON i.state_id = s.id
		 WHERE ci.sprint_id = $1 AND i.deleted_at IS NULL
		 GROUP BY s.id, s.name, s."group", s.color
		 ORDER BY s.sequence ASC`, sprintID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*api.IssuesByStateReport
	for rows.Next() {
		var r api.IssuesByStateReport
		if err := rows.Scan(&r.StateID, &r.StateName, &r.StateGroup, &r.Color, &r.Count); err != nil {
			return nil, err
		}
		results = append(results, &r)
	}
	if results == nil {
		results = []*api.IssuesByStateReport{}
	}
	return results, nil
}

func (r *Repository) CountCompletedSprintIssues(ctx context.Context, sprintID uuid.UUID) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*)
		 FROM sprint_issues ci
		 JOIN issues i ON ci.issue_id = i.id
		 WHERE ci.sprint_id = $1 AND i.deleted_at IS NULL AND i.completed_at IS NOT NULL`, sprintID,
	).Scan(&count)
	return count, err
}

func (r *Repository) CountTotalSprintIssues(ctx context.Context, sprintID uuid.UUID) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*)
		 FROM sprint_issues ci
		 JOIN issues i ON ci.issue_id = i.id
		 WHERE ci.sprint_id = $1 AND i.deleted_at IS NULL`, sprintID,
	).Scan(&count)
	return count, err
}

func (r *Repository) SumCompletedSprintPoints(ctx context.Context, sprintID uuid.UUID) (int64, error) {
	var sum int64
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(i.estimate_point), 0)::bigint
		 FROM sprint_issues ci
		 JOIN issues i ON ci.issue_id = i.id
		 WHERE ci.sprint_id = $1 AND i.deleted_at IS NULL AND i.completed_at IS NOT NULL AND i.estimate_point IS NOT NULL`, sprintID,
	).Scan(&sum)
	return sum, err
}

func (r *Repository) SumTotalSprintPoints(ctx context.Context, sprintID uuid.UUID) (int64, error) {
	var sum int64
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(i.estimate_point), 0)::bigint
		 FROM sprint_issues ci
		 JOIN issues i ON ci.issue_id = i.id
		 WHERE ci.sprint_id = $1 AND i.deleted_at IS NULL AND i.estimate_point IS NOT NULL`, sprintID,
	).Scan(&sum)
	return sum, err
}
