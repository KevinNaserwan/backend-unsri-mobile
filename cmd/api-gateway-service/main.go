package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"unsri-backend/internal/gateway/config"
	"unsri-backend/internal/gateway/handler"
	"unsri-backend/internal/gateway/middleware"
	"unsri-backend/internal/gateway/service"
	"unsri-backend/internal/shared/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	log := logger.New(cfg.LogLevel)
	log.Info("Starting API Gateway service...")

	// Initialize proxy service
	proxyService := service.NewProxyService(cfg, log)

	// Initialize handler
	gatewayHandler := handler.NewGatewayHandler(proxyService, log)

	// Setup router
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(middleware.CORS(cfg.CORS))
	router.Use(middleware.RequestLogger(log))
	router.Use(middleware.RateLimit(cfg.RateLimit))

	// Setup routes
	handler.SetupRoutes(router, gatewayHandler)

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", err)
		}
	}()

	log.Infof("API Gateway service started on port %s", cfg.Port)
	log.Info("Available routes:")
	log.Info("  GET  /health                 - Gateway health check")
	log.Info("  GET  /                       - Gateway info")
	log.Info("  ANY  /api/v1/*               - Proxy to microservices")
	log.Info("")
	log.Info("Service routing:")
	log.Info("  /api/v1/auth/*               -> auth-service:8081")
	log.Info("  /api/v1/users/*              -> user-service:8082")
	log.Info("  /api/v1/schedule/*           -> schedule-service:8083")
	log.Info("  /api/v1/attendance/*         -> attendance-service:8084")
	log.Info("  /api/v1/qr/*                 -> qr-service:8085")
	log.Info("  /api/v1/broadcast/*          -> broadcast-service:8086")
	log.Info("  /api/v1/notifications/*      -> notification-service:8087")
	log.Info("  /api/v1/calendar/*           -> calendar-service:8088")
	log.Info("  /api/v1/courses/*            -> course-service:8089")
	log.Info("  /api/v1/location/*           -> location-service:8090")
	log.Info("  /api/v1/locations/*          -> location-service:8090")
	log.Info("  /api/v1/access/*             -> access-service:8091")
	log.Info("  /api/v1/quick-actions/*      -> quick-actions-service:8092")
	log.Info("  /api/v1/files/*              -> file-storage-service:8093")
	log.Info("  /api/v1/search/*             -> search-service:8094")
	log.Info("  /api/v1/reports/*            -> report-service:8095")
	log.Info("  /api/v1/master-data/*        -> master-data-service:8096")
	log.Info("  /api/v1/leave/*              -> leave-service:8097")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down API Gateway...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", err)
	}

	log.Info("API Gateway exited")
}
