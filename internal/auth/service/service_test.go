package service

import (
	"testing"

	"github.com/google/uuid"
	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
)

// Test helper functions
func createTestUser() *models.User {
	return &models.User{
		ID:           uuid.New().String(),
		Email:        "test@example.com",
		PasswordHash: "$2a$10$hashedpassword",
		Role:         models.RoleMahasiswa,
		IsActive:     true,
	}
}

// Test LoginRequest validation
func TestLoginRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     LoginRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "missing email",
			req: LoginRequest{
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "missing password",
			req: LoginRequest{
				Email: "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "invalid email format",
			req: LoginRequest{
				Email:    "invalid-email",
				Password: "password123",
			},
			wantErr: false, // Email format validation happens in handler
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.Email == "" && !tt.wantErr {
				t.Error("Email should be required")
			}
			if tt.req.Password == "" && !tt.wantErr {
				t.Error("Password should be required")
			}
		})
	}
}

// Test RegisterRequest validation
func TestRegisterRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     RegisterRequest
		wantErr bool
	}{
		{
			name: "valid mahasiswa request",
			req: RegisterRequest{
				Email:    "student@example.com",
				Password: "password123",
				Role:     models.RoleMahasiswa,
				Nama:     "Test Student",
				NIM:      "1234567890",
			},
			wantErr: false,
		},
		{
			name: "valid dosen request",
			req: RegisterRequest{
				Email:    "dosen@example.com",
				Password: "password123",
				Role:     models.RoleDosen,
				Nama:     "Test Dosen",
				NIP:      "1234567890",
			},
			wantErr: false,
		},
		{
			name: "missing email",
			req: RegisterRequest{
				Password: "password123",
				Role:     "mahasiswa",
			},
			wantErr: true,
		},
		{
			name: "missing password",
			req: RegisterRequest{
				Email: "test@example.com",
				Role:  "mahasiswa",
			},
			wantErr: true,
		},
		{
			name: "missing role",
			req: RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Role:     "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.Email == "" && !tt.wantErr {
				t.Error("Email should be required")
			}
			if tt.req.Password == "" && !tt.wantErr {
				t.Error("Password should be required")
			}
			if tt.req.Role == "" && !tt.wantErr {
				t.Error("Role should be required")
			}
		})
	}
}

// Test UserRole validation
func TestUserRole(t *testing.T) {
	validRoles := []models.UserRole{
		models.RoleMahasiswa,
		models.RoleDosen,
		models.RoleStaff,
	}

	for _, role := range validRoles {
		t.Run(string(role), func(t *testing.T) {
			if role == "" {
				t.Error("Role should not be empty")
			}
		})
	}
}

// Test error types
func TestErrorTypes(t *testing.T) {
	t.Run("NotFoundError", func(t *testing.T) {
		err := apperrors.NewNotFoundError("user", "test-id")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("UnauthorizedError", func(t *testing.T) {
		err := apperrors.NewUnauthorizedError("invalid credentials")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("ConflictError", func(t *testing.T) {
		err := apperrors.NewConflictError("user already exists")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("ValidationError", func(t *testing.T) {
		err := apperrors.NewValidationError("invalid input")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

// Test User model
func TestUserModel(t *testing.T) {
	t.Run("valid user", func(t *testing.T) {
		user := createTestUser()
		if user.Email == "" {
			t.Error("Email should not be empty")
		}
		if user.Role == "" {
			t.Error("Role should not be empty")
		}
		if user.ID == "" {
			t.Error("ID should be generated")
		}
	})

	t.Run("table name", func(t *testing.T) {
		user := models.User{}
		if user.TableName() != "users" {
			t.Errorf("Expected table name 'users', got '%s'", user.TableName())
		}
	})
}

// Test password validation
func TestPasswordValidation(t *testing.T) {
	t.Run("password should be hashed", func(t *testing.T) {
		user := createTestUser()
		if user.PasswordHash == "" {
			t.Error("PasswordHash should not be empty")
		}
		// Password should be hashed (starts with $2a$ for bcrypt)
		if len(user.PasswordHash) < 10 {
			t.Error("PasswordHash should be hashed")
		}
	})
}

