package service

import (
	"context"
	"errors"
	"strings"

	"unsri-backend/internal/auth/repository"
	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
	"unsri-backend/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService handles authentication business logic
type AuthService struct {
	repo *repository.AuthRepository
	jwt  *jwt.JWT
}

// NewAuthService creates a new auth service
func NewAuthService(repo *repository.AuthRepository, jwtToken *jwt.JWT) *AuthService {
	return &AuthService{
		repo: repo,
		jwt:  jwtToken,
	}
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	User         *UserInfo `json:"user"`
}

// UserInfo represents user information in response
type UserInfo struct {
	ID        string            `json:"id"`
	Email     string            `json:"email"`
	Role      models.UserRole   `json:"role"`
	IsActive  bool              `json:"is_active"`
	Mahasiswa *models.Mahasiswa `json:"mahasiswa,omitempty"`
	Dosen     *models.Dosen     `json:"dosen,omitempty"`
	Staff     *models.Staff     `json:"staff,omitempty"`
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.NewUnauthorizedError("invalid email or password")
	}

	if !user.IsActive {
		return nil, apperrors.NewForbiddenError("account is inactive")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperrors.NewUnauthorizedError("invalid email or password")
	}

	accessToken, err := s.jwt.GenerateAccessToken(user.ID, string(user.Role), user.Email)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to generate access token", err)
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to generate refresh token", err)
	}

	userInfo := &UserInfo{
		ID:       user.ID,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
	}

	// Load role-specific data
	if user.Role == models.RoleMahasiswa && user.Mahasiswa != nil {
		userInfo.Mahasiswa = user.Mahasiswa
	} else if user.Role == models.RoleDosen && user.Dosen != nil {
		userInfo.Dosen = user.Dosen
	} else if user.Role == models.RoleStaff && user.Staff != nil {
		userInfo.Staff = user.Staff
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         userInfo,
	}, nil
}

// RegisterRequest represents registration request
type RegisterRequest struct {
	Email    string          `json:"email" binding:"required,email"`
	Password string          `json:"password" binding:"required,min=8"`
	Role     models.UserRole `json:"role" binding:"required,oneof=mahasiswa dosen staff"`
	NIM      string          `json:"nim,omitempty"` // For mahasiswa
	NIP      string          `json:"nip,omitempty"` // For dosen/staff
	Nama     string          `json:"nama" binding:"required"`
	Prodi    string          `json:"prodi,omitempty"`
	Angkatan int             `json:"angkatan,omitempty"` // For mahasiswa
	Jabatan  string          `json:"jabatan,omitempty"`  // For staff
	Unit     string          `json:"unit,omitempty"`     // For staff
}

