package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port      string
	LogLevel  string
	Timeout   time.Duration
	Services  ServicesConfig
	CORS      CORSConfig
	RateLimit RateLimitConfig
}

type ServicesConfig struct {
	AuthService         ServiceEndpoint
	UserService         ServiceEndpoint
	AttendanceService   ServiceEndpoint
	ScheduleService     ServiceEndpoint
	QRService           ServiceEndpoint
	CourseService       ServiceEndpoint
	BroadcastService    ServiceEndpoint
	NotificationService ServiceEndpoint
	CalendarService     ServiceEndpoint
	LocationService     ServiceEndpoint
	AccessService       ServiceEndpoint
	QuickActionsService ServiceEndpoint
	FileStorageService  ServiceEndpoint
	SearchService       ServiceEndpoint
	ReportService       ServiceEndpoint
	MasterDataService   ServiceEndpoint
	LeaveService        ServiceEndpoint
}

type ServiceEndpoint struct {
	Host string
	Port string
}

func (s ServiceEndpoint) URL() string {
	return "http://" + s.Host + ":" + s.Port
}

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

type RateLimitConfig struct {
	Enabled bool
	RPS     int
	Burst   int
}

func Load() *Config {
	viper.AutomaticEnv()

	// Create config with direct viper calls (bypassing unmarshal issues)
	cfg := &Config{
		Port:     viper.GetString("PORT"),
		LogLevel: viper.GetString("LOG_LEVEL"),
	}

	// Set defaults if not provided
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	// Parse timeout
	timeoutStr := viper.GetString("TIMEOUT")
	if timeoutStr == "" {
		timeoutStr = "30s"
	}
	if timeout, err := time.ParseDuration(timeoutStr); err == nil {
		cfg.Timeout = timeout
	} else {
		cfg.Timeout = 30 * time.Second
	}

	// Directly read service configurations from environment variables
	cfg.Services = ServicesConfig{
		AuthService:         ServiceEndpoint{Host: getEnvOrDefault("AUTH_SERVICE_HOST", "auth-service"), Port: getEnvOrDefault("AUTH_SERVICE_PORT", "8081")},
		UserService:         ServiceEndpoint{Host: getEnvOrDefault("USER_SERVICE_HOST", "user-service"), Port: getEnvOrDefault("USER_SERVICE_PORT", "8082")},
		AttendanceService:   ServiceEndpoint{Host: getEnvOrDefault("ATTENDANCE_SERVICE_HOST", "attendance-service"), Port: getEnvOrDefault("ATTENDANCE_SERVICE_PORT", "8084")},
		ScheduleService:     ServiceEndpoint{Host: getEnvOrDefault("SCHEDULE_SERVICE_HOST", "schedule-service"), Port: getEnvOrDefault("SCHEDULE_SERVICE_PORT", "8083")},
		QRService:           ServiceEndpoint{Host: getEnvOrDefault("QR_SERVICE_HOST", "qr-service"), Port: getEnvOrDefault("QR_SERVICE_PORT", "8085")},
		CourseService:       ServiceEndpoint{Host: getEnvOrDefault("COURSE_SERVICE_HOST", "course-service"), Port: getEnvOrDefault("COURSE_SERVICE_PORT", "8089")},
		BroadcastService:    ServiceEndpoint{Host: getEnvOrDefault("BROADCAST_SERVICE_HOST", "broadcast-service"), Port: getEnvOrDefault("BROADCAST_SERVICE_PORT", "8086")},
		NotificationService: ServiceEndpoint{Host: getEnvOrDefault("NOTIFICATION_SERVICE_HOST", "notification-service"), Port: getEnvOrDefault("NOTIFICATION_SERVICE_PORT", "8087")},
		CalendarService:     ServiceEndpoint{Host: getEnvOrDefault("CALENDAR_SERVICE_HOST", "calendar-service"), Port: getEnvOrDefault("CALENDAR_SERVICE_PORT", "8088")},
		LocationService:     ServiceEndpoint{Host: getEnvOrDefault("LOCATION_SERVICE_HOST", "location-service"), Port: getEnvOrDefault("LOCATION_SERVICE_PORT", "8090")},
		AccessService:       ServiceEndpoint{Host: getEnvOrDefault("ACCESS_SERVICE_HOST", "access-service"), Port: getEnvOrDefault("ACCESS_SERVICE_PORT", "8091")},
		QuickActionsService: ServiceEndpoint{Host: getEnvOrDefault("QUICK_ACTIONS_SERVICE_HOST", "quick-actions-service"), Port: getEnvOrDefault("QUICK_ACTIONS_SERVICE_PORT", "8092")},
		FileStorageService:  ServiceEndpoint{Host: getEnvOrDefault("FILE_STORAGE_SERVICE_HOST", "file-storage-service"), Port: getEnvOrDefault("FILE_STORAGE_SERVICE_PORT", "8093")},
		SearchService:       ServiceEndpoint{Host: getEnvOrDefault("SEARCH_SERVICE_HOST", "search-service"), Port: getEnvOrDefault("SEARCH_SERVICE_PORT", "8094")},
		ReportService:       ServiceEndpoint{Host: getEnvOrDefault("REPORT_SERVICE_HOST", "report-service"), Port: getEnvOrDefault("REPORT_SERVICE_PORT", "8095")},
		MasterDataService:   ServiceEndpoint{Host: getEnvOrDefault("MASTER_DATA_SERVICE_HOST", "master-data-service"), Port: getEnvOrDefault("MASTER_DATA_SERVICE_PORT", "8096")},
		LeaveService:        ServiceEndpoint{Host: getEnvOrDefault("LEAVE_SERVICE_HOST", "leave-service"), Port: getEnvOrDefault("LEAVE_SERVICE_PORT", "8097")},
	}

	// Set CORS configuration
	cfg.CORS = CORSConfig{
		AllowOrigins:     getStringSliceOrDefault("ALLOW_ORIGINS", []string{"*"}),
		AllowMethods:     getStringSliceOrDefault("ALLOW_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}),
		AllowHeaders:     getStringSliceOrDefault("ALLOW_HEADERS", []string{"*"}),
		ExposeHeaders:    getStringSliceOrDefault("EXPOSE_HEADERS", []string{"*"}),
		AllowCredentials: getBoolOrDefault("ALLOW_CREDENTIALS", true),
		MaxAge:           getIntOrDefault("MAX_AGE", 86400),
	}

	// Set Rate Limit configuration
	cfg.RateLimit = RateLimitConfig{
		Enabled: getBoolOrDefault("RATE_LIMIT_ENABLED", false),
		RPS:     getIntOrDefault("RATE_LIMIT_RPS", 100),
		Burst:   getIntOrDefault("RATE_LIMIT_BURST", 200),
	}

	return cfg
}

// Helper functions for environment variable handling
func getEnvOrDefault(key, defaultValue string) string {
	if value := viper.GetString(key); value != "" {
		return value
	}
	return defaultValue
}

func getStringSliceOrDefault(key string, defaultValue []string) []string {
	if value := viper.GetStringSlice(key); len(value) > 0 {
		return value
	}
	return defaultValue
}

func getBoolOrDefault(key string, defaultValue bool) bool {
	if viper.IsSet(key) {
		return viper.GetBool(key)
	}
	return defaultValue
}

func getIntOrDefault(key string, defaultValue int) int {
	if viper.IsSet(key) {
		return viper.GetInt(key)
	}
	return defaultValue
}
