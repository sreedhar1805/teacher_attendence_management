package service

import (
	"school-teacher-management/internal/model"
	"school-teacher-management/internal/repository"
)

type TeacherService struct {
	Repo *repository.TeacherRepository
}

func NewTeacherService(repo *repository.TeacherRepository) *TeacherService {
	return &TeacherService{Repo: repo}
}

func (s *TeacherService) CreateTeacher(teacher *model.Teacher) error {
	return s.Repo.Create(teacher)
}

func (s *TeacherService) UpdateTeacher(teacher *model.Teacher) error {
	return s.Repo.Update(teacher)
}

func (s *TeacherService) GetTeacher(id uint) (*model.Teacher, error) {
	return s.Repo.GetByID(id)
}

func (s *TeacherService) SearchTeachers(q string, subject string) ([]model.Teacher, error) {
	return s.Repo.SearchAllFields(q, subject)
}

func (s *TeacherService) CreateTeachers(req []model.TeacherRequest) error {
	var teachers []model.Teacher

	for _, t := range req {
		teachers = append(teachers, model.Teacher{
			FirstName: t.FirstName,
			LastName:  t.LastName,
			Email:     t.Email,
			Subject:   t.Subject,
			Phone:     t.Phone,
		})
	}

	return s.Repo.BulkCreate(teachers)
}
