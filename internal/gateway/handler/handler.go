package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"unsri-backend/internal/gateway/service"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/utils"
)

type GatewayHandler struct {
	proxyService *service.ProxyService
	logger       logger.Logger
}

func NewGatewayHandler(proxyService *service.ProxyService, logger logger.Logger) *GatewayHandler {
	return &GatewayHandler{
		proxyService: proxyService,
		logger:       logger,
	}
}

// Health check for the gateway itself
func (h *GatewayHandler) Health(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"service": "api-gateway",
		"status":  "healthy",
		"version": "1.0.0",
	})
}

// ProxyHandler handles all API requests and proxies them to appropriate services
func (h *GatewayHandler) ProxyHandler(c *gin.Context) {
	path := c.Request.URL.Path
	
	// Get target service URL
	targetURL, err := h.proxyService.GetServiceURL(path)
	if err != nil {
		h.logger.Errorf("[GATEWAY] Service not found for path %s: %v", path, err)
		utils.NotFoundResponse(c, "service", path)
		return
	}

	// Add query parameters to target URL
	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}

	// Proxy the request
	resp, err := h.proxyService.ProxyRequest(c.Request.Context(), targetURL, c.Request)
	if err != nil {
		h.logger.Errorf("[GATEWAY] Proxy request failed: %v", err)
		utils.ErrorResponse(c, http.StatusBadGateway, err)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for name, values := range resp.Header {
		for _, value := range values {
			c.Header(name, value)
		}
	}

	// Set status code
	c.Status(resp.StatusCode)

	// Copy response body
	if _, err := io.Copy(c.Writer, resp.Body); err != nil {
		h.logger.Errorf("[GATEWAY] Failed to copy response body: %v", err)
		return
	}
}

// SetupRoutes configures all routes for the API Gateway
func SetupRoutes(router *gin.Engine, handler *GatewayHandler) {
	// Gateway health check
	router.GET("/health", handler.Health)

	// Handle root path
	router.GET("/", func(c *gin.Context) {
		utils.SuccessResponse(c, http.StatusOK, gin.H{
			"service": "UNSRI Backend API Gateway",
			"version": "1.0.0",
			"status":  "running",
			"message": "Welcome to UNSRI Backend API Gateway",
		})
	})

	// Proxy all API requests
	router.Any("/api/v1/*path", handler.ProxyHandler)
}
