package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/open-pm/open-pm/server/internal/api"
)

func (r *Repository) CreateIssueResolution(ctx context.Context, projectID, workspaceID uuid.UUID, name, description string, isDefault bool, sortOrder float64) (*api.IssueResolution, error) {
	var res api.IssueResolution
	err := r.pool.QueryRow(ctx,
		`INSERT INTO issue_resolutions (project_id, workspace_id, name, description, is_default, sort_order)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, project_id, workspace_id, name, description, is_default, sort_order, created_at`,
		projectID, workspaceID, name, description, isDefault, sortOrder,
	).Scan(&res.ID, &res.ProjectID, &res.WorkspaceID, &res.Name, &res.Description, &res.IsDefault, &res.SortOrder, &res.CreatedAt)
	return &res, err
}

func (r *Repository) GetIssueResolutionByID(ctx context.Context, id uuid.UUID) (*api.IssueResolution, error) {
	var res api.IssueResolution
	err := r.pool.QueryRow(ctx,
		`SELECT id, project_id, workspace_id, name, description, is_default, sort_order, created_at
		 FROM issue_resolutions WHERE id = $1`, id,
	).Scan(&res.ID, &res.ProjectID, &res.WorkspaceID, &res.Name, &res.Description, &res.IsDefault, &res.SortOrder, &res.CreatedAt)
	return &res, err
}

func (r *Repository) ListIssueResolutionsByProject(ctx context.Context, projectID uuid.UUID) ([]*api.IssueResolution, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, project_id, workspace_id, name, description, is_default, sort_order, created_at
		 FROM issue_resolutions WHERE project_id = $1
		 ORDER BY sort_order ASC, created_at ASC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*api.IssueResolution
	for rows.Next() {
		var res api.IssueResolution
		if err := rows.Scan(&res.ID, &res.ProjectID, &res.WorkspaceID, &res.Name, &res.Description, &res.IsDefault, &res.SortOrder, &res.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, &res)
	}
	if results == nil {
		results = []*api.IssueResolution{}
	}
	return results, nil
}

func (r *Repository) UpdateIssueResolution(ctx context.Context, id uuid.UUID, name, description *string, isDefault *bool, sortOrder *float64) (*api.IssueResolution, error) {
	var res api.IssueResolution
	err := r.pool.QueryRow(ctx,
		`UPDATE issue_resolutions SET
			name = COALESCE($2, name),
			description = COALESCE($3, description),
			is_default = COALESCE($4, is_default),
			sort_order = COALESCE($5, sort_order)
		 WHERE id = $1
		 RETURNING id, project_id, workspace_id, name, description, is_default, sort_order, created_at`,
		id, name, description, isDefault, sortOrder,
	).Scan(&res.ID, &res.ProjectID, &res.WorkspaceID, &res.Name, &res.Description, &res.IsDefault, &res.SortOrder, &res.CreatedAt)
	return &res, err
}

func (r *Repository) DeleteIssueResolution(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM issue_resolutions WHERE id = $1`, id)
	return err
}
