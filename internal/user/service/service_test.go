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

// Test error types
func TestErrorTypes(t *testing.T) {
	t.Run("NotFoundError", func(t *testing.T) {
		err := apperrors.NewNotFoundError("user", "test-id")
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

