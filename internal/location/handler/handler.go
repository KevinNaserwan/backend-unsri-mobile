package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"unsri-backend/internal/location/service"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/utils"
)

// LocationHandler handles HTTP requests for location
type LocationHandler struct {
	service *service.LocationService
	logger  logger.Logger
}

// NewLocationHandler creates a new location handler
func NewLocationHandler(service *service.LocationService, logger logger.Logger) *LocationHandler {
	return &LocationHandler{
		service: service,
		logger:  logger,
	}
}

// GetCheckInStatus handles get check-in status request
func (h *LocationHandler) GetCheckInStatus(c *gin.Context) {
	userID := c.GetString("user_id")

	result, err := h.service.GetCheckInStatus(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// GetLocationHistory handles get location history request
func (h *LocationHandler) GetLocationHistory(c *gin.Context) {
	userID := c.GetString("user_id")

	var req service.GetLocationHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	history, total, err := h.service.GetLocationHistory(c.Request.Context(), userID, req)
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

	utils.PaginatedResponse(c, history, page, perPage, total)
}

// GetGeofences handles get geofences request
func (h *LocationHandler) GetGeofences(c *gin.Context) {
	result, err := h.service.GetGeofences(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// ValidateLocation handles validate location request
func (h *LocationHandler) ValidateLocation(c *gin.Context) {
	var req service.ValidateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.ValidateLocation(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// CreateGeofence handles create geofence request
func (h *LocationHandler) CreateGeofence(c *gin.Context) {
	var req service.CreateGeofenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.CreateGeofence(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, result)
}

// GetGeofence handles get geofence by id request
func (h *LocationHandler) GetGeofence(c *gin.Context) {
	id := c.Param("id")

	result, err := h.service.GetGeofence(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// UpdateGeofence handles update geofence request
func (h *LocationHandler) UpdateGeofence(c *gin.Context) {
	id := c.Param("id")

	var req service.UpdateGeofenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.UpdateGeofence(c.Request.Context(), id, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// DeleteGeofence handles delete geofence request
func (h *LocationHandler) DeleteGeofence(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteGeofence(c.Request.Context(), id); err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "geofence deleted"})
}
