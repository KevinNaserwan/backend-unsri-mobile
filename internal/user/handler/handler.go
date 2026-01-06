package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/utils"
	"unsri-backend/internal/user/service"
	apperrors "unsri-backend/internal/shared/errors"

	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests for user management
type UserHandler struct {
	service *service.UserService
	logger  logger.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(service *service.UserService, logger logger.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// GetProfile handles get user profile request
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	result, err := h.service.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// UpdateProfile handles update user profile request
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	var req service.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.UpdateUserProfile(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// GetUserByID handles get user by ID request
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	result, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// SearchUsers handles search users request
func (h *UserHandler) SearchUsers(c *gin.Context) {
	var req service.SearchUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	users, total, err := h.service.SearchUsers(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	utils.PaginatedResponse(c, users, page, perPage, total)
}

// GetMahasiswaByNIM handles get mahasiswa by NIM request
func (h *UserHandler) GetMahasiswaByNIM(c *gin.Context) {
	nim := c.Param("nim")

	result, err := h.service.GetMahasiswaByNIM(c.Request.Context(), nim)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// GetDosenByNIP handles get dosen by NIP request
func (h *UserHandler) GetDosenByNIP(c *gin.Context) {
	nip := c.Param("nip")

	result, err := h.service.GetDosenByNIP(c.Request.Context(), nip)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// GetStaffByNIP handles get staff by NIP request
func (h *UserHandler) GetStaffByNIP(c *gin.Context) {
	nip := c.Param("nip")

	result, err := h.service.GetStaffByNIP(c.Request.Context(), nip)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// UploadAvatar handles avatar upload request
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID := c.GetString("user_id")

	file, err := c.FormFile("photo")
	if err != nil {
		// Try alternative field name
		file, err = c.FormFile("avatar")
		if err != nil {
			utils.BadRequestResponse(c, "Photo file is required (use 'photo' or 'avatar' field)")
			return
		}
	}

	// Upload file to file storage service
	photoURL, err := h.uploadPhotoToStorage(c.Request.Context(), userID, file)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	req := service.UploadAvatarRequest{
		PhotoURL: photoURL,
	}

	avatarURL, err := h.service.UploadAvatar(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"profile_photo_url": avatarURL,
		"message":           "Profile photo uploaded successfully",
	})
}

// ListUsers handles list all users request (admin only)
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req service.ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	users, total, err := h.service.ListUsers(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	utils.PaginatedResponse(c, users, page, perPage, total)
}

// AdminUpdateUser handles admin update user request
func (h *UserHandler) AdminUpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var req service.AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.AdminUpdateUser(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// DeleteUser handles delete user request (admin only)
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	err := h.service.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// ActivateUser handles activate user request (admin only)
func (h *UserHandler) ActivateUser(c *gin.Context) {
	userID := c.Param("id")

	result, err := h.service.ActivateUser(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// DeactivateUser handles deactivate user request (admin only)
func (h *UserHandler) DeactivateUser(c *gin.Context) {
	userID := c.Param("id")

	result, err := h.service.DeactivateUser(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// CreateUser handles admin create user request
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req service.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.CreateUser(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, result)
}

// uploadPhotoToStorage uploads photo to file storage service
func (h *UserHandler) uploadPhotoToStorage(ctx context.Context, userID string, file *multipart.FileHeader) (string, error) {
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", apperrors.NewInternalError("failed to open uploaded file", err)
	}
	defer src.Close()

	// Create multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add file field
	part, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		return "", apperrors.NewInternalError("failed to create form file", err)
	}
	
	if _, err := io.Copy(part, src); err != nil {
		return "", apperrors.NewInternalError("failed to copy file data", err)
	}

	// Add file_type field
	if err := writer.WriteField("file_type", "profile"); err != nil {
		return "", apperrors.NewInternalError("failed to write file_type field", err)
	}

	// Add is_public field
	if err := writer.WriteField("is_public", "true"); err != nil {
		return "", apperrors.NewInternalError("failed to write is_public field", err)
	}

	writer.Close()

	// Create HTTP request to file storage service
	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:8093/api/v1/files/upload", &buf)
	if err != nil {
		return "", apperrors.NewInternalError("failed to create upload request", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-ID", userID) // Pass user ID for authentication

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", apperrors.NewInternalError("failed to upload file", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", apperrors.NewInternalError("file upload failed", fmt.Errorf("status: %d", resp.StatusCode))
	}

	// Parse response
	var uploadResp struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return "", apperrors.NewInternalError("failed to parse upload response", err)
	}

	return uploadResp.Data.URL, nil
}
