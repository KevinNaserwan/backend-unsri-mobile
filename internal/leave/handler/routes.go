package handler

import (
	"github.com/gin-gonic/gin"
	"unsri-backend/internal/leave/middleware"
	"unsri-backend/pkg/jwt"
)

// SetupRoutes sets up all routes for leave service
func SetupRoutes(router *gin.Engine, handler *LeaveHandler, jwtToken *jwt.JWT) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "leave-service"})
	})

	// Leave Requests routes
	leaveRequests := router.Group("/api/v1/leave-requests")
	leaveRequests.Use(middleware.AuthMiddleware(jwtToken))
	{
		leaveRequests.GET("", handler.GetLeaveRequests)
		leaveRequests.GET("/:id", handler.GetLeaveRequest)
		leaveRequests.POST("", handler.CreateLeaveRequest)
		leaveRequests.PUT("/:id/approve", middleware.RoleMiddleware("dosen", "staff"), handler.ApproveLeaveRequest)
		leaveRequests.PUT("/:id/reject", middleware.RoleMiddleware("dosen", "staff"), handler.RejectLeaveRequest)
		leaveRequests.PUT("/:id/cancel", handler.CancelLeaveRequest)
		leaveRequests.DELETE("/:id", middleware.RoleMiddleware("dosen", "staff"), handler.DeleteLeaveRequest)
		leaveRequests.GET("/by-user/:userId", handler.GetLeaveRequestsByUser)
	}

	// Leave Quotas routes
	leaveQuotas := router.Group("/api/v1/leave-quotas")
	leaveQuotas.Use(middleware.AuthMiddleware(jwtToken))
	{
		leaveQuotas.GET("", handler.GetLeaveQuotas)
		leaveQuotas.GET("/:id", handler.GetLeaveQuota)
		leaveQuotas.POST("", middleware.RoleMiddleware("dosen", "staff"), handler.CreateLeaveQuota)
		leaveQuotas.PUT("/:id", middleware.RoleMiddleware("dosen", "staff"), handler.UpdateLeaveQuota)
		leaveQuotas.DELETE("/:id", middleware.RoleMiddleware("dosen", "staff"), handler.DeleteLeaveQuota)
		leaveQuotas.GET("/by-user/:userId", handler.GetLeaveQuotasByUser)
	}
}

