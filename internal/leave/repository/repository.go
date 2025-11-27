package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"unsri-backend/internal/shared/models"
)

// LeaveRepository handles leave data operations
type LeaveRepository struct {
	db *gorm.DB
}

// NewLeaveRepository creates a new leave repository
func NewLeaveRepository(db *gorm.DB) *LeaveRepository {
	return &LeaveRepository{db: db}
}

// ========== Leave Request Methods ==========

// CreateLeaveRequest creates a new leave request
func (r *LeaveRepository) CreateLeaveRequest(ctx context.Context, leaveRequest *models.LeaveRequest) error {
	return r.db.WithContext(ctx).Create(leaveRequest).Error
}

// GetLeaveRequestByID gets a leave request by ID
func (r *LeaveRepository) GetLeaveRequestByID(ctx context.Context, id string) (*models.LeaveRequest, error) {
	var leaveRequest models.LeaveRequest
	if err := r.db.WithContext(ctx).Preload("User").Preload("Approver").Where("id = ?", id).First(&leaveRequest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("leave request not found")
		}
		return nil, err
	}
	return &leaveRequest, nil
}

// GetAllLeaveRequests gets all leave requests with filters
func (r *LeaveRepository) GetAllLeaveRequests(ctx context.Context, userID *string, leaveType *string, status *string, startDate, endDate *time.Time, limit, offset int) ([]models.LeaveRequest, int64, error) {
	var leaveRequests []models.LeaveRequest
	var total int64

	query := r.db.WithContext(ctx).Model(&models.LeaveRequest{})

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if leaveType != nil {
		query = query.Where("leave_type = ?", *leaveType)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if startDate != nil {
		query = query.Where("start_date >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("end_date <= ?", endDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("User").Preload("Approver").Limit(limit).Offset(offset).Order("created_at DESC").Find(&leaveRequests).Error; err != nil {
		return nil, 0, err
	}

	return leaveRequests, total, nil
}

// GetLeaveRequestsByUserID gets leave requests by user ID
func (r *LeaveRepository) GetLeaveRequestsByUserID(ctx context.Context, userID string, status *string) ([]models.LeaveRequest, error) {
	var leaveRequests []models.LeaveRequest
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Preload("Approver").Order("created_at DESC").Find(&leaveRequests).Error; err != nil {
		return nil, err
	}
	return leaveRequests, nil
}

// UpdateLeaveRequest updates a leave request
func (r *LeaveRepository) UpdateLeaveRequest(ctx context.Context, leaveRequest *models.LeaveRequest) error {
	return r.db.WithContext(ctx).Save(leaveRequest).Error
}

// DeleteLeaveRequest soft deletes a leave request
func (r *LeaveRepository) DeleteLeaveRequest(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.LeaveRequest{}, "id = ?", id).Error
}

// ========== Leave Quota Methods ==========

// CreateLeaveQuota creates a new leave quota
func (r *LeaveRepository) CreateLeaveQuota(ctx context.Context, quota *models.LeaveQuota) error {
	return r.db.WithContext(ctx).Create(quota).Error
}

// GetLeaveQuotaByID gets a leave quota by ID
func (r *LeaveRepository) GetLeaveQuotaByID(ctx context.Context, id string) (*models.LeaveQuota, error) {
	var quota models.LeaveQuota
	if err := r.db.WithContext(ctx).Preload("User").Where("id = ?", id).First(&quota).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("leave quota not found")
		}
		return nil, err
	}
	return &quota, nil
}

// GetLeaveQuotaByUserAndTypeAndYear gets leave quota by user, type, and year
func (r *LeaveRepository) GetLeaveQuotaByUserAndTypeAndYear(ctx context.Context, userID string, leaveType models.LeaveType, year int) (*models.LeaveQuota, error) {
	var quota models.LeaveQuota
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND leave_type = ? AND year = ?", userID, leaveType, year).
		First(&quota).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("leave quota not found")
		}
		return nil, err
	}
	return &quota, nil
}

// GetAllLeaveQuotas gets all leave quotas with filters
func (r *LeaveRepository) GetAllLeaveQuotas(ctx context.Context, userID *string, leaveType *string, year *int, limit, offset int) ([]models.LeaveQuota, int64, error) {
	var quotas []models.LeaveQuota
	var total int64

	query := r.db.WithContext(ctx).Model(&models.LeaveQuota{})

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if leaveType != nil {
		query = query.Where("leave_type = ?", *leaveType)
	}
	if year != nil {
		query = query.Where("year = ?", *year)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("User").Limit(limit).Offset(offset).Order("year DESC, leave_type ASC").Find(&quotas).Error; err != nil {
		return nil, 0, err
	}

	return quotas, total, nil
}

// GetLeaveQuotasByUserID gets leave quotas by user ID
func (r *LeaveRepository) GetLeaveQuotasByUserID(ctx context.Context, userID string, year *int) ([]models.LeaveQuota, error) {
	var quotas []models.LeaveQuota
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if year != nil {
		query = query.Where("year = ?", *year)
	}

	if err := query.Order("year DESC, leave_type ASC").Find(&quotas).Error; err != nil {
		return nil, err
	}
	return quotas, nil
}

// UpdateLeaveQuota updates a leave quota
func (r *LeaveRepository) UpdateLeaveQuota(ctx context.Context, quota *models.LeaveQuota) error {
	return r.db.WithContext(ctx).Save(quota).Error
}

// DeleteLeaveQuota deletes a leave quota
func (r *LeaveRepository) DeleteLeaveQuota(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.LeaveQuota{}, "id = ?", id).Error
}

