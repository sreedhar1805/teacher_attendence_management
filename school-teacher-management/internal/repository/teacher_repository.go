package repository

import (
	"school-teacher-management/internal/model"

	"gorm.io/gorm"
)

type TeacherRepository struct {
	DB *gorm.DB
}

func NewTeacherRepository(db *gorm.DB) *TeacherRepository {
	return &TeacherRepository{DB: db}
}

func (r *TeacherRepository) Create(teacher *model.Teacher) error {
	return r.DB.Create(teacher).Error
}

func (r *TeacherRepository) Update(teacher *model.Teacher) error {
	return r.DB.Save(teacher).Error
}

func (r *TeacherRepository) GetByID(id uint) (*model.Teacher, error) {
	var teacher model.Teacher
	err := r.DB.First(&teacher, id).Error
	return &teacher, err
}

func (r *TeacherRepository) SearchAllFields(q string, subject string) ([]model.Teacher, error) {
	var teachers []model.Teacher
	db := r.DB.Model(&model.Teacher{})

	if q != "" {
		likePattern := "%" + q + "%"
		db = db.Where(
			"first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR subject ILIKE ?",
			likePattern, likePattern, likePattern, likePattern,
		)
	}

	if subject != "" {
		db = db.Where(
			"subject ILIKE ?",
			subject,
		)
	}

	err := db.Find(&teachers).Error
	return teachers, err
}
