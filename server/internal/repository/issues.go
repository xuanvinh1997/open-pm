package repository

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/open-pm/open-pm/server/internal/api"
)

// --- States ---

func (r *Repository) CreateState(ctx context.Context, projectID, workspaceID uuid.UUID, name, description, color, group string, sequence float64, isDefault bool) (*api.State, error) {
	var s api.State
	err := r.pool.QueryRow(ctx,
		`INSERT INTO states (project_id, workspace_id, name, description, color, "group", sequence, is_default)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, project_id, workspace_id, name, description, color, "group", sequence, is_default, created_at, updated_at`,
		projectID, workspaceID, name, description, color, group, sequence, isDefault,
	).Scan(&s.ID, &s.ProjectID, &s.WorkspaceID, &s.Name, &s.Description, &s.Color, &s.Group, &s.Sequence, &s.IsDefault, &s.CreatedAt, &s.UpdatedAt)
	return &s, err
}

func (r *Repository) GetStateByID(ctx context.Context, id uuid.UUID) (*api.State, error) {
	var s api.State
	err := r.pool.QueryRow(ctx,
		`SELECT id, project_id, workspace_id, name, description, color, "group", sequence, is_default, created_at, updated_at
		 FROM states WHERE id = $1 AND deleted_at IS NULL`, id,
	).Scan(&s.ID, &s.ProjectID, &s.WorkspaceID, &s.Name, &s.Description, &s.Color, &s.Group, &s.Sequence, &s.IsDefault, &s.CreatedAt, &s.UpdatedAt)
	return &s, err
}

