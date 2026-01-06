package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"unsri-backend/internal/file-storage/repository"
	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
	"unsri-backend/internal/shared/storage"

	"github.com/google/uuid"
)

// FileStorageService handles file storage business logic
type FileStorageService struct {
	repo      *repository.FileRepository
	config    StorageConfig
	minioClient *storage.MinIOClient
}

// StorageConfig holds storage configuration
type StorageConfig struct {
	Type     string
	BasePath string
	BaseURL  string
	MaxSize  int64
	MinIO    MinIOStorageConfig
}

// MinIOStorageConfig holds MinIO storage configuration
type MinIOStorageConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	Region          string
	BucketProfiles  string
	BucketSelfies   string
	BucketDocuments string
}

// NewFileStorageService creates a new file storage service
func NewFileStorageService(repo *repository.FileRepository, config StorageConfig, minioClient *storage.MinIOClient) *FileStorageService {
	return &FileStorageService{
		repo:        repo,
		config:      config,
		minioClient: minioClient,
	}
}

// UploadFileRequest represents upload file request
type UploadFileRequest struct {
	File     *multipart.FileHeader
	FileType string
	IsPublic bool
}

// UploadFile uploads a file
func (s *FileStorageService) UploadFile(ctx context.Context, userID string, req UploadFileRequest) (*models.File, error) {
	// Validate file size
	if req.File.Size > s.config.MaxSize {
		return nil, apperrors.NewValidationError(fmt.Sprintf("file size exceeds maximum allowed size of %d bytes", s.config.MaxSize))
	}

	// Generate unique filename
	ext := filepath.Ext(req.File.Filename)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	objectName := fmt.Sprintf("%s/%s", req.FileType, fileName)

	// Open file
	src, err := req.File.Open()
	if err != nil {
		return nil, apperrors.NewInternalError("failed to open file", err)
	}
	defer src.Close()

	var fileURL string
	var filePath string

	// Upload based on storage type
	if s.config.Type == "minio" && s.minioClient != nil {
		// Upload to MinIO
		bucketName := s.getBucketForFileType(req.FileType)
		if err := s.minioClient.UploadFile(ctx, bucketName, objectName, src, req.File.Size, req.File.Header.Get("Content-Type")); err != nil {
			return nil, apperrors.NewInternalError("failed to upload file to MinIO", err)
		}

		// Generate public URL
		fileURL = s.minioClient.GetPublicURL(bucketName, objectName)
		filePath = fmt.Sprintf("%s/%s", bucketName, objectName)
	} else {
		// Fallback to local storage
		// Create directory if not exists
		fileDir := filepath.Join(s.config.BasePath, req.FileType)
		if err := os.MkdirAll(fileDir, 0755); err != nil {
			return nil, apperrors.NewInternalError("failed to create directory", err)
		}

		// Save file
		filePath = filepath.Join(fileDir, fileName)
		dst, err := os.Create(filePath)
		if err != nil {
			return nil, apperrors.NewInternalError("failed to create file", err)
		}
		defer dst.Close()

		// Reset reader position
		src.Seek(0, 0)

		if _, err := io.Copy(dst, src); err != nil {
			return nil, apperrors.NewInternalError("failed to save file", err)
		}

		fileURL = fmt.Sprintf("%s/%s/%s", s.config.BaseURL, req.FileType, fileName)
	}

	// Create file record
	file := &models.File{
		UserID:       userID,
		FileName:     fileName,
		OriginalName: req.File.Filename,
		FileType:     req.FileType,
		MimeType:     req.File.Header.Get("Content-Type"),
		Size:         req.File.Size,
		Path:         filePath,
		URL:          fileURL,
		IsPublic:     req.IsPublic,
	}

	if err := s.repo.CreateFile(ctx, file); err != nil {
		// Cleanup on error
		if s.config.Type == "minio" && s.minioClient != nil {
			bucketName := s.getBucketForFileType(req.FileType)
			_ = s.minioClient.DeleteFile(ctx, bucketName, objectName)
		} else {
			_ = os.Remove(filePath)
		}
		return nil, apperrors.NewInternalError("failed to create file record", err)
	}

	return file, nil
}

// getBucketForFileType returns the appropriate bucket name for a file type
func (s *FileStorageService) getBucketForFileType(fileType string) string {
	switch fileType {
	case "avatar", "profile":
		return s.config.MinIO.BucketProfiles
	case "selfie":
		return s.config.MinIO.BucketSelfies
	default:
		return s.config.MinIO.BucketDocuments
	}
}

// GetFileByID gets a file by ID
func (s *FileStorageService) GetFileByID(ctx context.Context, id string) (*models.File, error) {
	return s.repo.GetFileByID(ctx, id)
}

