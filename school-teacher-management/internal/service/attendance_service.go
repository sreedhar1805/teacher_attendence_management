package service

import (
	"errors"
	"school-teacher-management/internal/metrics"
	"school-teacher-management/internal/model"
	"school-teacher-management/internal/repository"
	"time"
)

type AttendanceService struct {
	Repo *repository.AttendanceRepository
}

func NewAttendanceService(repo *repository.AttendanceRepository) *AttendanceService {
	return &AttendanceService{Repo: repo}
}

func (s *AttendanceService) CreateAttendance(att *model.Attendance) error {
	return s.Repo.Create(att)
}

func (s *AttendanceService) GetAttendances() ([]model.Attendance, error) {
	return s.Repo.GetAll()
}

func (s *AttendanceService) GetAttendance(id uint) (*model.Attendance, error) {
	return s.Repo.GetByID(id)
}

func (s *AttendanceService) UpdateAttendance(att *model.Attendance) error {
	return s.Repo.Update(att)
}

func (s *AttendanceService) DeleteAttendance(id uint) error {
	return s.Repo.Delete(id)
}

func (s *AttendanceService) MarkAttendance(input *model.Attendance) error {
	var existing model.Attendance

	today := time.Now().Truncate(24 * time.Hour)

	// Check if attendance already exists for today
	err := s.Repo.FindByTeacherAndDate(
		input.TeacherID,
		today,
		&existing,
	)

	now := time.Now()

	if err != nil { // No record exists → allow only CHECK-IN
		if input.Status != "checkIn" {
			return errors.New("check-in required before check-out")
		}

		input.Date = today
		input.CheckIn = &now
		return s.Repo.Create(input)
	}

	// Record exists
	if input.Status == "checkIn" {
		return errors.New("already checked in for today")
	}

	if input.Status == "checkOut" {
		if existing.CheckIn == nil {
			metrics.AttendanceCheckInTotal.Inc()
			return errors.New("cannot checkout without check-in")
		}
		if existing.CheckOut != nil {
			metrics.AttendanceCheckOutTotal.Inc()
			return errors.New("already checked out")
		}

		existing.CheckOut = &now
		return s.Repo.Update(&existing)
	}

	return errors.New("invalid status value")
}

func (s *AttendanceService) GetAttendanceByTeacherMonth(teacherID uint, month time.Month, year int) (*model.AttendanceResponse, error) {
	// Get filtered attendance
	attList, err := s.Repo.FindByTeacherAndMonth(teacherID, month, year)
	if err != nil {
		return nil, err
	}

	// Map to DTO
	result := []model.AttendanceDTO{}
	for _, att := range attList {
		dto := model.AttendanceDTO{
			TeacherID:   att.Teacher.ID,
			TeacherName: att.Teacher.FirstName + " " + att.Teacher.LastName,
			CheckIn:     att.CheckIn,
			CheckOut:    att.CheckOut,
			Date:        att.Date.Format("02-01-2006"),
		}
		result = append(result, dto)
	}

	// Count checked-in today
	// countToday, err := s.Repo.CountCheckedInToday()
	// if err != nil {
	// 	return nil, err
	// }

	resp := &model.AttendanceResponse{
		AttendanceList: result,
		// CheckedInToday: int(countToday),
	}

	return resp, nil
}

func (s *AttendanceService) GetAttendanceByMonthAndDate(date time.Time) (*model.AttendanceResponse, error) {

	month := date.Month()
	year := date.Year()

	attList, err := s.Repo.FindByMonthAndDate(month, year, date)
	if err != nil {
		return nil, err
	}

	result := []model.AttendanceDTO{}

	for _, att := range attList {
		dto := model.AttendanceDTO{
			TeacherID:   att.Teacher.ID,
			TeacherName: att.Teacher.FirstName + " " + att.Teacher.LastName,
			CheckIn:     att.CheckIn,
			CheckOut:    att.CheckOut,
			Date:        att.Date.Format("02-01-2006"),
		}
		result = append(result, dto)
	}

	return &model.AttendanceResponse{
		AttendanceList: result,
	}, nil
}