func (r *Repository) ListStatesByProject(ctx context.Context, projectID uuid.UUID) ([]*api.State, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, project_id, workspace_id, name, description, color, "group", sequence, is_default, created_at, updated_at
		 FROM states WHERE project_id = $1 AND deleted_at IS NULL ORDER BY sequence ASC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var states []*api.State
	for rows.Next() {
		var s api.State
		if err := rows.Scan(&s.ID, &s.ProjectID, &s.WorkspaceID, &s.Name, &s.Description, &s.Color, &s.Group, &s.Sequence, &s.IsDefault, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		states = append(states, &s)
	}
	if states == nil {
		states = []*api.State{}
	}
	return states, nil
}

func (r *Repository) UpdateState(ctx context.Context, id uuid.UUID, name, description, color *string, group *string, sequence *float64, isDefault *bool) (*api.State, error) {
	var s api.State
	err := r.pool.QueryRow(ctx,
		`UPDATE states SET
			name = COALESCE($2, name), description = COALESCE($3, description),
			color = COALESCE($4, color), "group" = COALESCE($5, "group"),
			sequence = COALESCE($6, sequence), is_default = COALESCE($7, is_default)
		 WHERE id = $1 AND deleted_at IS NULL
		 RETURNING id, project_id, workspace_id, name, description, color, "group", sequence, is_default, created_at, updated_at`,
		id, name, description, color, group, sequence, isDefault,
	).Scan(&s.ID, &s.ProjectID, &s.WorkspaceID, &s.Name, &s.Description, &s.Color, &s.Group, &s.Sequence, &s.IsDefault, &s.CreatedAt, &s.UpdatedAt)
	return &s, err
}

func (r *Repository) DeleteState(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE states SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}

// --- Labels ---

func (r *Repository) CreateLabel(ctx context.Context, projectID *uuid.UUID, workspaceID uuid.UUID, parentID *uuid.UUID, name, description, color string, sortOrder float64) (*api.Label, error) {
	var l api.Label
	err := r.pool.QueryRow(ctx,
		`INSERT INTO labels (project_id, workspace_id, parent_id, name, description, color, sort_order)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, project_id, workspace_id, parent_id, name, description, color, sort_order, created_at, updated_at`,
		projectID, workspaceID, parentID, name, description, color, sortOrder,
	).Scan(&l.ID, &l.ProjectID, &l.WorkspaceID, &l.ParentID, &l.Name, &l.Description, &l.Color, &l.SortOrder, &l.CreatedAt, &l.UpdatedAt)
	return &l, err
}

func (r *Repository) GetLabelByID(ctx context.Context, id uuid.UUID) (*api.Label, error) {
	var l api.Label
	err := r.pool.QueryRow(ctx,
		`SELECT id, project_id, workspace_id, parent_id, name, description, color, sort_order, created_at, updated_at
		 FROM labels WHERE id = $1 AND deleted_at IS NULL`, id,
	).Scan(&l.ID, &l.ProjectID, &l.WorkspaceID, &l.ParentID, &l.Name, &l.Description, &l.Color, &l.SortOrder, &l.CreatedAt, &l.UpdatedAt)
	return &l, err
}

func (r *Repository) ListLabelsByProject(ctx context.Context, projectID uuid.UUID) ([]*api.Label, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, project_id, workspace_id, parent_id, name, description, color, sort_order, created_at, updated_at
		 FROM labels WHERE project_id = $1 AND deleted_at IS NULL ORDER BY sort_order ASC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labels []*api.Label
	for rows.Next() {
		var l api.Label
		if err := rows.Scan(&l.ID, &l.ProjectID, &l.WorkspaceID, &l.ParentID, &l.Name, &l.Description, &l.Color, &l.SortOrder, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		labels = append(labels, &l)
	}
	if labels == nil {
		labels = []*api.Label{}
	}
	return labels, nil
}

func (r *Repository) UpdateLabel(ctx context.Context, id uuid.UUID, name, description, color *string, parentID *uuid.UUID, sortOrder *float64) (*api.Label, error) {
	var l api.Label
	err := r.pool.QueryRow(ctx,
		`UPDATE labels SET
			name = COALESCE($2, name), description = COALESCE($3, description),
			color = COALESCE($4, color), parent_id = $5, sort_order = COALESCE($6, sort_order)
		 WHERE id = $1 AND deleted_at IS NULL
		 RETURNING id, project_id, workspace_id, parent_id, name, description, color, sort_order, created_at, updated_at`,
		id, name, description, color, parentID, sortOrder,
	).Scan(&l.ID, &l.ProjectID, &l.WorkspaceID, &l.ParentID, &l.Name, &l.Description, &l.Color, &l.SortOrder, &l.CreatedAt, &l.UpdatedAt)
	return &l, err
}

func (r *Repository) DeleteLabel(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE labels SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}

// --- Issues ---

func (r *Repository) CreateIssue(ctx context.Context, params api.CreateIssueParams) (*api.Issue, error) {
	var i api.Issue
	err := r.pool.QueryRow(ctx,
		`INSERT INTO issues (project_id, workspace_id, parent_id, state_id, name, description_html, description_json, description_stripped, priority, issue_type, start_date, target_date, is_draft, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		 RETURNING id, project_id, workspace_id, parent_id, state_id, name, description_html, description_json, description_stripped,
		           priority, issue_type, start_date, target_date, sequence_id, sort_order, completed_at, archived_at, is_draft, created_by, updated_by, created_at, updated_at`,
		params.ProjectID, params.WorkspaceID, params.ParentID, params.StateID,
		params.Name, params.DescriptionHTML, params.DescriptionJSON, params.DescriptionStripped,
		params.Priority, params.IssueType, params.StartDate, params.TargetDate, params.IsDraft, params.CreatedBy,
	).Scan(&i.ID, &i.ProjectID, &i.WorkspaceID, &i.ParentID, &i.StateID, &i.Name,
		&i.DescriptionHTML, &i.DescriptionJSON, &i.DescriptionStripped,
		&i.Priority, &i.IssueType, &i.StartDate, &i.TargetDate, &i.SequenceID, &i.SortOrder,
		&i.CompletedAt, &i.ArchivedAt, &i.IsDraft, &i.CreatedBy, &i.UpdatedBy, &i.CreatedAt, &i.UpdatedAt)
	return &i, err
}

func (r *Repository) GetIssueByID(ctx context.Context, id uuid.UUID) (*api.Issue, error) {
	var i api.Issue
	err := r.pool.QueryRow(ctx,
		`SELECT id, project_id, workspace_id, parent_id, state_id, name, description_html, description_json, description_stripped,
		        priority, issue_type, start_date, target_date, sequence_id, sort_order, completed_at, archived_at, is_draft, created_by, updated_by, created_at, updated_at
		 FROM issues WHERE id = $1 AND deleted_at IS NULL`, id,
	).Scan(&i.ID, &i.ProjectID, &i.WorkspaceID, &i.ParentID, &i.StateID, &i.Name,
		&i.DescriptionHTML, &i.DescriptionJSON, &i.DescriptionStripped,
		&i.Priority, &i.IssueType, &i.StartDate, &i.TargetDate, &i.SequenceID, &i.SortOrder,
		&i.CompletedAt, &i.ArchivedAt, &i.IsDraft, &i.CreatedBy, &i.UpdatedBy, &i.CreatedAt, &i.UpdatedAt)
	return &i, err
}

func (r *Repository) ListIssuesByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]*api.Issue, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, project_id, workspace_id, parent_id, state_id, name, description_html, description_json, description_stripped,
		        priority, issue_type, start_date, target_date, sequence_id, sort_order, completed_at, archived_at, is_draft, created_by, updated_by, created_at, updated_at
		 FROM issues WHERE project_id = $1 AND deleted_at IS NULL
		 ORDER BY sort_order ASC, created_at DESC
		 LIMIT $2 OFFSET $3`, projectID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []*api.Issue
	for rows.Next() {
		var i api.Issue
		if err := rows.Scan(&i.ID, &i.ProjectID, &i.WorkspaceID, &i.ParentID, &i.StateID, &i.Name,
			&i.DescriptionHTML, &i.DescriptionJSON, &i.DescriptionStripped,
			&i.Priority, &i.IssueType, &i.StartDate, &i.TargetDate, &i.SequenceID, &i.SortOrder,
			&i.CompletedAt, &i.ArchivedAt, &i.IsDraft, &i.CreatedBy, &i.UpdatedBy, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		issues = append(issues, &i)
	}
	if issues == nil {
		issues = []*api.Issue{}
	}
	return issues, nil
}

func (r *Repository) CountIssuesByProject(ctx context.Context, projectID uuid.UUID) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM issues WHERE project_id = $1 AND deleted_at IS NULL`, projectID).Scan(&count)
	return count, err
}

