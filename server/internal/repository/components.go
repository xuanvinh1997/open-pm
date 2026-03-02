package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/open-pm/open-pm/server/internal/api"
)

func (r *Repository) CreateComponent(ctx context.Context, projectID, workspaceID uuid.UUID, name, description string, leadID *uuid.UUID, defaultAssigneeType string, sortOrder float64) (*api.Component, error) {
	var c api.Component
	err := r.pool.QueryRow(ctx,
		`INSERT INTO components (project_id, workspace_id, name, description, lead_id, default_assignee_type, sort_order)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, project_id, workspace_id, name, description, lead_id, default_assignee_type, sort_order, created_at, updated_at`,
		projectID, workspaceID, name, description, leadID, defaultAssigneeType, sortOrder,
	).Scan(&c.ID, &c.ProjectID, &c.WorkspaceID, &c.Name, &c.Description, &c.LeadID, &c.DefaultAssigneeType, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

func (r *Repository) GetComponentByID(ctx context.Context, id uuid.UUID) (*api.Component, error) {
	var c api.Component
	err := r.pool.QueryRow(ctx,
		`SELECT id, project_id, workspace_id, name, description, lead_id, default_assignee_type, sort_order, created_at, updated_at
		 FROM components WHERE id = $1 AND deleted_at IS NULL`, id,
	).Scan(&c.ID, &c.ProjectID, &c.WorkspaceID, &c.Name, &c.Description, &c.LeadID, &c.DefaultAssigneeType, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

func (r *Repository) ListComponentsByProject(ctx context.Context, projectID uuid.UUID) ([]*api.ComponentWithLead, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT c.id, c.project_id, c.workspace_id, c.name, c.description, c.lead_id, c.default_assignee_type, c.sort_order, c.created_at, c.updated_at,
		        u.email, u.first_name, u.last_name, u.display_name, u.avatar_url
		 FROM components c
		 LEFT JOIN users u ON c.lead_id = u.id
		 WHERE c.project_id = $1 AND c.deleted_at IS NULL
		 ORDER BY c.sort_order ASC, c.name ASC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var components []*api.ComponentWithLead
	for rows.Next() {
		var c api.ComponentWithLead
		if err := rows.Scan(&c.ID, &c.ProjectID, &c.WorkspaceID, &c.Name, &c.Description, &c.LeadID, &c.DefaultAssigneeType, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt,
			&c.LeadEmail, &c.LeadFirstName, &c.LeadLastName, &c.LeadDisplayName, &c.LeadAvatarURL); err != nil {
			return nil, err
		}
		components = append(components, &c)
	}
	if components == nil {
		components = []*api.ComponentWithLead{}
	}
	return components, nil
}

func (r *Repository) UpdateComponent(ctx context.Context, id uuid.UUID, name, description *string, leadID *uuid.UUID, defaultAssigneeType *string, sortOrder *float64) (*api.Component, error) {
	var c api.Component
	err := r.pool.QueryRow(ctx,
		`UPDATE components SET
			name = COALESCE($2, name),
			description = COALESCE($3, description),
			lead_id = $4,
			default_assignee_type = COALESCE($5, default_assignee_type),
			sort_order = COALESCE($6, sort_order)
		 WHERE id = $1 AND deleted_at IS NULL
		 RETURNING id, project_id, workspace_id, name, description, lead_id, default_assignee_type, sort_order, created_at, updated_at`,
		id, name, description, leadID, defaultAssigneeType, sortOrder,
	).Scan(&c.ID, &c.ProjectID, &c.WorkspaceID, &c.Name, &c.Description, &c.LeadID, &c.DefaultAssigneeType, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

func (r *Repository) DeleteComponent(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE components SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *Repository) AddIssueComponent(ctx context.Context, issueID, componentID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO issue_components (issue_id, component_id) VALUES ($1, $2) ON CONFLICT (issue_id, component_id) DO NOTHING`,
		issueID, componentID)
	return err
}

func (r *Repository) RemoveIssueComponent(ctx context.Context, issueID, componentID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM issue_components WHERE issue_id = $1 AND component_id = $2`, issueID, componentID)
	return err
}

func (r *Repository) ListComponentsByIssue(ctx context.Context, issueID uuid.UUID) ([]*api.Component, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT c.id, c.project_id, c.workspace_id, c.name, c.description, c.lead_id, c.default_assignee_type, c.sort_order, c.created_at, c.updated_at
		 FROM components c JOIN issue_components ic ON ic.component_id = c.id
		 WHERE ic.issue_id = $1 AND c.deleted_at IS NULL
		 ORDER BY c.name ASC`, issueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var components []*api.Component
	for rows.Next() {
		var c api.Component
		if err := rows.Scan(&c.ID, &c.ProjectID, &c.WorkspaceID, &c.Name, &c.Description, &c.LeadID, &c.DefaultAssigneeType, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		components = append(components, &c)
	}
	if components == nil {
		components = []*api.Component{}
	}
	return components, nil
}
