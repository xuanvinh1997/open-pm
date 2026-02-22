package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/open-pm/open-pm/server/internal/config"
	"github.com/rs/zerolog/log"
)

// Storage wraps MinIO for file operations.
type Storage struct {
	client *minio.Client
	bucket string
}

// New creates a Storage client and ensures the bucket exists.
func New(cfg config.StorageConfig) (*Storage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("storage: failed to create minio client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("storage: failed to check bucket: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("storage: failed to create bucket: %w", err)
		}
		log.Info().Str("bucket", cfg.Bucket).Msg("created storage bucket")
	}

	return &Storage{client: client, bucket: cfg.Bucket}, nil
}

// Upload stores a file in the bucket.
func (s *Storage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	_, err := s.client.PutObject(ctx, s.bucket, key, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// Delete removes a file from the bucket.
func (s *Storage) Delete(ctx context.Context, key string) error {
	return s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{})
}

// PresignedGetURL returns a temporary download URL valid for the given duration.
func (s *Storage) PresignedGetURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	reqParams := make(url.Values)
	u, err := s.client.PresignedGetObject(ctx, s.bucket, key, expires, reqParams)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
