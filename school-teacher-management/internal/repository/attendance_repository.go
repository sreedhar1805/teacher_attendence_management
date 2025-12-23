package repository

import (
	"school-teacher-management/internal/model"
	"time"

	"gorm.io/gorm"
)

type AttendanceRepository struct {
	DB *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{DB: db}
}

func (r *AttendanceRepository) Create(att *model.Attendance) error {
	return r.DB.Create(att).Error
}

func (r *AttendanceRepository) GetAll() ([]model.Attendance, error) {
	var list []model.Attendance
	err := r.DB.Preload("Teacher").Find(&list).Error
	return list, err
}

func (r *AttendanceRepository) GetByID(id uint) (*model.Attendance, error) {
	var att model.Attendance
	err := r.DB.Preload("Teacher").First(&att, id).Error
	return &att, err
}

func (r *AttendanceRepository) Update(att *model.Attendance) error {
	return r.DB.Save(att).Error
}

func (r *AttendanceRepository) Delete(id uint) error {
	return r.DB.Delete(&model.Attendance{}, id).Error
}

func (r *AttendanceRepository) FindByTeacherAndDate(
	teacherID uint,
	date time.Time,
	attendance *model.Attendance,
) error {
	return r.DB.
		Where("teacher_id = ? AND date = ?", teacherID, date).
		First(attendance).Error
}

// FindByTeacherAndMonth finds all attendance for a teacher in a given month/year
func (r *AttendanceRepository) FindByTeacherAndMonth(
	teacherID uint,
	month time.Month,
	year int,
) ([]model.Attendance, error) {
	var list []model.Attendance
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0) // next month
	err := r.DB.
		Where("teacher_id = ? AND date >= ? AND date < ?", teacherID, start, end).
		Preload("Teacher").
		Find(&list).Error
	return list, err
}

// CountCheckedInToday counts how many teachers have checked in today
func (r *AttendanceRepository) CountCheckedInToday() (int64, error) {
	today := time.Now().Truncate(24 * time.Hour)
	var count int64
	err := r.DB.
		Model(&model.Attendance{}).
		Where("check_in IS NOT NULL AND date = ?", today).
		Count(&count).Error
	return count, err
}
