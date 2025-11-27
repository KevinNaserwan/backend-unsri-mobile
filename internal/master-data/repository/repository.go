package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"unsri-backend/internal/shared/models"
)

// MasterDataRepository handles master data operations
type MasterDataRepository struct {
	db *gorm.DB
}

// NewMasterDataRepository creates a new master data repository
func NewMasterDataRepository(db *gorm.DB) *MasterDataRepository {
	return &MasterDataRepository{db: db}
}

// ========== Study Program Methods ==========

// CreateStudyProgram creates a new study program
func (r *MasterDataRepository) CreateStudyProgram(ctx context.Context, studyProgram *models.StudyProgram) error {
	return r.db.WithContext(ctx).Create(studyProgram).Error
}

// GetStudyProgramByID gets a study program by ID
func (r *MasterDataRepository) GetStudyProgramByID(ctx context.Context, id string) (*models.StudyProgram, error) {
	var studyProgram models.StudyProgram
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&studyProgram).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("study program not found")
		}
		return nil, err
	}
	return &studyProgram, nil
}

// GetStudyProgramByCode gets a study program by code
func (r *MasterDataRepository) GetStudyProgramByCode(ctx context.Context, code string) (*models.StudyProgram, error) {
	var studyProgram models.StudyProgram
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&studyProgram).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("study program not found")
		}
		return nil, err
	}
	return &studyProgram, nil
}

// GetAllStudyPrograms gets all study programs with filters
func (r *MasterDataRepository) GetAllStudyPrograms(ctx context.Context, faculty *string, isActive *bool, limit, offset int) ([]models.StudyProgram, int64, error) {
	var studyPrograms []models.StudyProgram
	var total int64

	query := r.db.WithContext(ctx).Model(&models.StudyProgram{})

	if faculty != nil && *faculty != "" {
		query = query.Where("faculty = ?", *faculty)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Order("code ASC").Find(&studyPrograms).Error; err != nil {
		return nil, 0, err
	}

	return studyPrograms, total, nil
}

// UpdateStudyProgram updates a study program
func (r *MasterDataRepository) UpdateStudyProgram(ctx context.Context, studyProgram *models.StudyProgram) error {
	return r.db.WithContext(ctx).Save(studyProgram).Error
}

// DeleteStudyProgram soft deletes a study program
func (r *MasterDataRepository) DeleteStudyProgram(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.StudyProgram{}, "id = ?", id).Error
}

// ========== Academic Period Methods ==========

// CreateAcademicPeriod creates a new academic period
func (r *MasterDataRepository) CreateAcademicPeriod(ctx context.Context, academicPeriod *models.AcademicPeriod) error {
	return r.db.WithContext(ctx).Create(academicPeriod).Error
}

// GetAcademicPeriodByID gets an academic period by ID
func (r *MasterDataRepository) GetAcademicPeriodByID(ctx context.Context, id string) (*models.AcademicPeriod, error) {
	var academicPeriod models.AcademicPeriod
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&academicPeriod).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("academic period not found")
		}
		return nil, err
	}
	return &academicPeriod, nil
}

// GetAcademicPeriodByCode gets an academic period by code
func (r *MasterDataRepository) GetAcademicPeriodByCode(ctx context.Context, code string) (*models.AcademicPeriod, error) {
	var academicPeriod models.AcademicPeriod
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&academicPeriod).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("academic period not found")
		}
		return nil, err
	}
	return &academicPeriod, nil
}

// GetActiveAcademicPeriod gets the active academic period
func (r *MasterDataRepository) GetActiveAcademicPeriod(ctx context.Context) (*models.AcademicPeriod, error) {
	var academicPeriod models.AcademicPeriod
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).First(&academicPeriod).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no active academic period found")
		}
		return nil, err
	}
	return &academicPeriod, nil
}

// GetAllAcademicPeriods gets all academic periods with filters
func (r *MasterDataRepository) GetAllAcademicPeriods(ctx context.Context, academicYear *string, semesterType *string, isActive *bool, limit, offset int) ([]models.AcademicPeriod, int64, error) {
	var academicPeriods []models.AcademicPeriod
	var total int64

	query := r.db.WithContext(ctx).Model(&models.AcademicPeriod{})

	if academicYear != nil && *academicYear != "" {
		query = query.Where("academic_year = ?", *academicYear)
	}
	if semesterType != nil && *semesterType != "" {
		query = query.Where("semester_type = ?", *semesterType)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Order("academic_year DESC, semester_type ASC").Find(&academicPeriods).Error; err != nil {
		return nil, 0, err
	}

	return academicPeriods, total, nil
}

// UpdateAcademicPeriod updates an academic period
func (r *MasterDataRepository) UpdateAcademicPeriod(ctx context.Context, academicPeriod *models.AcademicPeriod) error {
	return r.db.WithContext(ctx).Save(academicPeriod).Error
}

// DeleteAcademicPeriod soft deletes an academic period
func (r *MasterDataRepository) DeleteAcademicPeriod(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.AcademicPeriod{}, "id = ?", id).Error
}

// ========== Room Methods ==========

// CreateRoom creates a new room
func (r *MasterDataRepository) CreateRoom(ctx context.Context, room *models.Room) error {
	return r.db.WithContext(ctx).Create(room).Error
}

// GetRoomByID gets a room by ID
func (r *MasterDataRepository) GetRoomByID(ctx context.Context, id string) (*models.Room, error) {
	var room models.Room
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("room not found")
		}
		return nil, err
	}
	return &room, nil
}

// GetRoomByCode gets a room by code
func (r *MasterDataRepository) GetRoomByCode(ctx context.Context, code string) (*models.Room, error) {
	var room models.Room
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("room not found")
		}
		return nil, err
	}
	return &room, nil
}

// GetAllRooms gets all rooms with filters
func (r *MasterDataRepository) GetAllRooms(ctx context.Context, building *string, roomType *string, isActive *bool, limit, offset int) ([]models.Room, int64, error) {
	var rooms []models.Room
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Room{})

	if building != nil && *building != "" {
		query = query.Where("building = ?", *building)
	}
	if roomType != nil && *roomType != "" {
		query = query.Where("room_type = ?", *roomType)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Order("building ASC, code ASC").Find(&rooms).Error; err != nil {
		return nil, 0, err
	}

	return rooms, total, nil
}

// UpdateRoom updates a room
func (r *MasterDataRepository) UpdateRoom(ctx context.Context, room *models.Room) error {
	return r.db.WithContext(ctx).Save(room).Error
}

// DeleteRoom soft deletes a room
func (r *MasterDataRepository) DeleteRoom(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Room{}, "id = ?", id).Error
}