// isConstraintViolation checks if error is a database constraint violation
func isConstraintViolation(err error) bool {
	if err == nil {
		return false
	}

	// Check for GORM duplicate entry error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	// Check for PostgreSQL unique constraint violation in error message
	errStr := strings.ToLower(err.Error())
	if strings.Contains(errStr, "unique constraint") ||
		strings.Contains(errStr, "duplicate key") ||
		strings.Contains(errStr, "violates unique constraint") ||
		strings.Contains(errStr, "23505") { // PostgreSQL unique violation error code
		return true
	}

	return false
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*UserInfo, error) {
	// Check if email already exists
	existingUser, _ := s.repo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, apperrors.NewConflictError("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to hash password", err)
	}

	// Use transaction to ensure atomicity - if any step fails, all changes are rolled back
	db := s.repo.GetDB()
	var createdUserID string

	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create user within transaction
		user := &models.User{
			Email:        req.Email,
			PasswordHash: string(hashedPassword),
			Role:         req.Role,
			IsActive:     true,
		}

		if err := tx.Create(user).Error; err != nil {
			if isConstraintViolation(err) {
				return apperrors.NewConflictError("email already registered")
			}
			return apperrors.NewInternalError("failed to create user", err)
		}

		// Store user ID for later use
		createdUserID = user.ID

		// Create role-specific record within transaction
		if req.Role == models.RoleMahasiswa {
			if req.NIM == "" {
				return apperrors.NewValidationError("NIM is required for mahasiswa")
			}
			// Check if NIM already exists
			var existingMahasiswa models.Mahasiswa
			if err := tx.Where("nim = ?", req.NIM).First(&existingMahasiswa).Error; err == nil {
				return apperrors.NewConflictError("NIM already registered")
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			mahasiswa := &models.Mahasiswa{
				UserID:   user.ID,
				NIM:      req.NIM,
				Nama:     req.Nama,
				Prodi:    req.Prodi,
				Angkatan: req.Angkatan,
			}
			if err := tx.Create(mahasiswa).Error; err != nil {
				if isConstraintViolation(err) {
					return apperrors.NewConflictError("NIM already registered")
				}
				return apperrors.NewInternalError("failed to create mahasiswa", err)
			}
		} else if req.Role == models.RoleDosen {
			if req.NIP == "" {
				return apperrors.NewValidationError("NIP is required for dosen")
			}
			// Check if NIP already exists
			var existingDosen models.Dosen
			if err := tx.Where("nip = ?", req.NIP).First(&existingDosen).Error; err == nil {
				return apperrors.NewConflictError("NIP already registered")
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			dosen := &models.Dosen{
				UserID: user.ID,
				NIP:    req.NIP,
				Nama:   req.Nama,
				Prodi:  req.Prodi,
			}
			if err := tx.Create(dosen).Error; err != nil {
				if isConstraintViolation(err) {
					return apperrors.NewConflictError("NIP already registered")
				}
				return apperrors.NewInternalError("failed to create dosen", err)
			}
		} else if req.Role == models.RoleStaff {
			if req.NIP == "" {
				return apperrors.NewValidationError("NIP is required for staff")
			}
			// Check if NIP already exists
			var existingStaff models.Staff
			if err := tx.Where("nip = ?", req.NIP).First(&existingStaff).Error; err == nil {
				return apperrors.NewConflictError("NIP already registered")
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			staff := &models.Staff{
				UserID:  user.ID,
				NIP:     req.NIP,
				Nama:    req.Nama,
				Jabatan: req.Jabatan,
				Unit:    req.Unit,
			}
			if err := tx.Create(staff).Error; err != nil {
				if isConstraintViolation(err) {
					return apperrors.NewConflictError("NIP already registered")
				}
				return apperrors.NewInternalError("failed to create staff", err)
			}
		}

		// Transaction will commit automatically if no error is returned
		return nil
	})

	if err != nil {
		// Transaction was rolled back automatically on error
		// Check if it's already an AppError
		if appErr, ok := err.(*apperrors.AppError); ok {
			return nil, appErr
		}
		return nil, err
	}

	// Query user again after transaction commits to get the full data with relations
	user, err := s.repo.FindByID(ctx, createdUserID)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to retrieve created user", err)
	}

	userInfo := &UserInfo{
		ID:       user.ID,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
	}

	if user.Mahasiswa != nil {
		userInfo.Mahasiswa = user.Mahasiswa
	} else if user.Dosen != nil {
		userInfo.Dosen = user.Dosen
	} else if user.Staff != nil {
		userInfo.Staff = user.Staff
	}

	return userInfo, nil
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse represents refresh token response
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshToken refreshes an access token
func (s *AuthService) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*RefreshTokenResponse, error) {
	claims, err := s.jwt.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, apperrors.NewUnauthorizedError("invalid refresh token")
	}

	user, err := s.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, apperrors.NewUnauthorizedError("user not found")
	}

	if !user.IsActive {
		return nil, apperrors.NewForbiddenError("account is inactive")
	}

	accessToken, err := s.jwt.GenerateAccessToken(user.ID, string(user.Role), user.Email)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to generate access token", err)
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to generate refresh token", err)
	}

	return &RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// VerifyToken verifies a JWT token
func (s *AuthService) VerifyToken(ctx context.Context, tokenString string) (*UserInfo, error) {
	claims, err := s.jwt.ValidateToken(tokenString)
	if err != nil {
		return nil, apperrors.NewUnauthorizedError("invalid token")
	}

	user, err := s.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, apperrors.NewUnauthorizedError("user not found")
	}

	if !user.IsActive {
		return nil, apperrors.NewForbiddenError("account is inactive")
	}

	userInfo := &UserInfo{
		ID:       user.ID,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
	}

	if user.Mahasiswa != nil {
		userInfo.Mahasiswa = user.Mahasiswa
	} else if user.Dosen != nil {
		userInfo.Dosen = user.Dosen
	} else if user.Staff != nil {
		userInfo.Staff = user.Staff
	}

	return userInfo, nil
}
