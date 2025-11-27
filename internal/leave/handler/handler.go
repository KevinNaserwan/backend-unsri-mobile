package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"unsri-backend/internal/leave/service"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/utils"
)

// LeaveHandler handles HTTP requests for leave management
type LeaveHandler struct {
	service *service.LeaveService
	logger  logger.Logger
}

// NewLeaveHandler creates a new leave handler
func NewLeaveHandler(service *service.LeaveService, logger logger.Logger) *LeaveHandler {
	return &LeaveHandler{
		service: service,
		logger:  logger,
	}
}

// ========== Leave Request Handlers ==========

// CreateLeaveRequest handles create leave request
func (h *LeaveHandler) CreateLeaveRequest(c *gin.Context) {
	userID := c.GetString("user_id")

	var req service.CreateLeaveRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.CreateLeaveRequest(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 201, result)
}

// GetLeaveRequest handles get leave request by ID request
func (h *LeaveHandler) GetLeaveRequest(c *gin.Context) {
	leaveID := c.Param("id")

	result, err := h.service.GetLeaveRequestByID(c.Request.Context(), leaveID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetLeaveRequests handles get leave requests request
func (h *LeaveHandler) GetLeaveRequests(c *gin.Context) {
	var req service.GetLeaveRequestsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	// If user_id not provided, use current user (for non-admin)
	userRole := c.GetString("user_role")
	if req.UserID == "" && userRole != "staff" {
		req.UserID = c.GetString("user_id")
	}

	leaveRequests, total, err := h.service.GetLeaveRequests(c.Request.Context(), req)
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

	utils.PaginatedResponse(c, leaveRequests, page, perPage, total)
}

// GetLeaveRequestsByUser handles get leave requests by user request
func (h *LeaveHandler) GetLeaveRequestsByUser(c *gin.Context) {
	userID := c.Param("userId")
	status := c.Query("status")

	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	result, err := h.service.GetLeaveRequestsByUser(c.Request.Context(), userID, statusPtr)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// ApproveLeaveRequest handles approve leave request
func (h *LeaveHandler) ApproveLeaveRequest(c *gin.Context) {
	leaveID := c.Param("id")
	approverID := c.GetString("user_id")

	var req service.ApproveLeaveRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.ApproveLeaveRequest(c.Request.Context(), leaveID, approverID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// RejectLeaveRequest handles reject leave request
func (h *LeaveHandler) RejectLeaveRequest(c *gin.Context) {
	leaveID := c.Param("id")
	approverID := c.GetString("user_id")

	var req service.RejectLeaveRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.RejectLeaveRequest(c.Request.Context(), leaveID, approverID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// CancelLeaveRequest handles cancel leave request
func (h *LeaveHandler) CancelLeaveRequest(c *gin.Context) {
	leaveID := c.Param("id")
	userID := c.GetString("user_id")

	result, err := h.service.CancelLeaveRequest(c.Request.Context(), leaveID, userID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// DeleteLeaveRequest handles delete leave request
func (h *LeaveHandler) DeleteLeaveRequest(c *gin.Context) {
	leaveID := c.Param("id")

	err := h.service.DeleteLeaveRequest(c.Request.Context(), leaveID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, gin.H{"message": "Leave request deleted successfully"})
}

// ========== Leave Quota Handlers ==========

// CreateLeaveQuota handles create leave quota request
func (h *LeaveHandler) CreateLeaveQuota(c *gin.Context) {
	var req service.CreateLeaveQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.CreateLeaveQuota(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 201, result)
}

// GetLeaveQuota handles get leave quota by ID request
func (h *LeaveHandler) GetLeaveQuota(c *gin.Context) {
	quotaID := c.Param("id")

	result, err := h.service.GetLeaveQuotaByID(c.Request.Context(), quotaID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetLeaveQuotas handles get leave quotas request
func (h *LeaveHandler) GetLeaveQuotas(c *gin.Context) {
	var req service.GetLeaveQuotasRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	quotas, total, err := h.service.GetLeaveQuotas(c.Request.Context(), req)
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

	utils.PaginatedResponse(c, quotas, page, perPage, total)
}

// GetLeaveQuotasByUser handles get leave quotas by user request
func (h *LeaveHandler) GetLeaveQuotasByUser(c *gin.Context) {
	userID := c.Param("userId")
	yearStr := c.Query("year")

	var yearPtr *int
	if yearStr != "" {
		year := 0
		if _, err := fmt.Sscanf(yearStr, "%d", &year); err == nil {
			yearPtr = &year
		}
	}

	result, err := h.service.GetLeaveQuotasByUser(c.Request.Context(), userID, yearPtr)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// UpdateLeaveQuota handles update leave quota request
func (h *LeaveHandler) UpdateLeaveQuota(c *gin.Context) {
	quotaID := c.Param("id")

	var req service.UpdateLeaveQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.UpdateLeaveQuota(c.Request.Context(), quotaID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// DeleteLeaveQuota handles delete leave quota request
func (h *LeaveHandler) DeleteLeaveQuota(c *gin.Context) {
	quotaID := c.Param("id")

	err := h.service.DeleteLeaveQuota(c.Request.Context(), quotaID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, gin.H{"message": "Leave quota deleted successfully"})
}

