package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"car-service/config"

	"github.com/minio/minio-go/v7"
)

type ImageStorageService interface {
	GeneratePresignedGetURL(ctx context.Context, bucketName, objectKey string) (string, error)
}

type minioImageStorageService struct {
	client        *minio.Client
	defaultBucket string
	presignTTL    time.Duration
}

func NewMinIOImageStorageService(client *minio.Client, cfg config.MinIOConfig) ImageStorageService {
	return &minioImageStorageService{
		client:        client,
		defaultBucket: strings.TrimSpace(cfg.DefaultBucket),
		presignTTL:    cfg.PresignTTL,
	}
}

func (s *minioImageStorageService) GeneratePresignedGetURL(
	ctx context.Context,
	bucketName, objectKey string,
) (string, error) {
	resolvedBucket := strings.TrimSpace(bucketName)
	if resolvedBucket == "" {
		resolvedBucket = s.defaultBucket
	}

	resolvedObjectKey := strings.TrimSpace(objectKey)
	if resolvedBucket == "" {
		return "", fmt.Errorf("bucket name is required")
	}

	if resolvedObjectKey == "" {
		return "", fmt.Errorf("object key is required")
	}

	presignedURL, err := s.client.PresignedGetObject(
		ctx,
		resolvedBucket,
		resolvedObjectKey,
		s.presignTTL,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("presign get object %s/%s: %w", resolvedBucket, resolvedObjectKey, err)
	}

	return presignedURL.String(), nil
}
