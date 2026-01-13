package handler

import (
	"unsri-backend/internal/attendance/middleware"
	"unsri-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all routes for attendance service
func SetupRoutes(router *gin.Engine, handler *AttendanceHandler, jwtToken *jwt.JWT) {
	// Public routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "attendance-service"})
	})

	// Protected routes
	v1 := router.Group("/api/v1/attendance")
	v1.Use(middleware.AuthMiddleware(jwtToken))
	{
		// QR code operations
		v1.POST("/qr/generate", middleware.RoleMiddleware("dosen", "staff"), handler.GenerateQR)
		v1.POST("/qr/scan", handler.ScanQR)

		// Attendance operations
		v1.GET("", handler.GetAttendances)
		v1.GET("/statistics", handler.GetStatistics)
		v1.GET("/overview", handler.GetOverview)
		v1.GET("/history", handler.GetAttendances) // Alias for GetAttendances
		v1.GET("/by-course/:courseId", handler.GetByCourse)
		v1.GET("/by-student/:studentId", handler.GetByStudent)
		v1.POST("/manual", middleware.RoleMiddleware("dosen", "staff"), handler.CreateManualAttendance)
		v1.PUT("/:id", middleware.RoleMiddleware("dosen", "staff"), handler.UpdateAttendance)

		// Campus attendance (tap in/out) - removed per request
	}

	// Schedule routes
	schedules := router.Group("/api/v1/schedules")
	schedules.Use(middleware.AuthMiddleware(jwtToken))
	{
		schedules.GET("", handler.GetSchedules)
		schedules.GET("/today", handler.GetTodaySchedules)
		schedules.GET("/:id", handler.GetSchedule)
		schedules.POST("", middleware.RoleMiddleware("dosen", "staff"), handler.CreateSchedule)
		schedules.PUT("/:id", middleware.RoleMiddleware("dosen", "staff"), handler.UpdateSchedule)
		schedules.DELETE("/:id", middleware.RoleMiddleware("dosen", "staff"), handler.DeleteSchedule)
	}

	// Work Attendance (Kepegawaian) routes
	workAttendance := router.Group("/api/v1/work-attendance")
	workAttendance.Use(middleware.AuthMiddleware(jwtToken))
	{
		// Check-in/out
<<<<<<< HEAD
		workAttendance.POST("/check-in", middleware.RoleMiddleware("admin", "dosen", "staff"), handler.CheckIn)
		workAttendance.POST("/check-out", middleware.RoleMiddleware("admin", "dosen", "staff"), handler.CheckOut)
		workAttendance.GET("/records", middleware.RoleMiddleware("admin", "dosen", "staff"), handler.GetWorkAttendanceRecords)
		workAttendance.GET("/statistics", middleware.RoleMiddleware("admin", "dosen", "staff"), handler.GetWorkAttendanceStatistics)
		workAttendance.GET("/summaries", middleware.RoleMiddleware("admin", "dosen", "staff"), handler.GetWorkAttendanceSummaries)
=======
		workAttendance.POST("/check-in", handler.CheckIn)
		workAttendance.POST("/check-out", handler.CheckOut)
		workAttendance.GET("/records", handler.GetWorkAttendanceRecords)
		workAttendance.GET("/all-records", middleware.RoleMiddleware("dosen", "staff"), handler.GetAllWorkAttendanceRecords)
>>>>>>> origin

		// Shift patterns (admin only)
		shifts := workAttendance.Group("/shifts")
		shifts.Use(middleware.RoleMiddleware("dosen", "staff"))
		{
			shifts.GET("", handler.GetShiftPatterns)
			shifts.GET("/:id", handler.GetShiftPattern)
			shifts.POST("", handler.CreateShiftPattern)
			shifts.PUT("/:id", handler.UpdateShiftPattern)
			shifts.DELETE("/:id", handler.DeleteShiftPattern)
		}

		// User shifts (admin only)
		userShifts := workAttendance.Group("/user-shifts")
		userShifts.Use(middleware.RoleMiddleware("dosen", "staff"))
		{
			userShifts.GET("/:userId", handler.GetUserShifts)
			userShifts.POST("", handler.CreateUserShift)
		}

		// Work schedules (admin only)
		workSchedules := workAttendance.Group("/schedules")
		workSchedules.Use(middleware.RoleMiddleware("dosen", "staff"))
		{
			workSchedules.GET("", handler.GetWorkSchedules)
			workSchedules.POST("", handler.CreateWorkSchedule)
		}
	}
}
