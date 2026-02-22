package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/open-pm/open-pm/server/internal/api"
)

func (r *Repository) AddIssueRelation(ctx context.Context, issueID, relatedIssueID uuid.UUID, relationType string, createdBy *uuid.UUID) (*api.IssueRelation, error) {
	var rel api.IssueRelation
	err := r.pool.QueryRow(ctx,
		`INSERT INTO issue_relations (issue_id, related_issue_id, relation_type, created_by)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (issue_id, related_issue_id, relation_type) DO NOTHING
		 RETURNING id, issue_id, related_issue_id, relation_type, created_by, created_at`,
		issueID, relatedIssueID, relationType, createdBy,
	).Scan(&rel.ID, &rel.IssueID, &rel.RelatedIssueID, &rel.RelationType, &rel.CreatedBy, &rel.CreatedAt)
	return &rel, err
}

func (r *Repository) GetIssueRelationByID(ctx context.Context, id uuid.UUID) (*api.IssueRelation, error) {
	var rel api.IssueRelation
	err := r.pool.QueryRow(ctx,
		`SELECT id, issue_id, related_issue_id, relation_type, created_by, created_at
		 FROM issue_relations WHERE id = $1`, id,
	).Scan(&rel.ID, &rel.IssueID, &rel.RelatedIssueID, &rel.RelationType, &rel.CreatedBy, &rel.CreatedAt)
	return &rel, err
}

func (r *Repository) RemoveIssueRelation(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM issue_relations WHERE id = $1`, id)
	return err
}

func (r *Repository) RemoveInverseIssueRelation(ctx context.Context, issueID, relatedIssueID uuid.UUID, relationType string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM issue_relations WHERE issue_id = $1 AND related_issue_id = $2 AND relation_type = $3`,
		issueID, relatedIssueID, relationType)
	return err
}

func (r *Repository) ListIssueRelations(ctx context.Context, issueID uuid.UUID) ([]*api.IssueRelationWithIssue, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT ir.id, ir.issue_id, ir.related_issue_id, ir.relation_type, ir.created_by, ir.created_at,
		        i.name, i.sequence_id, i.priority, i.state_id
		 FROM issue_relations ir
		 JOIN issues i ON ir.related_issue_id = i.id AND i.deleted_at IS NULL
		 WHERE ir.issue_id = $1
		 ORDER BY ir.created_at DESC`, issueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var relations []*api.IssueRelationWithIssue
	for rows.Next() {
		var rel api.IssueRelationWithIssue
		if err := rows.Scan(
			&rel.ID, &rel.IssueID, &rel.RelatedIssueID, &rel.RelationType, &rel.CreatedBy, &rel.CreatedAt,
			&rel.RelatedIssueName, &rel.RelatedIssueSequenceID, &rel.RelatedIssuePriority, &rel.RelatedIssueStateID,
		); err != nil {
			return nil, err
		}
		relations = append(relations, &rel)
	}
	if relations == nil {
		relations = []*api.IssueRelationWithIssue{}
	}
	return relations, nil
}
