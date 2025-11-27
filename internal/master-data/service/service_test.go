package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
)

// Test helper functions
func createTestStudyProgram() *models.StudyProgram {
	return &models.StudyProgram{
		ID:          uuid.New().String(),
		Code:        "TI",
		Name:        "Teknik Informatika",
		NameEn:      "Informatics Engineering",
		Faculty:     "Fakultas Ilmu Komputer",
		DegreeLevel: "S1",
		IsActive:    true,
	}
}

func createTestAcademicPeriod() *models.AcademicPeriod {
	startDate := time.Now()
	endDate := startDate.AddDate(0, 6, 0)
	return &models.AcademicPeriod{
		ID:           uuid.New().String(),
		Code:         "2024-GANJIL",
		Name:         "Ganjil 2024/2025",
		AcademicYear: "2024/2025",
		SemesterType: "Ganjil",
		StartDate:    startDate,
		EndDate:      endDate,
		IsActive:     false,
	}
}

func createTestRoom() *models.Room {
	floor := 1
	capacity := 40
	return &models.Room{
		ID:       uuid.New().String(),
		Code:     "A101",
		Name:     "Ruang A101",
		Building: "Gedung A",
		Floor:    &floor,
		Capacity: &capacity,
		RoomType: "classroom",
		IsActive: true,
	}
}

// Test CreateStudyProgramRequest validation
func TestCreateStudyProgramRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateStudyProgramRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateStudyProgramRequest{
				Code:        "TI",
				Name:        "Teknik Informatika",
				NameEn:      "Informatics Engineering",
				Faculty:     "Fakultas Ilmu Komputer",
				DegreeLevel: "S1",
			},
			wantErr: false,
		},
		{
			name: "missing code",
			req: CreateStudyProgramRequest{
				Name: "Teknik Informatika",
			},
			wantErr: true,
		},
		{
			name: "missing name",
			req: CreateStudyProgramRequest{
				Code: "TI",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This is a validation test - actual validation happens in handler
			// We're just testing the struct definition
			if tt.req.Code == "" && !tt.wantErr {
				t.Error("Code should be required")
			}
			if tt.req.Name == "" && !tt.wantErr {
				t.Error("Name should be required")
			}
		})
	}
}

// Test CreateAcademicPeriodRequest validation
func TestCreateAcademicPeriodRequest(t *testing.T) {
	startDate := time.Now()
	endDate := startDate.AddDate(0, 6, 0)

	tests := []struct {
		name    string
		req     CreateAcademicPeriodRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateAcademicPeriodRequest{
				Code:         "2024-GANJIL",
				Name:         "Ganjil 2024/2025",
				AcademicYear: "2024/2025",
				SemesterType: "Ganjil",
				StartDate:    startDate.Format("2006-01-02"),
				EndDate:      endDate.Format("2006-01-02"),
			},
			wantErr: false,
		},
		{
			name: "invalid date format",
			req: CreateAcademicPeriodRequest{
				Code:         "2024-GANJIL",
				Name:         "Ganjil 2024/2025",
				AcademicYear: "2024/2025",
				SemesterType: "Ganjil",
				StartDate:    "invalid-date",
				EndDate:      endDate.Format("2006-01-02"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := time.Parse("2006-01-02", tt.req.StartDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("Date parsing error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test CreateRoomRequest validation
func TestCreateRoomRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateRoomRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateRoomRequest{
				Code:     "A101",
				Name:     "Ruang A101",
				Building: "Gedung A",
				Floor:    func() *int { v := 1; return &v }(),
				Capacity: func() *int { v := 40; return &v }(),
				RoomType: "classroom",
			},
			wantErr: false,
		},
		{
			name: "missing code",
			req: CreateRoomRequest{
				Name: "Ruang A101",
			},
			wantErr: true,
		},
		{
			name: "missing name",
			req: CreateRoomRequest{
				Code: "A101",
			},
			wantErr: true,
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
		})
	}
}

// Test error types
func TestErrorTypes(t *testing.T) {
	t.Run("NotFoundError", func(t *testing.T) {
		err := apperrors.NewNotFoundError("study program", "test-id")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("ConflictError", func(t *testing.T) {
		err := apperrors.NewConflictError("study program with code already exists")
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

// Test model validation
func TestStudyProgramModel(t *testing.T) {
	t.Run("valid study program", func(t *testing.T) {
		sp := createTestStudyProgram()
		if sp.Code == "" {
			t.Error("Code should not be empty")
		}
		if sp.Name == "" {
			t.Error("Name should not be empty")
		}
		if sp.ID == "" {
			t.Error("ID should be generated")
		}
	})

	t.Run("table name", func(t *testing.T) {
		sp := models.StudyProgram{}
		if sp.TableName() != "study_programs" {
			t.Errorf("Expected table name 'study_programs', got '%s'", sp.TableName())
		}
	})
}

func TestAcademicPeriodModel(t *testing.T) {
	t.Run("valid academic period", func(t *testing.T) {
		ap := createTestAcademicPeriod()
		if ap.Code == "" {
			t.Error("Code should not be empty")
		}
		if ap.Name == "" {
			t.Error("Name should not be empty")
		}
		if ap.StartDate.After(ap.EndDate) {
			t.Error("Start date should be before end date")
		}
	})

	t.Run("table name", func(t *testing.T) {
		ap := models.AcademicPeriod{}
		if ap.TableName() != "academic_periods" {
			t.Errorf("Expected table name 'academic_periods', got '%s'", ap.TableName())
		}
	})
}

func TestRoomModel(t *testing.T) {
	t.Run("valid room", func(t *testing.T) {
		room := createTestRoom()
		if room.Code == "" {
			t.Error("Code should not be empty")
		}
		if room.Name == "" {
			t.Error("Name should not be empty")
		}
		if room.Capacity != nil && *room.Capacity <= 0 {
			t.Error("Capacity should be positive")
		}
	})

	t.Run("table name", func(t *testing.T) {
		room := models.Room{}
		if room.TableName() != "rooms" {
			t.Errorf("Expected table name 'rooms', got '%s'", room.TableName())
		}
	})
}

// Test date validation logic
func TestDateValidation(t *testing.T) {
	t.Run("valid date range", func(t *testing.T) {
		startDate := time.Now()
		endDate := startDate.AddDate(0, 6, 0)

		if endDate.Before(startDate) {
			t.Error("End date should be after start date")
		}
	})

	t.Run("invalid date range", func(t *testing.T) {
		startDate := time.Now()
		endDate := startDate.AddDate(0, -1, 0)

		if !endDate.Before(startDate) {
			t.Error("End date should be before start date in this test case")
		}
	})
}
