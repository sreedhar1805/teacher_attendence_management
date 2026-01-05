package model

import "time"

type Attendance struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	TeacherID uint       `json:"teacher_id" binding:"required"`
	Date      time.Time  `gorm:"type:date" json:"date"`
	Status    string     `json:"status" binding:"required"`
	CheckIn   *time.Time `json:"check_in,omitempty"`
	CheckOut  *time.Time `json:"check_out,omitempty"`
	Teacher   Teacher    `gorm:"foreignKey:TeacherID" json:"teacher"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type AttendanceDTO struct {
	TeacherID   uint       `json:"teacherId"`
	TeacherName string     `json:"teacherName"`
	CheckIn     *time.Time `json:"checkIn"`
	CheckOut    *time.Time `json:"checkOut"`
	Date        string     `json:"date"`
}

type AttendanceResponse struct {
	AttendanceList []AttendanceDTO `json:"attendanceList"`
	// CheckedInToday int             `json:"checkedInToday"`
}

type AttendanceRequest struct {
	TeacherID uint   `json:"teacher_id" binding:"required"`
	Status    string `json:"status" binding:"required"`
}
