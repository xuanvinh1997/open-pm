package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

const maxUploadSize = 10 << 20 // 10 MB

// UploadAsset handles POST .../workspaces/{slug}/assets
func (a *API) UploadAsset(w http.ResponseWriter, r *http.Request) error {
	if a.storage == nil {
		return internalServerError("file storage not configured")
	}

	ctx := r.Context()
	userID := getUserID(ctx)
	workspaceID := getWorkspaceID(ctx)

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		return badRequestError("file too large (max 10MB)")
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return badRequestError("missing file field")
	}
	defer file.Close()

	entityType := r.FormValue("entity_type")
	entityIDStr := r.FormValue("entity_id")
	if entityType == "" || entityIDStr == "" {
		return badRequestError("entity_type and entity_id are required")
	}
	entityID, err := uuid.FromString(entityIDStr)
	if err != nil {
		return badRequestError("invalid entity_id")
	}

	// Validate entity type
	validEntityTypes := map[string]bool{"issue": true, "page": true, "comment": true}
	if !validEntityTypes[entityType] {
		return badRequestError("invalid entity_type: must be issue, page, or comment")
	}

	// Generate storage key
	fileUUID, _ := uuid.NewV4()
	storageKey := fmt.Sprintf("%s/%s/%s-%s", workspaceID.String(), entityType, fileUUID.String(), header.Filename)
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload to storage
	if err := a.storage.Upload(ctx, storageKey, file, header.Size, contentType); err != nil {
		log.Error().Err(err).Str("key", storageKey).Msg("failed to upload file to storage")
		return internalServerError("failed to upload file")
	}

	// Create DB record
	asset, err := a.queries.CreateFileAsset(ctx, workspaceID, entityType, entityID, header.Filename, header.Size, contentType, storageKey, userID)
	if err != nil {
		log.Error().Err(err).Msg("failed to create file asset record")
		return internalServerError("failed to save file metadata")
	}

	// Generate download URL
	downloadURL, err := a.storage.PresignedGetURL(ctx, storageKey, 1*time.Hour)
	if err == nil {
		asset.DownloadURL = downloadURL
	}

	return sendJSON(w, http.StatusCreated, asset)
}

// ListAssetsByEntity handles GET .../workspaces/{slug}/assets?entity_type=x&entity_id=y
func (a *API) ListAssetsByEntity(w http.ResponseWriter, r *http.Request) error {
	entityType := r.URL.Query().Get("entity_type")
	entityIDStr := r.URL.Query().Get("entity_id")
	if entityType == "" || entityIDStr == "" {
		return badRequestError("entity_type and entity_id query params are required")
	}
	entityID, err := uuid.FromString(entityIDStr)
	if err != nil {
		return badRequestError("invalid entity_id")
	}

	assets, err := a.queries.ListFileAssetsByEntity(r.Context(), entityType, entityID)
	if err != nil {
		return internalServerError("failed to list assets")
	}

	// Generate download URLs
	if a.storage != nil {
		for _, asset := range assets {
			if url, err := a.storage.PresignedGetURL(r.Context(), asset.StorageKey, 1*time.Hour); err == nil {
				asset.DownloadURL = url
			}
		}
	}

	return sendJSON(w, http.StatusOK, map[string]interface{}{"results": assets})
}

// DeleteAsset handles DELETE .../workspaces/{slug}/assets/{assetID}
func (a *API) DeleteAsset(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)
	assetID, err := uuid.FromString(chi.URLParam(r, "assetID"))
	if err != nil {
		return badRequestError("invalid asset ID")
	}

	asset, err := a.queries.GetFileAssetByID(ctx, assetID)
	if err != nil {
		return notFoundError("asset not found")
	}

	// Only uploader or admin can delete
	if asset.UploadedBy != userID {
		return forbiddenError("only the uploader can delete this file")
	}

	// Delete from storage
	if a.storage != nil {
		if err := a.storage.Delete(ctx, asset.StorageKey); err != nil {
			log.Warn().Err(err).Str("key", asset.StorageKey).Msg("failed to delete file from storage")
		}
	}

	// Soft delete from DB
	if err := a.queries.DeleteFileAsset(ctx, assetID); err != nil {
		return internalServerError("failed to delete asset")
	}
	return sendEmpty(w, http.StatusNoContent)
}
