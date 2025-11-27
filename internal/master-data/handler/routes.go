package handler

import (
	"github.com/gin-gonic/gin"
	"unsri-backend/internal/master-data/middleware"
	"unsri-backend/pkg/jwt"
)

// SetupRoutes sets up all routes for master data service
func SetupRoutes(router *gin.Engine, handler *MasterDataHandler, jwtToken *jwt.JWT) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "master-data-service"})
	})

	// Study Programs routes
	studyPrograms := router.Group("/api/v1/study-programs")
	studyPrograms.Use(middleware.AuthMiddleware(jwtToken))
	{
		studyPrograms.GET("", handler.GetStudyPrograms)
		studyPrograms.GET("/:id", handler.GetStudyProgram)
		studyPrograms.POST("", middleware.RoleMiddleware("dosen", "staff"), handler.CreateStudyProgram)
		studyPrograms.PUT("/:id", middleware.RoleMiddleware("dosen", "staff"), handler.UpdateStudyProgram)
		studyPrograms.DELETE("/:id", middleware.RoleMiddleware("dosen", "staff"), handler.DeleteStudyProgram)
	}

	// Academic Periods routes
	academicPeriods := router.Group("/api/v1/academic-periods")
	academicPeriods.Use(middleware.AuthMiddleware(jwtToken))
	{
		academicPeriods.GET("", handler.GetAcademicPeriods)
		academicPeriods.GET("/active", handler.GetActiveAcademicPeriod)
		academicPeriods.GET("/:id", handler.GetAcademicPeriod)
		academicPeriods.POST("", middleware.RoleMiddleware("dosen", "staff"), handler.CreateAcademicPeriod)
		academicPeriods.PUT("/:id", middleware.RoleMiddleware("dosen", "staff"), handler.UpdateAcademicPeriod)
		academicPeriods.DELETE("/:id", middleware.RoleMiddleware("dosen", "staff"), handler.DeleteAcademicPeriod)
	}

	// Rooms routes
	rooms := router.Group("/api/v1/rooms")
	rooms.Use(middleware.AuthMiddleware(jwtToken))
	{
		rooms.GET("", handler.GetRooms)
		rooms.GET("/:id", handler.GetRoom)
		rooms.POST("", middleware.RoleMiddleware("dosen", "staff"), handler.CreateRoom)
		rooms.PUT("/:id", middleware.RoleMiddleware("dosen", "staff"), handler.UpdateRoom)
		rooms.DELETE("/:id", middleware.RoleMiddleware("dosen", "staff"), handler.DeleteRoom)
	}
}

