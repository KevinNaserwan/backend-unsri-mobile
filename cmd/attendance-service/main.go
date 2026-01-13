package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"unsri-backend/internal/attendance/config"
	"unsri-backend/internal/attendance/handler"
	"unsri-backend/internal/attendance/repository"
	"unsri-backend/internal/attendance/service"
	fileRepo "unsri-backend/internal/file-storage/repository"
	fileSvc "unsri-backend/internal/file-storage/service"
	locationRepo "unsri-backend/internal/location/repository"
	"unsri-backend/internal/shared/database"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/models"
	"unsri-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	log := logger.New(cfg.LogLevel)
	log.Info("Starting attendance service...")

	// Initialize database
	db, err := database.NewPostgres(database.Config{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		DBName:          cfg.Database.DBName,
		SSLMode:         cfg.Database.SSLMode,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	})
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	// Auto migrate models
	if err := db.AutoMigrate(
		&models.Attendance{},
		&models.AttendanceSession{},
		&models.Schedule{},
		// Work Attendance (HRIS) models
		&models.ShiftPattern{},
		&models.UserShift{},
		&models.WorkSchedule{},
		&models.WorkAttendanceSession{},
		&models.WorkAttendanceRecord{},
	); err != nil {
		log.Fatal("Failed to migrate database", err)
	}

	// Initialize JWT
	jwtToken := jwt.NewJWT(
		cfg.JWT.SecretKey,
		15*time.Minute, // Access token TTL
		7*24*time.Hour, // Refresh token TTL
	)

	// Initialize repository
	attendanceRepo := repository.NewAttendanceRepository(db)
	locationRepository := locationRepo.NewLocationRepository(db)

	// Initialize file storage (local or minio based on env)
	viper.SetDefault("STORAGE_TYPE", "minio")
	viper.SetDefault("STORAGE_BASE_PATH", "./storage")
	viper.SetDefault("STORAGE_BASE_URL", "http://localhost:8093/files")
	viper.SetDefault("STORAGE_MAX_SIZE", int64(10<<20)) // 10MB
	viper.SetDefault("MINIO_ENDPOINT", "localhost:9000")
	viper.SetDefault("MINIO_ACCESS_KEY", "minioadmin")
	viper.SetDefault("MINIO_SECRET_KEY", "minioadmin")
	viper.SetDefault("MINIO_BUCKET", "unsri")
	viper.SetDefault("MINIO_USE_SSL", false)
	viper.SetDefault("MINIO_REGION", "")
	viper.AutomaticEnv()
	// Ensure storage directory exists (for local)
	_ = os.MkdirAll(viper.GetString("STORAGE_BASE_PATH"), 0755)
	filesRepository := fileRepo.NewFileRepository(db)
	filesService := fileSvc.NewFileStorageService(filesRepository, fileSvc.StorageConfig{
		Type:           viper.GetString("STORAGE_TYPE"),
		BasePath:       viper.GetString("STORAGE_BASE_PATH"),
		BaseURL:        viper.GetString("STORAGE_BASE_URL"),
		MaxSize:        viper.GetInt64("STORAGE_MAX_SIZE"),
		MinioEndpoint:  viper.GetString("MINIO_ENDPOINT"),
		MinioAccessKey: viper.GetString("MINIO_ACCESS_KEY"),
		MinioSecretKey: viper.GetString("MINIO_SECRET_KEY"),
		MinioBucket:    viper.GetString("MINIO_BUCKET"),
		MinioUseSSL:    viper.GetBool("MINIO_USE_SSL"),
		MinioRegion:    viper.GetString("MINIO_REGION"),
	})

	// Initialize service
	attendanceService := service.NewAttendanceService(attendanceRepo, locationRepository, jwtToken)

	// Initialize handler
	attendanceHandler := handler.NewAttendanceHandler(attendanceService, log, filesService)

	// Setup router
	router := gin.Default()
	router.Use(gin.Recovery())
	handler.SetupRoutes(router, attendanceHandler, jwtToken)

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

	log.Infof("Attendance service started on port %s", cfg.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", err)
	}

	log.Info("Server exited")
}