func (r *Repository) UpdateIssue(ctx context.Context, id uuid.UUID, params api.UpdateIssueParams) (*api.Issue, error) {
	var i api.Issue
	err := r.pool.QueryRow(ctx,
		`UPDATE issues SET
			state_id = COALESCE($2, state_id), name = COALESCE($3, name),
			description_html = COALESCE($4, description_html),
			description_json = COALESCE($5, description_json),
			description_stripped = COALESCE($6, description_stripped),
			priority = COALESCE($7, priority),
			issue_type = COALESCE($8, issue_type),
			start_date = $9, target_date = $10, parent_id = $11,
			sort_order = COALESCE($12, sort_order), is_draft = COALESCE($13, is_draft),
			completed_at = $14, archived_at = $15, updated_by = $16
		 WHERE id = $1 AND deleted_at IS NULL
		 RETURNING id, project_id, workspace_id, parent_id, state_id, name, description_html, description_json, description_stripped,
		           priority, issue_type, start_date, target_date, sequence_id, sort_order, completed_at, archived_at, is_draft, created_by, updated_by, created_at, updated_at`,
		id, params.StateID, params.Name, params.DescriptionHTML, params.DescriptionJSON,
		params.DescriptionStripped, params.Priority, params.IssueType, params.StartDate, params.TargetDate,
		params.ParentID, params.SortOrder, params.IsDraft, params.CompletedAt, params.ArchivedAt, params.UpdatedBy,
	).Scan(&i.ID, &i.ProjectID, &i.WorkspaceID, &i.ParentID, &i.StateID, &i.Name,
		&i.DescriptionHTML, &i.DescriptionJSON, &i.DescriptionStripped,
		&i.Priority, &i.IssueType, &i.StartDate, &i.TargetDate, &i.SequenceID, &i.SortOrder,
		&i.CompletedAt, &i.ArchivedAt, &i.IsDraft, &i.CreatedBy, &i.UpdatedBy, &i.CreatedAt, &i.UpdatedAt)
	return &i, err
}

func (r *Repository) DeleteIssue(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE issues SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}

// --- Issue Assignees ---

func (r *Repository) AddIssueAssignee(ctx context.Context, issueID, assigneeID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO issue_assignees (issue_id, assignee_id) VALUES ($1, $2) ON CONFLICT (issue_id, assignee_id) DO NOTHING`, issueID, assigneeID)
	return err
}

func (r *Repository) RemoveIssueAssignee(ctx context.Context, issueID, assigneeID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM issue_assignees WHERE issue_id = $1 AND assignee_id = $2`, issueID, assigneeID)
	return err
}