// GetFilesRequest represents get files request
type GetFilesRequest struct {
	FileType string `form:"file_type"`
	Page     int    `form:"page,default=1"`
	PerPage  int    `form:"per_page,default=20"`
}

// GetFiles gets files for a user
func (s *FileStorageService) GetFiles(ctx context.Context, userID string, req GetFilesRequest) ([]models.File, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	var fileTypePtr *string
	if req.FileType != "" {
		fileTypePtr = &req.FileType
	}

	return s.repo.GetFilesByUserID(ctx, userID, fileTypePtr, perPage, (page-1)*perPage)
}

// DeleteFile deletes a file
func (s *FileStorageService) DeleteFile(ctx context.Context, id string, userID string) error {
	file, err := s.repo.GetFileByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFoundError("file", id)
	}

	// Check ownership
	if file.UserID != userID {
		return apperrors.NewForbiddenError("not authorized to delete this file")
	}

	// Delete physical file
	if s.config.Type == "minio" && s.minioClient != nil {
		bucketName := s.getBucketForFileType(file.FileType)
		objectName := fmt.Sprintf("%s/%s", file.FileType, file.FileName)
		if err := s.minioClient.DeleteFile(ctx, bucketName, objectName); err != nil {
			return apperrors.NewInternalError("failed to delete file from MinIO", err)
		}
	} else {
		if err := os.Remove(file.Path); err != nil && !os.IsNotExist(err) {
			return apperrors.NewInternalError("failed to delete file", err)
		}
	}

	// Delete record
	return s.repo.DeleteFile(ctx, id)
}

// UploadAvatarRequest represents upload avatar request
type UploadAvatarRequest struct {
	File *multipart.FileHeader
}

// UploadAvatar uploads user avatar
func (s *FileStorageService) UploadAvatar(ctx context.Context, userID string, req UploadAvatarRequest) (*models.File, error) {
	// Delete old avatar
	oldAvatar, _ := s.repo.GetAvatarByUserID(ctx, userID)
	if oldAvatar != nil {
		if s.config.Type == "minio" && s.minioClient != nil {
			// Extract bucket and object name from path
			bucketName := s.getBucketForFileType("avatar")
			objectName := fmt.Sprintf("avatar/%s", oldAvatar.FileName)
			_ = s.minioClient.DeleteFile(ctx, bucketName, objectName)
		} else {
			_ = os.Remove(oldAvatar.Path)
		}
		_ = s.repo.DeleteFile(ctx, oldAvatar.ID)
	}

	uploadReq := UploadFileRequest{
		File:     req.File,
		FileType: "avatar",
		IsPublic: true,
	}

	return s.UploadFile(ctx, userID, uploadReq)
}

// UploadDocumentRequest represents upload document request
type UploadDocumentRequest struct {
	File     *multipart.FileHeader
	IsPublic bool
}

// UploadDocument uploads a document
func (s *FileStorageService) UploadDocument(ctx context.Context, userID string, req UploadDocumentRequest) (*models.File, error) {
	uploadReq := UploadFileRequest{
		File:     req.File,
		FileType: "document",
		IsPublic: req.IsPublic,
	}

	return s.UploadFile(ctx, userID, uploadReq)
}

// GetFileContent gets file content for download
func (s *FileStorageService) GetFileContent(ctx context.Context, id string) ([]byte, string, error) {
	file, err := s.repo.GetFileByID(ctx, id)
	if err != nil {
		return nil, "", apperrors.NewNotFoundError("file", id)
	}

	if s.config.Type == "minio" && s.minioClient != nil {
		// For MinIO, return URL instead of content
		// Client should fetch from URL directly
		return nil, "", apperrors.NewBadRequestError("file stored in MinIO, use URL to access")
	}

	content, err := os.ReadFile(file.Path)
	if err != nil {
		return nil, "", apperrors.NewInternalError("failed to read file", err)
	}

	return content, file.MimeType, nil
}

// UploadToMinIO uploads a file directly to MinIO (helper method)
func (s *FileStorageService) UploadToMinIO(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64, contentType string) error {
	if s.minioClient == nil {
		return apperrors.NewInternalError("MinIO client not initialized", nil)
	}
	return s.minioClient.UploadFile(ctx, bucketName, objectName, reader, size, contentType)
}

// GetMinIOURL generates a URL for a file in MinIO
func (s *FileStorageService) GetMinIOURL(ctx context.Context, bucketName, objectName string, expires time.Duration) (string, error) {
	if s.minioClient == nil {
		return "", apperrors.NewInternalError("MinIO client not initialized", nil)
	}
	if expires == 0 {
		expires = 7 * 24 * time.Hour // Default 7 days
	}
	return s.minioClient.GetFileURL(ctx, bucketName, objectName, expires)
}
