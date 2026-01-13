package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"unsri-backend/internal/file-storage/config"
	"unsri-backend/internal/file-storage/handler"
	"unsri-backend/internal/file-storage/repository"
	"unsri-backend/internal/file-storage/service"
	"unsri-backend/internal/shared/database"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/models"
	"unsri-backend/internal/shared/storage"
	"unsri-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel)
	log.Info("Starting file storage service...")

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

	if err := db.AutoMigrate(
		&models.File{},
	); err != nil {
		log.Fatal("Failed to migrate database", err)
	}

	// Create storage directory (for local storage fallback)
	if err := os.MkdirAll(cfg.Storage.BasePath, 0755); err != nil {
		log.Fatal("Failed to create storage directory", err)
	}

	// Initialize MinIO client if storage type is minio
	var minioClient *storage.MinIOClient
	if cfg.Storage.Type == "minio" {
		minioClient, err = storage.NewMinIOClient(storage.MinIOConfig{
			Endpoint:        cfg.Storage.MinIO.Endpoint,
			AccessKeyID:     cfg.Storage.MinIO.AccessKeyID,
			SecretAccessKey: cfg.Storage.MinIO.SecretAccessKey,
			UseSSL:          cfg.Storage.MinIO.UseSSL,
			Region:          cfg.Storage.MinIO.Region,
		})
		if err != nil {
			log.Fatal("Failed to initialize MinIO client", err)
		}

		// Create buckets if they don't exist
		ctx := context.Background()
		if err := minioClient.CreateBucketIfNotExists(ctx, cfg.Storage.MinIO.BucketProfiles); err != nil {
			log.Fatal("Failed to create profiles bucket", err)
		}
		if err := minioClient.CreateBucketIfNotExists(ctx, cfg.Storage.MinIO.BucketSelfies); err != nil {
			log.Fatal("Failed to create selfies bucket", err)
		}
		if err := minioClient.CreateBucketIfNotExists(ctx, cfg.Storage.MinIO.BucketDocuments); err != nil {
			log.Fatal("Failed to create documents bucket", err)
		}

		log.Info("MinIO client initialized and buckets created")
	}

	jwtToken := jwt.NewJWT(
		cfg.JWT.SecretKey,
		15*time.Minute,
		7*24*time.Hour,
	)

	fileRepo := repository.NewFileRepository(db)
	fileService := service.NewFileStorageService(fileRepo, service.StorageConfig{
<<<<<<< HEAD
		Type:     cfg.Storage.Type,
		BasePath: cfg.Storage.BasePath,
		BaseURL:  cfg.Storage.BaseURL,
		MaxSize:  cfg.Storage.MaxSize,
		MinIO: service.MinIOStorageConfig{
			Endpoint:        cfg.Storage.MinIO.Endpoint,
			AccessKeyID:     cfg.Storage.MinIO.AccessKeyID,
			SecretAccessKey: cfg.Storage.MinIO.SecretAccessKey,
			UseSSL:          cfg.Storage.MinIO.UseSSL,
			Region:          cfg.Storage.MinIO.Region,
			BucketProfiles:  cfg.Storage.MinIO.BucketProfiles,
			BucketSelfies:   cfg.Storage.MinIO.BucketSelfies,
			BucketDocuments: cfg.Storage.MinIO.BucketDocuments,
		},
	}, minioClient)
=======
		Type:           cfg.Storage.Type,
		BasePath:       cfg.Storage.BasePath,
		BaseURL:        cfg.Storage.BaseURL,
		MaxSize:        cfg.Storage.MaxSize,
		MinioEndpoint:  cfg.Storage.MinioEndpoint,
		MinioAccessKey: cfg.Storage.MinioAccessKey,
		MinioSecretKey: cfg.Storage.MinioSecretKey,
		MinioBucket:    cfg.Storage.MinioBucket,
		MinioUseSSL:    cfg.Storage.MinioUseSSL,
		MinioRegion:    cfg.Storage.MinioRegion,
	})
>>>>>>> origin
	fileHandler := handler.NewFileStorageHandler(fileService, log)

	router := gin.Default()
	router.Use(gin.Recovery())
	router.MaxMultipartMemory = 10 << 20 // 10 MB

	// Serve static files if storage type is local
	if cfg.Storage.Type == "local" {
		log.Infof("Serving static files from %s at /files", cfg.Storage.BasePath)
		router.Static("/files", cfg.Storage.BasePath)
	}

	handler.SetupRoutes(router, fileHandler, jwtToken)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", err)
		}
	}()

	log.Infof("File storage service started on port %s", cfg.Port)

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
