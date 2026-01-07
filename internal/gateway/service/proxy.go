package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"unsri-backend/internal/gateway/config"
	"unsri-backend/internal/shared/logger"
)

type ProxyService struct {
	config *config.Config
	client *http.Client
	logger logger.Logger
}

func NewProxyService(cfg *config.Config, log logger.Logger) *ProxyService {
	return &ProxyService{
		config: cfg,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
		logger: log,
	}
}

func (p *ProxyService) ProxyRequest(ctx context.Context, targetURL string, originalReq *http.Request) (*http.Response, error) {
	// Create new request
	body, err := io.ReadAll(originalReq.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	// Parse target URL
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse target URL: %w", err)
	}

	// Create new request with target URL
	proxyReq, err := http.NewRequestWithContext(ctx, originalReq.Method, target.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create proxy request: %w", err)
	}

	// Copy headers
	for name, values := range originalReq.Header {
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}

	// Add/modify headers for internal communication
	proxyReq.Header.Set("X-Forwarded-For", originalReq.RemoteAddr)
	proxyReq.Header.Set("X-Forwarded-Proto", "http")
	if originalReq.Host != "" {
		proxyReq.Header.Set("X-Forwarded-Host", originalReq.Host)
	}

	// Log the proxy request
	p.logger.Infof("[PROXY] %s %s -> %s", originalReq.Method, originalReq.URL.Path, target.String())

	// Make the request
	resp, err := p.client.Do(proxyReq)
	if err != nil {
		p.logger.Errorf("[PROXY] Failed to proxy request: %v", err)
		return nil, fmt.Errorf("failed to proxy request: %w", err)
	}

	return resp, nil
}

func (p *ProxyService) GetServiceURL(path string) (string, error) {
	// Handle health endpoints specially - they are at root level, not under /api/v1/
	switch path {
	case "/api/v1/auth/health":
		return p.config.Services.AuthService.URL() + "/health", nil
	case "/api/v1/users/health":
		return p.config.Services.UserService.URL() + "/health", nil
	case "/api/v1/work-attendance/health":
		return p.config.Services.AttendanceService.URL() + "/health", nil
	case "/api/v1/attendance/health":
		return p.config.Services.AttendanceService.URL() + "/health", nil
	case "/api/v1/schedule/health":
		return p.config.Services.ScheduleService.URL() + "/health", nil
	case "/api/v1/qr/health":
		return p.config.Services.QRService.URL() + "/health", nil
	case "/api/v1/courses/health":
		return p.config.Services.CourseService.URL() + "/health", nil
	case "/api/v1/broadcast/health":
		return p.config.Services.BroadcastService.URL() + "/health", nil
	case "/api/v1/notifications/health":
		return p.config.Services.NotificationService.URL() + "/health", nil
	case "/api/v1/calendar/health":
		return p.config.Services.CalendarService.URL() + "/health", nil
	case "/api/v1/locations/health":
		return p.config.Services.LocationService.URL() + "/health", nil
	case "/api/v1/access/health":
		return p.config.Services.AccessService.URL() + "/health", nil
	case "/api/v1/quick-actions/health":
		return p.config.Services.QuickActionsService.URL() + "/health", nil
	case "/api/v1/files/health":
		return p.config.Services.FileStorageService.URL() + "/health", nil
	case "/api/v1/search/health":
		return p.config.Services.SearchService.URL() + "/health", nil
	case "/api/v1/reports/health":
		return p.config.Services.ReportService.URL() + "/health", nil
	case "/api/v1/master-data/health":
		return p.config.Services.MasterDataService.URL() + "/health", nil
	case "/api/v1/leave/health":
		return p.config.Services.LeaveService.URL() + "/health", nil
	}

	// Handle regular API endpoints - keep the full path since backend services expect the full API path
	switch {
	case strings.HasPrefix(path, "/api/v1/auth"):
		return p.config.Services.AuthService.URL() + path, nil
	case strings.HasPrefix(path, "/api/v1/users"):
		return p.config.Services.UserService.URL() + path, nil
	case strings.HasPrefix(path, "/api/v1/work-attendance"):
		return p.config.Services.AttendanceService.URL() + path, nil
	case strings.HasPrefix(path, "/api/v1/attendance"):
		return p.config.Services.AttendanceService.URL() + path, nil
	case strings.HasPrefix(path, "/api/v1/schedule"):
		return p.config.Services.ScheduleService.URL() + path, nil
	case strings.HasPrefix(path, "/api/v1/qr"):
		return p.config.Services.QRService.URL() + path, nil
	case strings.HasPrefix(path, "/api/v1/courses"):
		return p.config.Services.CourseService.URL() + path, nil
	case strings.HasPrefix(path, "/api/v1/broadcast"):
		return p.config.Services.BroadcastService.URL() + path, nil
	case strings.HasPrefix(path, "/api/v1/notifications"):
		return p.config.Services.NotificationService.URL() + path, nil
	case strings.HasPrefix(path, "/api/v1/calendar"):
		return p.config.Services.CalendarService.URL() + path, nil
	case strings.HasPrefix(path, "/api/v1/locations"):
		return p.config.Services.LocationService.URL() + path, nil
	case strings.HasPrefix(path, "/api/v1/access"):
		return p.config.Services.AccessService.URL() + path, nil
	case strings.HasPrefix(path, "/api/v1/quick-actions"):
		return p.config.Services.QuickActionsService.URL() + path, nil
	case strings.HasPrefix(path, "/api/v1/files"):
		return p.config.Services.FileStorageService.URL() + path, nil
	case strings.HasPrefix(path, "/api/v1/search"):
		return p.config.Services.SearchService.URL() + path, nil
	case strings.HasPrefix(path, "/api/v1/reports"):
		return p.config.Services.ReportService.URL() + path, nil
	case strings.HasPrefix(path, "/api/v1/master-data"):
		return p.config.Services.MasterDataService.URL() + path, nil
	case strings.HasPrefix(path, "/api/v1/leave"):
		return p.config.Services.LeaveService.URL() + path, nil
	default:
		return "", fmt.Errorf("no service found for path: %s", path)
	}
}
