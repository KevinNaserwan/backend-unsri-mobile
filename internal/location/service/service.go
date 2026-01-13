package service

import (
	"context"

	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
	"unsri-backend/internal/location/repository"
)

// LocationService handles location business logic
type LocationService struct {
	repo *repository.LocationRepository
}

// NewLocationService creates a new location service
func NewLocationService(repo *repository.LocationRepository) *LocationService {
	return &LocationService{repo: repo}
}

// GetCheckInStatus gets current check-in status
func (s *LocationService) GetCheckInStatus(ctx context.Context, userID string) (*models.LocationHistory, error) {
	return s.repo.GetCurrentTapInStatus(ctx, userID)
}

// GetLocationHistoryRequest represents get location history request
type GetLocationHistoryRequest struct {
	Page    int `form:"page,default=1"`
	PerPage int `form:"per_page,default=20"`
}

// GetLocationHistory gets location history
func (s *LocationService) GetLocationHistory(ctx context.Context, userID string, req GetLocationHistoryRequest) ([]models.LocationHistory, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	return s.repo.GetLocationHistory(ctx, userID, perPage, (page-1)*perPage)
}

// GetGeofences gets all geofences
func (s *LocationService) GetGeofences(ctx context.Context) ([]models.Geofence, error) {
	return s.repo.GetAllGeofences(ctx)
}

// ValidateLocationRequest represents validate location request
type ValidateLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// ValidateLocation validates if location is within geofence
func (s *LocationService) ValidateLocation(ctx context.Context, req ValidateLocationRequest) (*models.Geofence, error) {
	geofence, err := s.repo.CheckLocationInGeofence(ctx, req.Latitude, req.Longitude)
	if err != nil {
		return nil, apperrors.NewBadRequestError("location not within allowed area")
	}
	return geofence, nil
}

// CreateGeofenceRequest represents create geofence request
type CreateGeofenceRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description,omitempty"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Radius      float64 `json:"radius" binding:"required"` // in meters
}

// CreateGeofence creates a new geofence
func (s *LocationService) CreateGeofence(ctx context.Context, req CreateGeofenceRequest) (*models.Geofence, error) {
	geofence := &models.Geofence{
		Name:        req.Name,
		Description: req.Description,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Radius:      req.Radius,
		IsActive:    true,
	}

	if err := s.repo.CreateGeofence(ctx, geofence); err != nil {
		return nil, apperrors.NewInternalError("failed to create geofence", err)
	}

	return geofence, nil
}

