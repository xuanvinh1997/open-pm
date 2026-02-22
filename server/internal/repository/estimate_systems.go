package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/open-pm/open-pm/server/internal/api"
)

func (r *Repository) CreateEstimateSystem(ctx context.Context, projectID, workspaceID uuid.UUID, systemType string, estimates []byte) (*api.EstimateSystem, error) {
	var es api.EstimateSystem
	err := r.pool.QueryRow(ctx,
		`INSERT INTO estimate_systems (project_id, workspace_id, type, estimates)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, project_id, workspace_id, type, estimates, created_at, updated_at`,
		projectID, workspaceID, systemType, estimates,
	).Scan(&es.ID, &es.ProjectID, &es.WorkspaceID, &es.Type, &es.Estimates, &es.CreatedAt, &es.UpdatedAt)
	return &es, err
}

func (r *Repository) GetEstimateSystemByProject(ctx context.Context, projectID uuid.UUID) (*api.EstimateSystem, error) {
	var es api.EstimateSystem
	err := r.pool.QueryRow(ctx,
		`SELECT id, project_id, workspace_id, type, estimates, created_at, updated_at
		 FROM estimate_systems WHERE project_id = $1 AND deleted_at IS NULL`, projectID,
	).Scan(&es.ID, &es.ProjectID, &es.WorkspaceID, &es.Type, &es.Estimates, &es.CreatedAt, &es.UpdatedAt)
	return &es, err
}

func (r *Repository) UpdateEstimateSystem(ctx context.Context, id uuid.UUID, systemType *string, estimates []byte) (*api.EstimateSystem, error) {
	var es api.EstimateSystem
	err := r.pool.QueryRow(ctx,
		`UPDATE estimate_systems SET
			type = COALESCE($2, type),
			estimates = COALESCE($3, estimates)
		 WHERE id = $1 AND deleted_at IS NULL
		 RETURNING id, project_id, workspace_id, type, estimates, created_at, updated_at`,
		id, systemType, estimates,
	).Scan(&es.ID, &es.ProjectID, &es.WorkspaceID, &es.Type, &es.Estimates, &es.CreatedAt, &es.UpdatedAt)
	return &es, err
}

func (r *Repository) DeleteEstimateSystem(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE estimate_systems SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}
