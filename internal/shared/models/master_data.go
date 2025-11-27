package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StudyProgram represents a study program (Program Studi)
type StudyProgram struct {
	ID            string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Code          string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name          string    `gorm:"type:varchar(255);not null" json:"name"`
	NameEn        string    `gorm:"type:varchar(255)" json:"name_en"`
	Faculty       string    `gorm:"type:varchar(255)" json:"faculty"`
	DegreeLevel   string    `gorm:"type:varchar(50)" json:"degree_level"` // S1, S2, S3, D3, etc.
	Accreditation string    `gorm:"type:varchar(10)" json:"accreditation"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name
func (StudyProgram) TableName() string {
	return "study_programs"
}

// BeforeCreate hook
func (s *StudyProgram) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

// AcademicPeriod represents an academic period (Semester)
type AcademicPeriod struct {
	ID                string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Code              string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name              string    `gorm:"type:varchar(255);not null" json:"name"`
	AcademicYear      string    `gorm:"type:varchar(20);not null" json:"academic_year"`
	SemesterType      string    `gorm:"type:varchar(20);not null" json:"semester_type"` // GANJIL, GENAP, PENDEK
	StartDate         time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate           time.Time `gorm:"type:date;not null" json:"end_date"`
	RegistrationStart *time.Time `gorm:"type:date" json:"registration_start,omitempty"`
	RegistrationEnd   *time.Time `gorm:"type:date" json:"registration_end,omitempty"`
	IsActive          bool      `gorm:"default:false" json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name
func (AcademicPeriod) TableName() string {
	return "academic_periods"
}

// BeforeCreate hook
func (a *AcademicPeriod) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

// Room represents a room (Ruangan)
type Room struct {
	ID         string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Code       string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name       string    `gorm:"type:varchar(255);not null" json:"name"`
	Building   string    `gorm:"type:varchar(255)" json:"building"`
	Floor      *int      `gorm:"type:integer" json:"floor,omitempty"`
	Capacity   *int      `gorm:"type:integer" json:"capacity,omitempty"`
	RoomType   string    `gorm:"type:varchar(50)" json:"room_type"` // classroom, lab, auditorium, etc.
	Facilities string    `gorm:"type:text" json:"facilities"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name
func (Room) TableName() string {
	return "rooms"
}

// BeforeCreate hook
func (r *Room) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

