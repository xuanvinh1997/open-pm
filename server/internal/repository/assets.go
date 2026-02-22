package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/open-pm/open-pm/server/internal/api"
)

func (r *Repository) CreateFileAsset(ctx context.Context, workspaceID uuid.UUID, entityType string, entityID uuid.UUID, fileName string, fileSize int64, contentType, storageKey string, uploadedBy uuid.UUID) (*api.FileAsset, error) {
	var a api.FileAsset
	err := r.pool.QueryRow(ctx,
		`INSERT INTO file_assets (workspace_id, entity_type, entity_id, file_name, file_size, content_type, storage_key, uploaded_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, workspace_id, entity_type, entity_id, file_name, file_size, content_type, storage_key, uploaded_by, created_at`,
		workspaceID, entityType, entityID, fileName, fileSize, contentType, storageKey, uploadedBy,
	).Scan(&a.ID, &a.WorkspaceID, &a.EntityType, &a.EntityID, &a.FileName, &a.FileSize, &a.ContentType, &a.StorageKey, &a.UploadedBy, &a.CreatedAt)
	return &a, err
}

func (r *Repository) GetFileAssetByID(ctx context.Context, id uuid.UUID) (*api.FileAsset, error) {
	var a api.FileAsset
	err := r.pool.QueryRow(ctx,
		`SELECT id, workspace_id, entity_type, entity_id, file_name, file_size, content_type, storage_key, uploaded_by, created_at
		 FROM file_assets WHERE id = $1 AND deleted_at IS NULL`, id,
	).Scan(&a.ID, &a.WorkspaceID, &a.EntityType, &a.EntityID, &a.FileName, &a.FileSize, &a.ContentType, &a.StorageKey, &a.UploadedBy, &a.CreatedAt)
	return &a, err
}

func (r *Repository) ListFileAssetsByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*api.FileAsset, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, workspace_id, entity_type, entity_id, file_name, file_size, content_type, storage_key, uploaded_by, created_at
		 FROM file_assets WHERE entity_type = $1 AND entity_id = $2 AND deleted_at IS NULL
		 ORDER BY created_at DESC`, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []*api.FileAsset
	for rows.Next() {
		var a api.FileAsset
		if err := rows.Scan(&a.ID, &a.WorkspaceID, &a.EntityType, &a.EntityID, &a.FileName, &a.FileSize, &a.ContentType, &a.StorageKey, &a.UploadedBy, &a.CreatedAt); err != nil {
			return nil, err
		}
		assets = append(assets, &a)
	}
	if assets == nil {
		assets = []*api.FileAsset{}
	}
	return assets, nil
}

func (r *Repository) DeleteFileAsset(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE file_assets SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}
