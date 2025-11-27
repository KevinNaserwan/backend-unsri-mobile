package handler

import (
	"github.com/gin-gonic/gin"
	"unsri-backend/internal/course/service"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/utils"
)

// CourseHandler handles HTTP requests for course management
type CourseHandler struct {
	service *service.CourseService
	logger  logger.Logger
}

// NewCourseHandler creates a new course handler
func NewCourseHandler(service *service.CourseService, logger logger.Logger) *CourseHandler {
	return &CourseHandler{
		service: service,
		logger:  logger,
	}
}

// CreateCourse handles create course request
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var req service.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.CreateCourse(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 201, result)
}

// GetCourse handles get course by ID request
func (h *CourseHandler) GetCourse(c *gin.Context) {
	courseID := c.Param("id")

	result, err := h.service.GetCourseByID(c.Request.Context(), courseID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetCourses handles get courses request
func (h *CourseHandler) GetCourses(c *gin.Context) {
	var req service.GetCoursesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	courses, total, err := h.service.GetCourses(c.Request.Context(), req)
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

	utils.PaginatedResponse(c, courses, page, perPage, total)
}

// UpdateCourse handles update course request
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	courseID := c.Param("id")

	var req service.UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.UpdateCourse(c.Request.Context(), courseID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// DeleteCourse handles delete course request
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	courseID := c.Param("id")

	err := h.service.DeleteCourse(c.Request.Context(), courseID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, gin.H{"message": "Course deleted successfully"})
}

// CreateClass handles create class request
func (h *CourseHandler) CreateClass(c *gin.Context) {
	var req service.CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.CreateClass(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 201, result)
}

// GetClass handles get class by ID request
func (h *CourseHandler) GetClass(c *gin.Context) {
	classID := c.Param("id")

	result, err := h.service.GetClassByID(c.Request.Context(), classID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetClasses handles get classes request
func (h *CourseHandler) GetClasses(c *gin.Context) {
	var req service.GetClassesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	classes, total, err := h.service.GetClasses(c.Request.Context(), req)
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

	utils.PaginatedResponse(c, classes, page, perPage, total)
}

// GetClassesByStudent handles get classes by student request
func (h *CourseHandler) GetClassesByStudent(c *gin.Context) {
	studentID := c.Param("studentId")

	result, err := h.service.GetClassesByStudent(c.Request.Context(), studentID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetClassesByLecturer handles get classes by lecturer request
func (h *CourseHandler) GetClassesByLecturer(c *gin.Context) {
	lecturerID := c.Param("lecturerId")

	result, err := h.service.GetClassesByLecturer(c.Request.Context(), lecturerID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// ========== Enrollment Handlers ==========

// CreateEnrollment handles create enrollment request
func (h *CourseHandler) CreateEnrollment(c *gin.Context) {
	var req service.CreateEnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.CreateEnrollment(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 201, result)
}

// GetEnrollment handles get enrollment by ID request
func (h *CourseHandler) GetEnrollment(c *gin.Context) {
	enrollmentID := c.Param("id")

	result, err := h.service.GetEnrollmentByID(c.Request.Context(), enrollmentID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetEnrollments handles get enrollments request
func (h *CourseHandler) GetEnrollments(c *gin.Context) {
	var req service.GetEnrollmentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	enrollments, total, err := h.service.GetEnrollments(c.Request.Context(), req)
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

	utils.PaginatedResponse(c, enrollments, page, perPage, total)
}

// GetEnrollmentsByStudent handles get enrollments by student request
func (h *CourseHandler) GetEnrollmentsByStudent(c *gin.Context) {
	studentID := c.Param("studentId")

	result, err := h.service.GetEnrollmentsByStudent(c.Request.Context(), studentID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetEnrollmentsByClass handles get enrollments by class request
func (h *CourseHandler) GetEnrollmentsByClass(c *gin.Context) {
	classID := c.Param("classId")

	result, err := h.service.GetEnrollmentsByClass(c.Request.Context(), classID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// UpdateEnrollmentStatus handles update enrollment status request (approve/reject)
func (h *CourseHandler) UpdateEnrollmentStatus(c *gin.Context) {
	enrollmentID := c.Param("id")

	var req service.UpdateEnrollmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.UpdateEnrollmentStatus(c.Request.Context(), enrollmentID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// UpdateEnrollmentGrade handles update enrollment grade request
func (h *CourseHandler) UpdateEnrollmentGrade(c *gin.Context) {
	enrollmentID := c.Param("id")

	var req service.UpdateEnrollmentGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.UpdateEnrollmentGrade(c.Request.Context(), enrollmentID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// DeleteEnrollment handles delete enrollment request
func (h *CourseHandler) DeleteEnrollment(c *gin.Context) {
	enrollmentID := c.Param("id")

	err := h.service.DeleteEnrollment(c.Request.Context(), enrollmentID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, gin.H{"message": "Enrollment deleted successfully"})
}

