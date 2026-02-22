package repository

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/open-pm/open-pm/server/internal/api"
)

func (r *Repository) CreateWorkLog(ctx context.Context, issueID, projectID, workspaceID uuid.UUID, durationMinutes int, description string, startDate, endDate time.Time, loggedBy uuid.UUID) (*api.WorkLog, error) {
	var w api.WorkLog
	err := r.pool.QueryRow(ctx,
		`INSERT INTO work_logs (issue_id, project_id, workspace_id, duration_minutes, description, start_date, end_date, logged_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, issue_id, project_id, workspace_id, duration_minutes, description, start_date, end_date, logged_by, created_at, updated_at`,
		issueID, projectID, workspaceID, durationMinutes, description, startDate, endDate, loggedBy,
	).Scan(&w.ID, &w.IssueID, &w.ProjectID, &w.WorkspaceID, &w.DurationMinutes, &w.Description, &w.StartDate, &w.EndDate, &w.LoggedBy, &w.CreatedAt, &w.UpdatedAt)
	return &w, err
}

func (r *Repository) GetWorkLogByID(ctx context.Context, id uuid.UUID) (*api.WorkLog, error) {
	var w api.WorkLog
	err := r.pool.QueryRow(ctx,
		`SELECT id, issue_id, project_id, workspace_id, duration_minutes, description, start_date, end_date, logged_by, created_at, updated_at
		 FROM work_logs WHERE id = $1 AND deleted_at IS NULL`, id,
	).Scan(&w.ID, &w.IssueID, &w.ProjectID, &w.WorkspaceID, &w.DurationMinutes, &w.Description, &w.StartDate, &w.EndDate, &w.LoggedBy, &w.CreatedAt, &w.UpdatedAt)
	return &w, err
}

func (r *Repository) ListWorkLogsByIssue(ctx context.Context, issueID uuid.UUID) ([]*api.WorkLogWithUser, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT w.id, w.issue_id, w.project_id, w.workspace_id, w.duration_minutes, w.description, w.start_date, w.end_date, w.logged_by, w.created_at, w.updated_at,
		        u.first_name, u.last_name, u.display_name, u.avatar_url
		 FROM work_logs w
		 JOIN users u ON u.id = w.logged_by
		 WHERE w.issue_id = $1 AND w.deleted_at IS NULL
		 ORDER BY w.start_date DESC, w.created_at DESC`, issueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*api.WorkLogWithUser
	for rows.Next() {
		var w api.WorkLogWithUser
		if err := rows.Scan(
			&w.ID, &w.IssueID, &w.ProjectID, &w.WorkspaceID, &w.DurationMinutes, &w.Description, &w.StartDate, &w.EndDate, &w.LoggedBy, &w.CreatedAt, &w.UpdatedAt,
			&w.FirstName, &w.LastName, &w.DisplayName, &w.AvatarURL,
		); err != nil {
			return nil, err
		}
		logs = append(logs, &w)
	}
	if logs == nil {
		logs = []*api.WorkLogWithUser{}
	}
	return logs, nil
}

func (r *Repository) UpdateWorkLog(ctx context.Context, id uuid.UUID, durationMinutes *int, description *string, startDate, endDate *time.Time) (*api.WorkLog, error) {
	var w api.WorkLog
	err := r.pool.QueryRow(ctx,
		`UPDATE work_logs SET
			duration_minutes = COALESCE($2, duration_minutes),
			description = COALESCE($3, description),
			start_date = COALESCE($4, start_date),
			end_date = COALESCE($5, end_date)
		 WHERE id = $1 AND deleted_at IS NULL
		 RETURNING id, issue_id, project_id, workspace_id, duration_minutes, description, start_date, end_date, logged_by, created_at, updated_at`,
		id, durationMinutes, description, startDate, endDate,
	).Scan(&w.ID, &w.IssueID, &w.ProjectID, &w.WorkspaceID, &w.DurationMinutes, &w.Description, &w.StartDate, &w.EndDate, &w.LoggedBy, &w.CreatedAt, &w.UpdatedAt)
	return &w, err
}

func (r *Repository) DeleteWorkLog(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE work_logs SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *Repository) SumWorkLogMinutesByIssue(ctx context.Context, issueID uuid.UUID) (int64, error) {
	var total int64
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(duration_minutes), 0) FROM work_logs WHERE issue_id = $1 AND deleted_at IS NULL`, issueID,
	).Scan(&total)
	return total, err
}
