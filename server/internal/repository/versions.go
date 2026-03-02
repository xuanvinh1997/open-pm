package repository

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/open-pm/open-pm/server/internal/api"
)

func (r *Repository) CreateVersion(ctx context.Context, projectID, workspaceID uuid.UUID, name, description string, startDate, releaseDate *time.Time, createdBy *uuid.UUID) (*api.Version, error) {
	var v api.Version
	err := r.pool.QueryRow(ctx,
		`INSERT INTO versions (project_id, workspace_id, name, description, start_date, release_date, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, project_id, workspace_id, name, description, start_date, release_date, released, released_at, archived_at, sort_order, created_by, created_at, updated_at`,
		projectID, workspaceID, name, description, startDate, releaseDate, createdBy,
	).Scan(&v.ID, &v.ProjectID, &v.WorkspaceID, &v.Name, &v.Description, &v.StartDate, &v.ReleaseDate,
		&v.Released, &v.ReleasedAt, &v.ArchivedAt, &v.SortOrder, &v.CreatedBy, &v.CreatedAt, &v.UpdatedAt)
	return &v, err
}

func (r *Repository) GetVersionByID(ctx context.Context, id uuid.UUID) (*api.Version, error) {
	var v api.Version
	err := r.pool.QueryRow(ctx,
		`SELECT id, project_id, workspace_id, name, description, start_date, release_date, released, released_at, archived_at, sort_order, created_by, created_at, updated_at
		 FROM versions WHERE id = $1 AND deleted_at IS NULL`, id,
	).Scan(&v.ID, &v.ProjectID, &v.WorkspaceID, &v.Name, &v.Description, &v.StartDate, &v.ReleaseDate,
		&v.Released, &v.ReleasedAt, &v.ArchivedAt, &v.SortOrder, &v.CreatedBy, &v.CreatedAt, &v.UpdatedAt)
	return &v, err
}

func (r *Repository) ListVersionsByProject(ctx context.Context, projectID uuid.UUID) ([]*api.Version, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, project_id, workspace_id, name, description, start_date, release_date, released, released_at, archived_at, sort_order, created_by, created_at, updated_at
		 FROM versions WHERE project_id = $1 AND deleted_at IS NULL
		 ORDER BY sort_order ASC, release_date ASC NULLS LAST`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []*api.Version
	for rows.Next() {
		var v api.Version
		if err := rows.Scan(&v.ID, &v.ProjectID, &v.WorkspaceID, &v.Name, &v.Description, &v.StartDate, &v.ReleaseDate,
			&v.Released, &v.ReleasedAt, &v.ArchivedAt, &v.SortOrder, &v.CreatedBy, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		versions = append(versions, &v)
	}
	if versions == nil {
		versions = []*api.Version{}
	}
	return versions, nil
}

func (r *Repository) UpdateVersion(ctx context.Context, id uuid.UUID, name, description *string, startDate, releaseDate *time.Time, released *bool, sortOrder *float64) (*api.Version, error) {
	var v api.Version
	err := r.pool.QueryRow(ctx,
		`UPDATE versions SET
			name = COALESCE($2, name),
			description = COALESCE($3, description),
			start_date = $4,
			release_date = $5,
			released = COALESCE($6, released),
			released_at = CASE WHEN $6 IS TRUE AND released_at IS NULL THEN NOW() ELSE released_at END,
			sort_order = COALESCE($7, sort_order)
		 WHERE id = $1 AND deleted_at IS NULL
		 RETURNING id, project_id, workspace_id, name, description, start_date, release_date, released, released_at, archived_at, sort_order, created_by, created_at, updated_at`,
		id, name, description, startDate, releaseDate, released, sortOrder,
	).Scan(&v.ID, &v.ProjectID, &v.WorkspaceID, &v.Name, &v.Description, &v.StartDate, &v.ReleaseDate,
		&v.Released, &v.ReleasedAt, &v.ArchivedAt, &v.SortOrder, &v.CreatedBy, &v.CreatedAt, &v.UpdatedAt)
	return &v, err
}

func (r *Repository) DeleteVersion(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE versions SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *Repository) AddIssueFixVersion(ctx context.Context, issueID, versionID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO issue_fix_versions (issue_id, version_id) VALUES ($1, $2) ON CONFLICT (issue_id, version_id) DO NOTHING`,
		issueID, versionID)
	return err
}

func (r *Repository) RemoveIssueFixVersion(ctx context.Context, issueID, versionID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM issue_fix_versions WHERE issue_id = $1 AND version_id = $2`, issueID, versionID)
	return err
}

func (r *Repository) AddIssueAffectsVersion(ctx context.Context, issueID, versionID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO issue_affects_versions (issue_id, version_id) VALUES ($1, $2) ON CONFLICT (issue_id, version_id) DO NOTHING`,
		issueID, versionID)
	return err
}

func (r *Repository) RemoveIssueAffectsVersion(ctx context.Context, issueID, versionID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM issue_affects_versions WHERE issue_id = $1 AND version_id = $2`, issueID, versionID)
	return err
}

func (r *Repository) ListFixVersionsByIssue(ctx context.Context, issueID uuid.UUID) ([]*api.Version, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT v.id, v.project_id, v.workspace_id, v.name, v.description, v.start_date, v.release_date, v.released, v.released_at, v.archived_at, v.sort_order, v.created_by, v.created_at, v.updated_at
		 FROM versions v JOIN issue_fix_versions ifv ON ifv.version_id = v.id
		 WHERE ifv.issue_id = $1 AND v.deleted_at IS NULL
		 ORDER BY v.release_date ASC NULLS LAST`, issueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []*api.Version
	for rows.Next() {
		var v api.Version
		if err := rows.Scan(&v.ID, &v.ProjectID, &v.WorkspaceID, &v.Name, &v.Description, &v.StartDate, &v.ReleaseDate,
			&v.Released, &v.ReleasedAt, &v.ArchivedAt, &v.SortOrder, &v.CreatedBy, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		versions = append(versions, &v)
	}
	if versions == nil {
		versions = []*api.Version{}
	}
	return versions, nil
}

func (r *Repository) ListAffectsVersionsByIssue(ctx context.Context, issueID uuid.UUID) ([]*api.Version, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT v.id, v.project_id, v.workspace_id, v.name, v.description, v.start_date, v.release_date, v.released, v.released_at, v.archived_at, v.sort_order, v.created_by, v.created_at, v.updated_at
		 FROM versions v JOIN issue_affects_versions iav ON iav.version_id = v.id
		 WHERE iav.issue_id = $1 AND v.deleted_at IS NULL
		 ORDER BY v.release_date ASC NULLS LAST`, issueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []*api.Version
	for rows.Next() {
		var v api.Version
		if err := rows.Scan(&v.ID, &v.ProjectID, &v.WorkspaceID, &v.Name, &v.Description, &v.StartDate, &v.ReleaseDate,
			&v.Released, &v.ReleasedAt, &v.ArchivedAt, &v.SortOrder, &v.CreatedBy, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		versions = append(versions, &v)
	}
	if versions == nil {
		versions = []*api.Version{}
	}
	return versions, nil
}
