package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"unsri-backend/internal/master-data/service"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/utils"
)

// MasterDataHandler handles HTTP requests for master data management
type MasterDataHandler struct {
	service *service.MasterDataService
	logger  logger.Logger
}

// NewMasterDataHandler creates a new master data handler
func NewMasterDataHandler(service *service.MasterDataService, logger logger.Logger) *MasterDataHandler {
	return &MasterDataHandler{
		service: service,
		logger:  logger,
	}
}

// ========== Study Program Handlers ==========

// CreateStudyProgram handles create study program request
func (h *MasterDataHandler) CreateStudyProgram(c *gin.Context) {
	var req service.CreateStudyProgramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.CreateStudyProgram(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, result)
}

// GetStudyProgram handles get study program by ID request
func (h *MasterDataHandler) GetStudyProgram(c *gin.Context) {
	studyProgramID := c.Param("id")

	result, err := h.service.GetStudyProgramByID(c.Request.Context(), studyProgramID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// GetStudyPrograms handles get study programs request
func (h *MasterDataHandler) GetStudyPrograms(c *gin.Context) {
	var req service.GetStudyProgramsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	studyPrograms, total, err := h.service.GetStudyPrograms(c.Request.Context(), req)
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

	utils.PaginatedResponse(c, studyPrograms, page, perPage, total)
}

// UpdateStudyProgram handles update study program request
func (h *MasterDataHandler) UpdateStudyProgram(c *gin.Context) {
	studyProgramID := c.Param("id")

	var req service.UpdateStudyProgramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.UpdateStudyProgram(c.Request.Context(), studyProgramID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// DeleteStudyProgram handles delete study program request
func (h *MasterDataHandler) DeleteStudyProgram(c *gin.Context) {
	studyProgramID := c.Param("id")

	err := h.service.DeleteStudyProgram(c.Request.Context(), studyProgramID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Study program deleted successfully"})
}

// ========== Academic Period Handlers ==========

// CreateAcademicPeriod handles create academic period request
func (h *MasterDataHandler) CreateAcademicPeriod(c *gin.Context) {
	var req service.CreateAcademicPeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.CreateAcademicPeriod(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, result)
}

// GetAcademicPeriod handles get academic period by ID request
func (h *MasterDataHandler) GetAcademicPeriod(c *gin.Context) {
	academicPeriodID := c.Param("id")

	result, err := h.service.GetAcademicPeriodByID(c.Request.Context(), academicPeriodID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// GetActiveAcademicPeriod handles get active academic period request
func (h *MasterDataHandler) GetActiveAcademicPeriod(c *gin.Context) {
	result, err := h.service.GetActiveAcademicPeriod(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// GetAcademicPeriods handles get academic periods request
func (h *MasterDataHandler) GetAcademicPeriods(c *gin.Context) {
	var req service.GetAcademicPeriodsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	academicPeriods, total, err := h.service.GetAcademicPeriods(c.Request.Context(), req)
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

	utils.PaginatedResponse(c, academicPeriods, page, perPage, total)
}

// UpdateAcademicPeriod handles update academic period request
func (h *MasterDataHandler) UpdateAcademicPeriod(c *gin.Context) {
	academicPeriodID := c.Param("id")

	var req service.UpdateAcademicPeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.UpdateAcademicPeriod(c.Request.Context(), academicPeriodID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// DeleteAcademicPeriod handles delete academic period request
func (h *MasterDataHandler) DeleteAcademicPeriod(c *gin.Context) {
	academicPeriodID := c.Param("id")

	err := h.service.DeleteAcademicPeriod(c.Request.Context(), academicPeriodID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, gin.H{"message": "Academic period deleted successfully"})
}

// ========== Room Handlers ==========

// CreateRoom handles create room request
func (h *MasterDataHandler) CreateRoom(c *gin.Context) {
	var req service.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.CreateRoom(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, result)
}

// GetRoom handles get room by ID request
func (h *MasterDataHandler) GetRoom(c *gin.Context) {
	roomID := c.Param("id")

	result, err := h.service.GetRoomByID(c.Request.Context(), roomID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// GetRooms handles get rooms request
func (h *MasterDataHandler) GetRooms(c *gin.Context) {
	var req service.GetRoomsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	rooms, total, err := h.service.GetRooms(c.Request.Context(), req)
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

	utils.PaginatedResponse(c, rooms, page, perPage, total)
}

// UpdateRoom handles update room request
func (h *MasterDataHandler) UpdateRoom(c *gin.Context) {
	roomID := c.Param("id")

	var req service.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.UpdateRoom(c.Request.Context(), roomID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// DeleteRoom handles delete room request
func (h *MasterDataHandler) DeleteRoom(c *gin.Context) {
	roomID := c.Param("id")

	err := h.service.DeleteRoom(c.Request.Context(), roomID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, gin.H{"message": "Room deleted successfully"})
}

