package handlers

import (
	"net/http"
	"strconv"

	"github.com/ansarctica/domashka4/internal/models"
	"github.com/labstack/echo/v4"
)

// GetStudents retrieves a list of students
// @Summary Get all students
// @Description Retrieve students with optional filtering by group, major, and course year
// @Tags Students
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param group_id query int false "Filter by Group ID"
// @Param major query string false "Filter by Major"
// @Param course_year query int false "Filter by Course Year"
// @Param limit query int false "Limit number of results (default 20)"
// @Param offset query int false "Offset for pagination"
// @Success 200 {array} models.Student
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /students [get]
func (h *Handler) GetStudents(c echo.Context) error {
	var params struct {
		GroupID    *int    `query:"group_id"`
		Major      *string `query:"major"`
		CourseYear *int    `query:"course_year"`
		Limit      int     `query:"limit"`
		Offset     int     `query:"offset"`
	}

	if err := c.Bind(&params); err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	if params.Limit == 0 {
		params.Limit = 20
	}

	filter := models.StudentFilter{
		GroupID:    params.GroupID,
		Major:      params.Major,
		CourseYear: params.CourseYear,
		Limit:      params.Limit,
		Offset:     params.Offset,
	}

	students, err := h.service.GetAllStudents(c.Request().Context(), filter)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, students)
}

// GetStudent retrieves a specific student
// @Summary Get a student
// @Description Get details of a student by their ID
// @Tags Students
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} models.Student
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /students/{id} [get]
func (h *Handler) GetStudent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	student, err := h.service.GetStudent(c.Request().Context(), id)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, student)
}

// CreateStudent adds a new student
// @Summary Create a student
// @Description Add a new student record to the database
// @Tags Students
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body models.Student true "Student Data"
// @Success 201 {object} map[string]int "Returns created ID"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /students [post]
func (h *Handler) CreateStudent(c echo.Context) error {
	var student models.Student
	if err := c.Bind(&student); err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	student.ID = 0

	id, err := h.service.CreateStudent(c.Request().Context(), &student)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, map[string]int{"id": id})
}

// UpdateStudent updates an existing student
// @Summary Update a student
// @Description Update details of an existing student
// @Tags Students
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Param input body models.Student true "Updated Student Data"
// @Success 200 {object} map[string]string "Returns status"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /students/{id} [patch]
func (h *Handler) UpdateStudent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	var student models.Student
	if err := c.Bind(&student); err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}
	student.ID = id

	if err := h.service.UpdateStudent(c.Request().Context(), &student); err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "updated"})
}

// DeleteStudent removes a student
// @Summary Delete a student
// @Description Remove a student record from the database
// @Tags Students
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} map[string]string "Returns status"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /students/{id} [delete]
func (h *Handler) DeleteStudent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	if err := h.service.DeleteStudent(c.Request().Context(), id); err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

// GetStudentGPA calculates a student's GPA
// @Summary Get Student GPA
// @Description Calculate and return the GPA for a specific student
// @Tags Students
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /students/{id}/gpa [get]
func (h *Handler) GetStudentGPA(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	gpa, err := h.service.GetGPA(c.Request().Context(), id)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"student_id": id,
		"gpa":        gpa,
	})
}
