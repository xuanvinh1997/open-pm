package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/open-pm/open-pm/server/internal/api"
)

func (r *Repository) CreateIssueLink(ctx context.Context, issueID uuid.UUID, title, url string, createdBy *uuid.UUID) (*api.IssueLink, error) {
	var l api.IssueLink
	err := r.pool.QueryRow(ctx,
		`INSERT INTO issue_links (issue_id, title, url, created_by)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, issue_id, title, url, created_by, created_at, updated_at`,
		issueID, title, url, createdBy,
	).Scan(&l.ID, &l.IssueID, &l.Title, &l.URL, &l.CreatedBy, &l.CreatedAt, &l.UpdatedAt)
	return &l, err
}

func (r *Repository) ListIssueLinksByIssue(ctx context.Context, issueID uuid.UUID) ([]*api.IssueLink, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, issue_id, title, url, created_by, created_at, updated_at
		 FROM issue_links WHERE issue_id = $1 AND deleted_at IS NULL
		 ORDER BY created_at DESC`, issueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*api.IssueLink
	for rows.Next() {
		var l api.IssueLink
		if err := rows.Scan(&l.ID, &l.IssueID, &l.Title, &l.URL, &l.CreatedBy, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		links = append(links, &l)
	}
	if links == nil {
		links = []*api.IssueLink{}
	}
	return links, nil
}

func (r *Repository) UpdateIssueLink(ctx context.Context, id uuid.UUID, title, url *string) (*api.IssueLink, error) {
	var l api.IssueLink
	err := r.pool.QueryRow(ctx,
		`UPDATE issue_links SET title = COALESCE($2, title), url = COALESCE($3, url)
		 WHERE id = $1 AND deleted_at IS NULL
		 RETURNING id, issue_id, title, url, created_by, created_at, updated_at`,
		id, title, url,
	).Scan(&l.ID, &l.IssueID, &l.Title, &l.URL, &l.CreatedBy, &l.CreatedAt, &l.UpdatedAt)
	return &l, err
}

func (r *Repository) DeleteIssueLink(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE issue_links SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}
