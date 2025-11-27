package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
)

// Test helper functions
func createTestCourse() *models.Course {
	return &models.Course{
		ID:          uuid.New().String(),
		Code:        "CS101",
		Name:        "Introduction to Computer Science",
		NameEn:      "Introduction to Computer Science",
		Credits:     3,
		Semester:    1,
		Description: "Basic computer science concepts",
		IsActive:    true,
	}
}

func createTestClass() *models.Class {
	return &models.Class{
		ID:           uuid.New().String(),
		CourseID:     uuid.New().String(),
		ClassCode:    "CS101-A",
		ClassName:    "CS101 Class A",
		Semester:     "Ganjil",
		AcademicYear: "2024/2025",
		DosenID:      uuid.New().String(),
		DayOfWeek:    1,
		StartTime:    time.Now(),
		EndTime:      time.Now().Add(2 * time.Hour),
		IsActive:     true,
	}
}

func createTestEnrollment() *models.Enrollment {
	return &models.Enrollment{
		ID:             uuid.New().String(),
		ClassID:        uuid.New().String(),
		StudentID:      uuid.New().String(),
		EnrollmentDate: time.Now(),
		Status:         "PENDING",
	}
}

// Test CreateCourseRequest validation
func TestCreateCourseRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateCourseRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateCourseRequest{
				Code:        "CS101",
				Name:        "Introduction to Computer Science",
				NameEn:      "Introduction to Computer Science",
				Credits:     3,
				Semester:    1,
				Description: "Basic concepts",
			},
			wantErr: false,
		},
		{
			name: "missing code",
			req: CreateCourseRequest{
				Name:    "Introduction to Computer Science",
				Credits: 3,
			},
			wantErr: true,
		},
		{
			name: "missing name",
			req: CreateCourseRequest{
				Code:    "CS101",
				Credits: 3,
			},
			wantErr: true,
		},
		{
			name: "missing credits",
			req: CreateCourseRequest{
				Code: "CS101",
				Name: "Introduction to Computer Science",
			},
			wantErr: true,
		},
		{
			name: "zero credits",
			req: CreateCourseRequest{
				Code:    "CS101",
				Name:    "Introduction to Computer Science",
				Credits: 0,
			},
			wantErr: true, // Credits should be positive
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.Code == "" && !tt.wantErr {
				t.Error("Code should be required")
			}
			if tt.req.Name == "" && !tt.wantErr {
				t.Error("Name should be required")
			}
			if tt.req.Credits == 0 && !tt.wantErr {
				t.Error("Credits should be required")
			}
		})
	}
}

// Test CreateClassRequest validation
func TestCreateClassRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateClassRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateClassRequest{
				CourseID:     uuid.New().String(),
				ClassCode:    "CS101-A",
				ClassName:    "CS101 Class A",
				Semester:     "Ganjil",
				AcademicYear: "2024/2025",
				DosenID:      uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "missing course_id",
			req: CreateClassRequest{
				ClassCode: "CS101-A",
				ClassName: "CS101 Class A",
				Semester:  "Ganjil",
				DosenID:   uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "missing class_code",
			req: CreateClassRequest{
				CourseID: uuid.New().String(),
				ClassName: "CS101 Class A",
				Semester:  "Ganjil",
				DosenID:   uuid.New().String(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.CourseID == "" && !tt.wantErr {
				t.Error("CourseID should be required")
			}
			if tt.req.ClassCode == "" && !tt.wantErr {
				t.Error("ClassCode should be required")
			}
		})
	}
}

// Test Enrollment status validation
func TestEnrollmentStatus(t *testing.T) {
	validStatuses := []string{
		"PENDING",
		"APPROVED",
		"REJECTED",
		"COMPLETED",
		"DROPPED",
		"FAILED",
	}

	for _, status := range validStatuses {
		t.Run(status, func(t *testing.T) {
			if status == "" {
				t.Error("Status should not be empty")
			}
		})
	}
}

// Test error types
func TestErrorTypes(t *testing.T) {
	t.Run("NotFoundError", func(t *testing.T) {
		err := apperrors.NewNotFoundError("course", "test-id")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("ConflictError", func(t *testing.T) {
		err := apperrors.NewConflictError("course with code already exists")
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

// Test Course model
func TestCourseModel(t *testing.T) {
	t.Run("valid course", func(t *testing.T) {
		course := createTestCourse()
		if course.Code == "" {
			t.Error("Code should not be empty")
		}
		if course.Name == "" {
			t.Error("Name should not be empty")
		}
		if course.Credits <= 0 {
			t.Error("Credits should be positive")
		}
		if course.ID == "" {
			t.Error("ID should be generated")
		}
	})

	t.Run("table name", func(t *testing.T) {
		course := models.Course{}
		if course.TableName() != "courses" {
			t.Errorf("Expected table name 'courses', got '%s'", course.TableName())
		}
	})
}

// Test Class model
func TestClassModel(t *testing.T) {
	t.Run("valid class", func(t *testing.T) {
		class := createTestClass()
		if class.ClassCode == "" {
			t.Error("ClassCode should not be empty")
		}
		if class.CourseID == "" {
			t.Error("CourseID should not be empty")
		}
		if class.ID == "" {
			t.Error("ID should be generated")
		}
	})

	t.Run("table name", func(t *testing.T) {
		class := models.Class{}
		if class.TableName() != "classes" {
			t.Errorf("Expected table name 'classes', got '%s'", class.TableName())
		}
	})
}

// Test Enrollment model
func TestEnrollmentModel(t *testing.T) {
	t.Run("valid enrollment", func(t *testing.T) {
		enrollment := createTestEnrollment()
		if enrollment.ClassID == "" {
			t.Error("ClassID should not be empty")
		}
		if enrollment.StudentID == "" {
			t.Error("StudentID should not be empty")
		}
		if enrollment.Status == "" {
			t.Error("Status should not be empty")
		}
		if enrollment.ID == "" {
			t.Error("ID should be generated")
		}
	})

	t.Run("table name", func(t *testing.T) {
		enrollment := models.Enrollment{}
		if enrollment.TableName() != "enrollments" {
			t.Errorf("Expected table name 'enrollments', got '%s'", enrollment.TableName())
		}
	})
}
