package handler

import (
	"net/http"
	"strconv"
	"time"

	"school-teacher-management/internal/model"
	"school-teacher-management/internal/service"

	"school-teacher-management/internal/metrics"

	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	Service *service.AttendanceService
}

func NewAttendanceHandler(s *service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{Service: s}
}

// CreateAttendance godoc
// @Summary      Create attendance record
// @Description  Creates a new attendance entry (teacher_id and status are mandatory)
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Param        attendance  body      model.AttendanceRequest  true  "Attendance request"
// @Success      201         {object}  map[string]string
// @Failure      400         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /attendance [post]
func (h *AttendanceHandler) CreateAttendance(c *gin.Context) {
	var input model.AttendanceRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	metrics.AttendanceCreatedTotal.Inc()

	if err := h.Service.MarkAttendance(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	switch input.Status {
	case "checkIn":
		c.JSON(http.StatusCreated, gin.H{
			"message": "You have checked in successfully",
		})
	case "checkOut":
		c.JSON(http.StatusCreated, gin.H{
			"message": "You have checked out successfully",
		})
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid status value",
		})
	}
}

// GetAttendances godoc
// @Summary      Get all attendance records
// @Tags         attendance
// @Produce      json
// @Success      200  {array}   model.Attendance
// @Router       /attendance [get]
func (h *AttendanceHandler) GetAttendances(c *gin.Context) {
	list, err := h.Service.GetAttendances()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	result := []model.AttendanceDTO{}
	for _, att := range list {
		dto := model.AttendanceDTO{
			TeacherID:   att.Teacher.ID,
			TeacherName: att.Teacher.FirstName + " " + att.Teacher.LastName,
			CheckIn:     att.CheckIn,
			CheckOut:    att.CheckOut,
			Date:        att.Date.Format("02-01-2006"),
		}
		result = append(result, dto)
	}
	c.JSON(http.StatusOK, result)
}

// GetAttendanceByID godoc
// @Summary      Get attendance by ID
// @Tags         attendance
// @Param        id   path      int  true  "Attendance ID"
// @Produce      json
// @Success      200  {object}  model.Attendance
// @Failure      404  {object}  map[string]string
// @Router       /attendance/{id} [get]
func (h *AttendanceHandler) GetAttendanceByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID",
		})
		return
	}

	att, err := h.Service.GetAttendance(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{
			"error": "Attendance not found",
		})
		return
	}

	c.JSON(http.StatusOK, att)
}

// UpdateAttendance godoc
// @Summary      Update attendance
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Param        id          path      int               true  "Attendance ID"
// @Param        attendance  body      model.Attendance  true  "Updated attendance"
// @Success      200         {object}  model.Attendance
// @Failure      400         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /attendance/{id} [put]
func (h *AttendanceHandler) UpdateAttendance(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID",
		})
		return
	}

	var input model.Attendance
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	input.ID = uint(id)

	if err := h.Service.UpdateAttendance(&input); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, input)
}

// DeleteAttendance godoc
// @Summary      Delete attendance
// @Tags         attendance
// @Param        id   path  int  true  "Attendance ID"
// @Success      204  "No Content"
// @Failure      404  {object}  map[string]string
// @Router       /attendance/{id} [delete]
func (h *AttendanceHandler) DeleteAttendance(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID",
		})
		return
	}

	if err := h.Service.DeleteAttendance(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, map[string]string{
			"error": "Attendance not found",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAttendance godoc
// @Summary Get attendance by teacher
// @Description Get attendance for a teacher filtered by month/year and checked-in count today
// @Tags attendance
// @Param teacherId query int true "Teacher ID"
// @Param month query int false "Month (1-12), default current month"
// @Param year query int false "Year, default current year"
// @Success 200 {object} model.AttendanceResponse
// @Router /attendanceByDate [get]
func (h *AttendanceHandler) GetAttendanceByDate(c *gin.Context) {
	teacherIDStr := c.Query("teacherId")
	if teacherIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "teacherId is required"})
		return
	}

	teacherID, err := strconv.ParseUint(teacherIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid teacherId"})
		return
	}

	// Get month/year or default to current
	monthStr := c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month())))
	yearStr := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))

	monthInt, _ := strconv.Atoi(monthStr)
	yearInt, _ := strconv.Atoi(yearStr)

	resp, err := h.Service.GetAttendanceByTeacherMonth(uint(teacherID), time.Month(monthInt), yearInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetAttendance godoc
// @Summary Get attendance for a date
// @Description Get attendance filtered by date/month/year (teacher independent)
// @Tags attendance
// @Param date  query int false "Day of month (1-31)"
// @Param month query int false "Month (1-12)"
// @Param year  query int false "Year (YYYY)"
// @Success 200 {object} model.AttendanceResponse
// @Router /attendanceByFilterDate [get]
func (h *AttendanceHandler) GetAttendanceByFilterDate(c *gin.Context) {

	now := time.Now()

	dayStr := c.DefaultQuery("date", strconv.Itoa(now.Day()))
	monthStr := c.DefaultQuery("month", strconv.Itoa(int(now.Month())))
	yearStr := c.DefaultQuery("year", strconv.Itoa(now.Year()))

	day, err := strconv.Atoi(dayStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
		return
	}

	monthInt, err := strconv.Atoi(monthStr)
	if monthInt < 1 || monthInt > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "month must be 1-12"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month"})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}

	date := time.Date(
		year,
		time.Month(monthInt),
		day,
		0, 0, 0, 0,
		time.Local,
	)

	resp, err := h.Service.GetAttendanceByMonthAndDate(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
