package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOConfig holds MinIO configuration
type MinIOConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	Region          string
}

// MinIOClient wraps MinIO client with helper methods
type MinIOClient struct {
	client *minio.Client
	config MinIOConfig
}

// NewMinIOClient creates a new MinIO client
func NewMinIOClient(config MinIOConfig) (*MinIOClient, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
		Region: config.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	return &MinIOClient{
		client: client,
		config: config,
	}, nil
}

// CreateBucketIfNotExists creates a bucket if it doesn't exist
func (m *MinIOClient) CreateBucketIfNotExists(ctx context.Context, bucketName string) error {
	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
			Region: m.config.Region,
		})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}

		// Set bucket policy to public read (optional, adjust as needed)
		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::` + bucketName + `/*"]
				}
			]
		}`
		err = m.client.SetBucketPolicy(ctx, bucketName, policy)
		if err != nil {
			// Log warning but don't fail - policy setting is optional
			_ = err
		}
	}

	return nil
}

// UploadFile uploads a file to MinIO
func (m *MinIOClient) UploadFile(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64, contentType string) error {
	_, err := m.client.PutObject(ctx, bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

// GetFileURL generates a presigned URL for a file (expires in 7 days by default)
func (m *MinIOClient) GetFileURL(ctx context.Context, bucketName, objectName string, expires time.Duration) (string, error) {
	if expires == 0 {
		expires = 7 * 24 * time.Hour // Default 7 days
	}

	url, err := m.client.PresignedGetObject(ctx, bucketName, objectName, expires, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}

// GetPublicURL generates a public URL for a file (if bucket is public)
func (m *MinIOClient) GetPublicURL(bucketName, objectName string) string {
	protocol := "http"
	if m.config.UseSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, m.config.Endpoint, bucketName, objectName)
}

// DeleteFile deletes a file from MinIO
func (m *MinIOClient) DeleteFile(ctx context.Context, bucketName, objectName string) error {
	err := m.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// FileExists checks if a file exists in MinIO
func (m *MinIOClient) FileExists(ctx context.Context, bucketName, objectName string) (bool, error) {
	_, err := m.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}

	return true, nil
}

// GetClient returns the underlying MinIO client (for advanced usage)
func (m *MinIOClient) GetClient() *minio.Client {
	return m.client
}

