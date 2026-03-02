package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/open-pm/open-pm/server/internal/api"
)

func (r *Repository) CreateView(ctx context.Context, projectID *uuid.UUID, workspaceID uuid.UUID, name, description string, filters, displayFilters, displayProperties []byte, access int16, createdBy *uuid.UUID) (*api.View, error) {
	var v api.View
	err := r.pool.QueryRow(ctx,
		`INSERT INTO views (project_id, workspace_id, name, description, filters, display_filters, display_properties, access, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id, project_id, workspace_id, name, description, filters, display_filters, display_properties, access, sort_order, created_by, created_at, updated_at`,
		projectID, workspaceID, name, description, filters, displayFilters, displayProperties, access, createdBy,
	).Scan(&v.ID, &v.ProjectID, &v.WorkspaceID, &v.Name, &v.Description, &v.Filters, &v.DisplayFilters, &v.DisplayProperties, &v.Access, &v.SortOrder, &v.CreatedBy, &v.CreatedAt, &v.UpdatedAt)
	return &v, err
}

func (r *Repository) GetViewByID(ctx context.Context, id uuid.UUID) (*api.View, error) {
	var v api.View
	err := r.pool.QueryRow(ctx,
		`SELECT id, project_id, workspace_id, name, description, filters, display_filters, display_properties, access, sort_order, created_by, created_at, updated_at
		 FROM views WHERE id = $1 AND deleted_at IS NULL`, id,
	).Scan(&v.ID, &v.ProjectID, &v.WorkspaceID, &v.Name, &v.Description, &v.Filters, &v.DisplayFilters, &v.DisplayProperties, &v.Access, &v.SortOrder, &v.CreatedBy, &v.CreatedAt, &v.UpdatedAt)
	return &v, err
}

func (r *Repository) ListViewsByProject(ctx context.Context, projectID uuid.UUID) ([]*api.View, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, project_id, workspace_id, name, description, filters, display_filters, display_properties, access, sort_order, created_by, created_at, updated_at
		 FROM views WHERE project_id = $1 AND deleted_at IS NULL
		 ORDER BY sort_order ASC, created_at DESC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []*api.View
	for rows.Next() {
		var v api.View
		if err := rows.Scan(&v.ID, &v.ProjectID, &v.WorkspaceID, &v.Name, &v.Description, &v.Filters, &v.DisplayFilters, &v.DisplayProperties, &v.Access, &v.SortOrder, &v.CreatedBy, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		views = append(views, &v)
	}
	if views == nil {
		views = []*api.View{}
	}
	return views, nil
}

func (r *Repository) UpdateView(ctx context.Context, id uuid.UUID, name, description *string, filters, displayFilters, displayProperties []byte, access *int16) (*api.View, error) {
	var v api.View
	err := r.pool.QueryRow(ctx,
		`UPDATE views SET
			name = COALESCE($2, name),
			description = COALESCE($3, description),
			filters = COALESCE($4, filters),
			display_filters = COALESCE($5, display_filters),
			display_properties = COALESCE($6, display_properties),
			access = COALESCE($7, access),
			updated_at = NOW()
		 WHERE id = $1 AND deleted_at IS NULL
		 RETURNING id, project_id, workspace_id, name, description, filters, display_filters, display_properties, access, sort_order, created_by, created_at, updated_at`,
		id, name, description, filters, displayFilters, displayProperties, access,
	).Scan(&v.ID, &v.ProjectID, &v.WorkspaceID, &v.Name, &v.Description, &v.Filters, &v.DisplayFilters, &v.DisplayProperties, &v.Access, &v.SortOrder, &v.CreatedBy, &v.CreatedAt, &v.UpdatedAt)
	return &v, err
}

func (r *Repository) DeleteView(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE views SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}
