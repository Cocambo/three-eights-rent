package storage

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"car-service/config"

	"github.com/minio/minio-go/v7"
)

type StoredObject struct {
	BucketName string
	ObjectKey  string
}

type ImageStorageService interface {
	UploadObject(
		ctx context.Context,
		bucketName, objectKey string,
		reader io.Reader,
		size int64,
		contentType string,
	) (StoredObject, error)
	RemoveObject(ctx context.Context, bucketName, objectKey string) error
	GeneratePresignedGetURL(ctx context.Context, bucketName, objectKey string) (string, error)
}

type minioImageStorageService struct {
	client        *minio.Client
	presignClient *minio.Client
	defaultBucket string
	presignTTL    time.Duration
}

func NewMinIOImageStorageService(
	client *minio.Client,
	presignClient *minio.Client,
	cfg config.MinIOConfig,
) ImageStorageService {
	if presignClient == nil {
		presignClient = client
	}

	return &minioImageStorageService{
		client:        client,
		presignClient: presignClient,
		defaultBucket: strings.TrimSpace(cfg.DefaultBucket),
		presignTTL:    cfg.PresignTTL,
	}
}

func (s *minioImageStorageService) UploadObject(
	ctx context.Context,
	bucketName, objectKey string,
	reader io.Reader,
	size int64,
	contentType string,
) (StoredObject, error) {
	resolvedBucket, resolvedObjectKey, err := s.resolveLocation(bucketName, objectKey)
	if err != nil {
		return StoredObject{}, err
	}

	if size < 0 {
		return StoredObject{}, fmt.Errorf("object size must be greater than or equal to zero")
	}

	_, err = s.client.PutObject(
		ctx,
		resolvedBucket,
		resolvedObjectKey,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType: strings.TrimSpace(contentType),
		},
	)
	if err != nil {
		return StoredObject{}, fmt.Errorf("put object %s/%s: %w", resolvedBucket, resolvedObjectKey, err)
	}

	return StoredObject{
		BucketName: resolvedBucket,
		ObjectKey:  resolvedObjectKey,
	}, nil
}

func (s *minioImageStorageService) RemoveObject(
	ctx context.Context,
	bucketName, objectKey string,
) error {
	resolvedBucket, resolvedObjectKey, err := s.resolveLocation(bucketName, objectKey)
	if err != nil {
		return err
	}

	if err := s.client.RemoveObject(
		ctx,
		resolvedBucket,
		resolvedObjectKey,
		minio.RemoveObjectOptions{},
	); err != nil {
		return fmt.Errorf("remove object %s/%s: %w", resolvedBucket, resolvedObjectKey, err)
	}

	return nil
}

func (s *minioImageStorageService) GeneratePresignedGetURL(
	ctx context.Context,
	bucketName, objectKey string,
) (string, error) {
	resolvedBucket, resolvedObjectKey, err := s.resolveLocation(bucketName, objectKey)
	if err != nil {
		return "", err
	}

	presignedURL, err := s.presignClient.PresignedGetObject(
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

func (s *minioImageStorageService) resolveLocation(
	bucketName, objectKey string,
) (string, string, error) {
	resolvedBucket := strings.TrimSpace(bucketName)
	if resolvedBucket == "" {
		resolvedBucket = s.defaultBucket
	}

	resolvedObjectKey := strings.TrimSpace(objectKey)
	if resolvedBucket == "" {
		return "", "", fmt.Errorf("bucket name is required")
	}

	if resolvedObjectKey == "" {
		return "", "", fmt.Errorf("object key is required")
	}

	return resolvedBucket, resolvedObjectKey, nil
}