func (r *Repository) ListIssueAssignees(ctx context.Context, issueID uuid.UUID) ([]*api.UserSummary, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT u.id, u.email, u.first_name, u.last_name, u.display_name, u.avatar_url
		 FROM issue_assignees ia JOIN users u ON ia.assignee_id = u.id
		 WHERE ia.issue_id = $1`, issueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*api.UserSummary
	for rows.Next() {
		var u api.UserSummary
		if err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.DisplayName, &u.AvatarURL); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	if users == nil {
		users = []*api.UserSummary{}
	}
	return users, nil
}

// --- Issue Labels ---

func (r *Repository) AddIssueLabel(ctx context.Context, issueID, labelID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO issue_labels (issue_id, label_id) VALUES ($1, $2) ON CONFLICT (issue_id, label_id) DO NOTHING`, issueID, labelID)
	return err
}

func (r *Repository) RemoveIssueLabel(ctx context.Context, issueID, labelID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM issue_labels WHERE issue_id = $1 AND label_id = $2`, issueID, labelID)
	return err
}

func (r *Repository) ListIssueLabels(ctx context.Context, issueID uuid.UUID) ([]*api.Label, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT l.id, l.project_id, l.workspace_id, l.parent_id, l.name, l.description, l.color, l.sort_order, l.created_at, l.updated_at
		 FROM issue_labels il JOIN labels l ON il.label_id = l.id
		 WHERE il.issue_id = $1 AND l.deleted_at IS NULL`, issueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labels []*api.Label
	for rows.Next() {
		var l api.Label
		if err := rows.Scan(&l.ID, &l.ProjectID, &l.WorkspaceID, &l.ParentID, &l.Name, &l.Description, &l.Color, &l.SortOrder, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		labels = append(labels, &l)
	}
	if labels == nil {
		labels = []*api.Label{}
	}
	return labels, nil
}

// --- Issue Comments ---

func (r *Repository) CreateIssueComment(ctx context.Context, issueID, workspaceID uuid.UUID, commentHTML string, commentJSON []byte, commentStripped string, actorID *uuid.UUID) (*api.IssueComment, error) {
	var c api.IssueComment
	err := r.pool.QueryRow(ctx,
		`INSERT INTO issue_comments (issue_id, workspace_id, comment_html, comment_json, comment_stripped, actor_id)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, issue_id, workspace_id, comment_html, comment_json, comment_stripped, actor_id, created_at, updated_at`,
		issueID, workspaceID, commentHTML, commentJSON, commentStripped, actorID,
	).Scan(&c.ID, &c.IssueID, &c.WorkspaceID, &c.CommentHTML, &c.CommentJSON, &c.CommentStripped, &c.ActorID, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

func (r *Repository) GetIssueCommentByID(ctx context.Context, id uuid.UUID) (*api.IssueComment, error) {
	var c api.IssueComment
	err := r.pool.QueryRow(ctx,
		`SELECT id, issue_id, workspace_id, comment_html, comment_json, comment_stripped, actor_id, created_at, updated_at
		 FROM issue_comments WHERE id = $1 AND deleted_at IS NULL`, id,
	).Scan(&c.ID, &c.IssueID, &c.WorkspaceID, &c.CommentHTML, &c.CommentJSON, &c.CommentStripped, &c.ActorID, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

func (r *Repository) ListIssueComments(ctx context.Context, issueID uuid.UUID) ([]*api.IssueCommentWithUser, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT ic.id, ic.issue_id, ic.workspace_id, ic.comment_html, ic.comment_json, ic.comment_stripped, ic.actor_id, ic.created_at, ic.updated_at,
		        u.email, u.first_name, u.last_name, u.display_name, u.avatar_url
		 FROM issue_comments ic LEFT JOIN users u ON ic.actor_id = u.id
		 WHERE ic.issue_id = $1 AND ic.deleted_at IS NULL
		 ORDER BY ic.created_at ASC`, issueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*api.IssueCommentWithUser
	for rows.Next() {
		var c api.IssueCommentWithUser
		if err := rows.Scan(&c.ID, &c.IssueID, &c.WorkspaceID, &c.CommentHTML, &c.CommentJSON, &c.CommentStripped, &c.ActorID, &c.CreatedAt, &c.UpdatedAt,
			&c.Email, &c.FirstName, &c.LastName, &c.DisplayName, &c.AvatarURL); err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}
	if comments == nil {
		comments = []*api.IssueCommentWithUser{}
	}
	return comments, nil
}

func (r *Repository) UpdateIssueComment(ctx context.Context, id uuid.UUID, commentHTML string, commentJSON []byte, commentStripped string) (*api.IssueComment, error) {
	var c api.IssueComment
	err := r.pool.QueryRow(ctx,
		`UPDATE issue_comments SET comment_html = $2, comment_json = $3, comment_stripped = $4
		 WHERE id = $1 AND deleted_at IS NULL
		 RETURNING id, issue_id, workspace_id, comment_html, comment_json, comment_stripped, actor_id, created_at, updated_at`,
		id, commentHTML, commentJSON, commentStripped,
	).Scan(&c.ID, &c.IssueID, &c.WorkspaceID, &c.CommentHTML, &c.CommentJSON, &c.CommentStripped, &c.ActorID, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

func (r *Repository) DeleteIssueComment(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE issue_comments SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}

// --- Issue Activities ---

func (r *Repository) CreateIssueActivity(ctx context.Context, params api.CreateActivityParams) (*api.IssueActivity, error) {
	var a api.IssueActivity
	err := r.pool.QueryRow(ctx,
		`INSERT INTO issue_activities (issue_id, project_id, workspace_id, verb, field, old_value, new_value, old_identifier, new_identifier, comment, actor_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		 RETURNING id, issue_id, project_id, workspace_id, verb, field, old_value, new_value, old_identifier, new_identifier, comment, actor_id, created_at`,
		params.IssueID, params.ProjectID, params.WorkspaceID, params.Verb, params.Field,
		params.OldValue, params.NewValue, params.OldIdentifier, params.NewIdentifier, params.Comment, params.ActorID,
	).Scan(&a.ID, &a.IssueID, &a.ProjectID, &a.WorkspaceID, &a.Verb, &a.Field,
		&a.OldValue, &a.NewValue, &a.OldIdentifier, &a.NewIdentifier, &a.Comment, &a.ActorID, &a.CreatedAt)
	return &a, err
}

func (r *Repository) ListIssueActivities(ctx context.Context, issueID uuid.UUID) ([]*api.IssueActivityWithUser, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT ia.id, ia.issue_id, ia.project_id, ia.workspace_id, ia.verb, ia.field, ia.old_value, ia.new_value,
		        ia.old_identifier, ia.new_identifier, ia.comment, ia.actor_id, ia.created_at,
		        u.email, u.first_name, u.last_name, u.display_name, u.avatar_url
		 FROM issue_activities ia LEFT JOIN users u ON ia.actor_id = u.id
		 WHERE ia.issue_id = $1 ORDER BY ia.created_at ASC`, issueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []*api.IssueActivityWithUser
	for rows.Next() {
		var a api.IssueActivityWithUser
		if err := rows.Scan(&a.ID, &a.IssueID, &a.ProjectID, &a.WorkspaceID, &a.Verb, &a.Field,
			&a.OldValue, &a.NewValue, &a.OldIdentifier, &a.NewIdentifier, &a.Comment, &a.ActorID, &a.CreatedAt,
			&a.Email, &a.FirstName, &a.LastName, &a.DisplayName, &a.AvatarURL); err != nil {
			return nil, err
		}
		activities = append(activities, &a)
	}
	if activities == nil {
		activities = []*api.IssueActivityWithUser{}
	}
	return activities, nil
}

// --- Cycles ---

func (r *Repository) CreateCycle(ctx context.Context, projectID, workspaceID uuid.UUID, name, description string, startDate, endDate *time.Time, ownedBy uuid.UUID) (*api.Cycle, error) {
	var c api.Cycle
	err := r.pool.QueryRow(ctx,
		`INSERT INTO cycles (project_id, workspace_id, name, description, start_date, end_date, owned_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, project_id, workspace_id, name, description, start_date, end_date, owned_by, sort_order, progress_snapshot, archived_at, created_at, updated_at`,
		projectID, workspaceID, name, description, startDate, endDate, ownedBy,
	).Scan(&c.ID, &c.ProjectID, &c.WorkspaceID, &c.Name, &c.Description, &c.StartDate, &c.EndDate, &c.OwnedBy, &c.SortOrder, &c.ProgressSnapshot, &c.ArchivedAt, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

func (r *Repository) GetCycleByID(ctx context.Context, id uuid.UUID) (*api.Cycle, error) {
	var c api.Cycle
	err := r.pool.QueryRow(ctx,
		`SELECT id, project_id, workspace_id, name, description, start_date, end_date, owned_by, sort_order, progress_snapshot, archived_at, created_at, updated_at
		 FROM cycles WHERE id = $1 AND deleted_at IS NULL`, id,
	).Scan(&c.ID, &c.ProjectID, &c.WorkspaceID, &c.Name, &c.Description, &c.StartDate, &c.EndDate, &c.OwnedBy, &c.SortOrder, &c.ProgressSnapshot, &c.ArchivedAt, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

func (r *Repository) ListCyclesByProject(ctx context.Context, projectID uuid.UUID) ([]*api.Cycle, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, project_id, workspace_id, name, description, start_date, end_date, owned_by, sort_order, progress_snapshot, archived_at, created_at, updated_at
		 FROM cycles WHERE project_id = $1 AND deleted_at IS NULL ORDER BY start_date DESC NULLS LAST`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cycles []*api.Cycle
	for rows.Next() {
		var c api.Cycle
		if err := rows.Scan(&c.ID, &c.ProjectID, &c.WorkspaceID, &c.Name, &c.Description, &c.StartDate, &c.EndDate, &c.OwnedBy, &c.SortOrder, &c.ProgressSnapshot, &c.ArchivedAt, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		cycles = append(cycles, &c)
	}
	if cycles == nil {
		cycles = []*api.Cycle{}
	}
	return cycles, nil
}

func (r *Repository) UpdateCycle(ctx context.Context, id uuid.UUID, name, description *string, startDate, endDate *time.Time, sortOrder *float64) (*api.Cycle, error) {
	var c api.Cycle
	err := r.pool.QueryRow(ctx,
		`UPDATE cycles SET name = COALESCE($2, name), description = COALESCE($3, description), start_date = $4, end_date = $5, sort_order = COALESCE($6, sort_order)
		 WHERE id = $1 AND deleted_at IS NULL
		 RETURNING id, project_id, workspace_id, name, description, start_date, end_date, owned_by, sort_order, progress_snapshot, archived_at, created_at, updated_at`,
		id, name, description, startDate, endDate, sortOrder,
	).Scan(&c.ID, &c.ProjectID, &c.WorkspaceID, &c.Name, &c.Description, &c.StartDate, &c.EndDate, &c.OwnedBy, &c.SortOrder, &c.ProgressSnapshot, &c.ArchivedAt, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

func (r *Repository) DeleteCycle(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE cycles SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *Repository) AddIssueToCycle(ctx context.Context, cycleID, issueID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO cycle_issues (cycle_id, issue_id) VALUES ($1, $2) ON CONFLICT (cycle_id, issue_id) DO NOTHING`, cycleID, issueID)
	return err
}

func (r *Repository) RemoveIssueFromCycle(ctx context.Context, cycleID, issueID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM cycle_issues WHERE cycle_id = $1 AND issue_id = $2`, cycleID, issueID)
	return err
}

// --- Modules ---

func (r *Repository) CreateModule(ctx context.Context, projectID, workspaceID uuid.UUID, name, description string, startDate, targetDate *time.Time, status string, leadID, createdBy *uuid.UUID) (*api.Module, error) {
	var m api.Module
	err := r.pool.QueryRow(ctx,
		`INSERT INTO modules (project_id, workspace_id, name, description, start_date, target_date, status, lead_id, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id, project_id, workspace_id, name, description, description_html, start_date, target_date, status, lead_id, sort_order, archived_at, created_by, created_at, updated_at`,
		projectID, workspaceID, name, description, startDate, targetDate, status, leadID, createdBy,
	).Scan(&m.ID, &m.ProjectID, &m.WorkspaceID, &m.Name, &m.Description, &m.DescriptionHTML, &m.StartDate, &m.TargetDate, &m.Status, &m.LeadID, &m.SortOrder, &m.ArchivedAt, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt)
	return &m, err
}

func (r *Repository) GetModuleByID(ctx context.Context, id uuid.UUID) (*api.Module, error) {
	var m api.Module
	err := r.pool.QueryRow(ctx,
		`SELECT id, project_id, workspace_id, name, description, description_html, start_date, target_date, status, lead_id, sort_order, archived_at, created_by, created_at, updated_at
		 FROM modules WHERE id = $1 AND deleted_at IS NULL`, id,
	).Scan(&m.ID, &m.ProjectID, &m.WorkspaceID, &m.Name, &m.Description, &m.DescriptionHTML, &m.StartDate, &m.TargetDate, &m.Status, &m.LeadID, &m.SortOrder, &m.ArchivedAt, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt)
	return &m, err
}

func (r *Repository) ListModulesByProject(ctx context.Context, projectID uuid.UUID) ([]*api.Module, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, project_id, workspace_id, name, description, description_html, start_date, target_date, status, lead_id, sort_order, archived_at, created_by, created_at, updated_at
		 FROM modules WHERE project_id = $1 AND deleted_at IS NULL ORDER BY sort_order ASC, created_at DESC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modules []*api.Module
	for rows.Next() {
		var m api.Module
		if err := rows.Scan(&m.ID, &m.ProjectID, &m.WorkspaceID, &m.Name, &m.Description, &m.DescriptionHTML, &m.StartDate, &m.TargetDate, &m.Status, &m.LeadID, &m.SortOrder, &m.ArchivedAt, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		modules = append(modules, &m)
	}
	if modules == nil {
		modules = []*api.Module{}
	}
	return modules, nil
}

func (r *Repository) UpdateModule(ctx context.Context, id uuid.UUID, name, description *string, startDate, targetDate *time.Time, status *string, leadID *uuid.UUID, sortOrder *float64) (*api.Module, error) {
	var m api.Module
	err := r.pool.QueryRow(ctx,
		`UPDATE modules SET name = COALESCE($2, name), description = COALESCE($3, description),
		 start_date = $4, target_date = $5, status = COALESCE($6, status), lead_id = $7, sort_order = COALESCE($8, sort_order)
		 WHERE id = $1 AND deleted_at IS NULL
		 RETURNING id, project_id, workspace_id, name, description, description_html, start_date, target_date, status, lead_id, sort_order, archived_at, created_by, created_at, updated_at`,
		id, name, description, startDate, targetDate, status, leadID, sortOrder,
	).Scan(&m.ID, &m.ProjectID, &m.WorkspaceID, &m.Name, &m.Description, &m.DescriptionHTML, &m.StartDate, &m.TargetDate, &m.Status, &m.LeadID, &m.SortOrder, &m.ArchivedAt, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt)
	return &m, err
}

func (r *Repository) DeleteModule(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE modules SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *Repository) AddIssueToModule(ctx context.Context, moduleID, issueID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO module_issues (module_id, issue_id) VALUES ($1, $2) ON CONFLICT (module_id, issue_id) DO NOTHING`, moduleID, issueID)
	return err
}

func (r *Repository) RemoveIssueFromModule(ctx context.Context, moduleID, issueID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM module_issues WHERE module_id = $1 AND issue_id = $2`, moduleID, issueID)
	return err
}

// --- Pages ---

func (r *Repository) CreatePage(ctx context.Context, projectID *uuid.UUID, workspaceID uuid.UUID, name, descHTML string, descJSON []byte, descStripped string, ownedBy uuid.UUID, parentID *uuid.UUID) (*api.Page, error) {
	var p api.Page
	err := r.pool.QueryRow(ctx,
		`INSERT INTO pages (project_id, workspace_id, name, description_html, description_json, description_stripped, owned_by, parent_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, project_id, workspace_id, name, description_html, description_json, description_stripped, color, is_locked, archived_at, owned_by, parent_id, created_at, updated_at`,
		projectID, workspaceID, name, descHTML, descJSON, descStripped, ownedBy, parentID,
	).Scan(&p.ID, &p.ProjectID, &p.WorkspaceID, &p.Name, &p.DescriptionHTML, &p.DescriptionJSON, &p.DescriptionStripped, &p.Color, &p.IsLocked, &p.ArchivedAt, &p.OwnedBy, &p.ParentID, &p.CreatedAt, &p.UpdatedAt)
	return &p, err
}

func (r *Repository) GetPageByID(ctx context.Context, id uuid.UUID) (*api.Page, error) {
	var p api.Page
	err := r.pool.QueryRow(ctx,
		`SELECT id, project_id, workspace_id, name, description_html, description_json, description_stripped, color, is_locked, archived_at, owned_by, parent_id, created_at, updated_at
		 FROM pages WHERE id = $1 AND deleted_at IS NULL`, id,
	).Scan(&p.ID, &p.ProjectID, &p.WorkspaceID, &p.Name, &p.DescriptionHTML, &p.DescriptionJSON, &p.DescriptionStripped, &p.Color, &p.IsLocked, &p.ArchivedAt, &p.OwnedBy, &p.ParentID, &p.CreatedAt, &p.UpdatedAt)
	return &p, err
}

func (r *Repository) ListPagesByProject(ctx context.Context, projectID uuid.UUID) ([]*api.Page, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, project_id, workspace_id, name, description_html, description_json, description_stripped, color, is_locked, archived_at, owned_by, parent_id, created_at, updated_at
		 FROM pages WHERE project_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []*api.Page
	for rows.Next() {
		var p api.Page
		if err := rows.Scan(&p.ID, &p.ProjectID, &p.WorkspaceID, &p.Name, &p.DescriptionHTML, &p.DescriptionJSON, &p.DescriptionStripped, &p.Color, &p.IsLocked, &p.ArchivedAt, &p.OwnedBy, &p.ParentID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		pages = append(pages, &p)
	}
	if pages == nil {
		pages = []*api.Page{}
	}
	return pages, nil
}

func (r *Repository) UpdatePage(ctx context.Context, id uuid.UUID, name, descHTML *string, descJSON []byte, descStripped *string, color *string, isLocked *bool, parentID *uuid.UUID) (*api.Page, error) {
	var p api.Page
	err := r.pool.QueryRow(ctx,
		`UPDATE pages SET name = COALESCE($2, name), description_html = COALESCE($3, description_html),
		 description_json = COALESCE($4, description_json), description_stripped = COALESCE($5, description_stripped),
		 color = COALESCE($6, color), is_locked = COALESCE($7, is_locked), parent_id = $8
		 WHERE id = $1 AND deleted_at IS NULL
		 RETURNING id, project_id, workspace_id, name, description_html, description_json, description_stripped, color, is_locked, archived_at, owned_by, parent_id, created_at, updated_at`,
		id, name, descHTML, descJSON, descStripped, color, isLocked, parentID,
	).Scan(&p.ID, &p.ProjectID, &p.WorkspaceID, &p.Name, &p.DescriptionHTML, &p.DescriptionJSON, &p.DescriptionStripped, &p.Color, &p.IsLocked, &p.ArchivedAt, &p.OwnedBy, &p.ParentID, &p.CreatedAt, &p.UpdatedAt)
	return &p, err
}

func (r *Repository) DeletePage(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE pages SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}

// --- Notifications ---

func (r *Repository) ListNotificationsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*api.Notification, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, workspace_id, project_id, title, message, data, entity_type, entity_id, sender_id, receiver_id, read_at, snoozed_till, archived_at, created_at
		 FROM notifications WHERE receiver_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*api.Notification
	for rows.Next() {
		var n api.Notification
		if err := rows.Scan(&n.ID, &n.WorkspaceID, &n.ProjectID, &n.Title, &n.Message, &n.Data, &n.EntityType, &n.EntityID, &n.SenderID, &n.ReceiverID, &n.ReadAt, &n.SnoozedTill, &n.ArchivedAt, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, &n)
	}
	if notifications == nil {
		notifications = []*api.Notification{}
	}
	return notifications, nil
}

func (r *Repository) CountUnreadNotifications(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM notifications WHERE receiver_id = $1 AND read_at IS NULL`, userID).Scan(&count)
	return count, err
}

func (r *Repository) MarkNotificationRead(ctx context.Context, id, receiverID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE notifications SET read_at = NOW() WHERE id = $1 AND receiver_id = $2`, id, receiverID)
	return err
}

func (r *Repository) MarkAllNotificationsRead(ctx context.Context, receiverID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE notifications SET read_at = NOW() WHERE receiver_id = $1 AND read_at IS NULL`, receiverID)
	return err
}
