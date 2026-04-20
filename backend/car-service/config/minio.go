package config

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinIOClient(cfg *Config) (*minio.Client, error) {
	return newMinIOClient(cfg.MinIO.Endpoint, cfg.MinIO.UseSSL, cfg.MinIO.AccessKey, cfg.MinIO.SecretKey)
}

func NewMinIOPublicClient(cfg *Config) (*minio.Client, error) {
	return minio.New(cfg.MinIO.PublicEndpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(cfg.MinIO.AccessKey, cfg.MinIO.SecretKey, ""),
		Secure:       cfg.MinIO.PublicUseSSL,
		Region:       "us-east-1",
		BucketLookup: minio.BucketLookupPath,
	})
}

func newMinIOClient(endpoint string, useSSL bool, accessKey string, secretKey string) (*minio.Client, error) {
	return minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
}
