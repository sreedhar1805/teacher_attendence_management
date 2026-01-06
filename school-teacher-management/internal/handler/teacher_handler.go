package handler

import (
	"net/http"
	"strconv"

	"school-teacher-management/internal/model"
	"school-teacher-management/internal/service"

	"school-teacher-management/internal/metrics"

	"github.com/gin-gonic/gin"
)

type TeacherHandler struct {
	Service *service.TeacherService
}

func NewTeacherHandler(s *service.TeacherService) *TeacherHandler {
	return &TeacherHandler{Service: s}
}

// CreateTeacher godoc
// @Summary      Create teacher
// @Description  Create a new teacher
// @Tags         teachers
// @Accept       json
// @Produce      json
// @Param        teacher  body      model.Teacher  true  "Teacher data"
// @Success      201      {object}  model.Teacher
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /teachers [post]
func (h *TeacherHandler) CreateTeacher(c *gin.Context) {
	var input model.Teacher

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := h.Service.CreateTeacher(&input); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	metrics.TeachersCreatedTotal.Inc()
	metrics.TeachersTotal.Inc()

	c.JSON(http.StatusCreated, input)
}

// UpdateTeacher godoc
// @Summary      Update teacher
// @Description  Update a teacher by ID
// @Tags         teachers
// @Accept       json
// @Produce      json
// @Param        id       path      int            true  "Teacher ID"
// @Param        teacher  body      model.Teacher  true  "Teacher data"
// @Success      200      {object}  model.Teacher
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /teachers/{id} [put]
func (h *TeacherHandler) UpdateTeacher(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID",
		})
		return
	}

	var input model.Teacher
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	input.ID = uint(id)

	if err := h.Service.UpdateTeacher(&input); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, input)
}

// SearchTeachers godoc
// @Summary      Search teachers
// @Description  Search teachers across first name, last name, email, subject
// @Tags         teachers
// @Produce      json
// @Param        q   query     string  false  "Search keyword"
// @Param        subject   query     string  false  "Search keyword"
// @Success      200 {array}   model.Teacher
// @Failure      500 {object}  map[string]string
// @Router       /teachers [get]
func (h *TeacherHandler) SearchTeachers(c *gin.Context) {
	q := c.Query("q")
	subject := c.Query("subject")

	teachers, err := h.Service.SearchTeachers(q, subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, teachers)
}

// GetTeacherByID godoc
// @Summary      Get teacher by ID
// @Tags         teachers
// @Param        id   path      int  true  "Teacher ID"
// @Produce      json
// @Success      200  {object}  model.Teacher
// @Failure      404  {object}  map[string]string
// @Router       /teachers/{id} [get]
func (h *TeacherHandler) GetTeacherByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID",
		})
		return
	}

	teacher, err := h.Service.GetTeacher(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{
			"error": "Teacher not found",
		})
		return
	}

	c.JSON(http.StatusOK, teacher)
}

// CreateTeachers godoc
// @Summary Create multiple teachers
// @Tags teachers
// @Accept json
// @Produce json
// @Param teachers body []model.TeacherRequest true "List of teachers"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /teachers/bulk [post]
func (h *TeacherHandler) CreateTeachers(c *gin.Context) {
	var input []model.TeacherRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(input) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "teacher list cannot be empty",
		})
		return
	}

	if err := h.Service.CreateTeachers(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Teachers created successfully",
		"count":   len(input),
	})
}
