package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ansarctica/domashka4/internal/models"
	"github.com/labstack/echo/v4"
)

type AttendanceInput struct {
	SubjectName string `json:"subject_name"`
	VisitDay    string `json:"visit_day"`
	Visited     bool   `json:"visited"`
	StudentID   int    `json:"student_id"`
}

// GetAttendance retrieves attendance records
// @Summary Get attendance
// @Description Get attendance records filtered by student ID or subject name
// @Tags Attendance
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param student_id query int false "Filter by Student ID"
// @Param subject_name query string false "Filter by Subject Name"
// @Success 200 {array} models.Attendance
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /attendance [get]
func (h *Handler) GetAttendance(c echo.Context) error {
	var params struct {
		StudentID   *int    `query:"student_id"`
		SubjectName *string `query:"subject_name"`
	}

	if err := c.Bind(&params); err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	attendanceList, err := h.service.GetAttendance(c.Request().Context(), params.StudentID, params.SubjectName)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, attendanceList)
}

// CreateAttendance records a new attendance entry
// @Summary Record attendance
// @Description Create a new attendance record for a student
// @Tags Attendance
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body handlers.AttendanceInput true "Attendance Data"
// @Success 201 {object} map[string]int "Returns created ID"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /attendance [post]
func (h *Handler) CreateAttendance(c echo.Context) error {
	var input AttendanceInput

	if err := c.Bind(&input); err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	parsedDate, err := time.Parse("02.01.2006", input.VisitDay)
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	attendance := &models.Attendance{
		SubjectName: input.SubjectName,
		VisitDay:    parsedDate,
		Visited:     input.Visited,
		StudentID:   input.StudentID,
	}

	id, err := h.service.NewAttendance(c.Request().Context(), attendance)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, map[string]int{"id": id})
}

// UpdateAttendance modifies an attendance record
// @Summary Update attendance
// @Description Update an existing attendance record
// @Tags Attendance
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Attendance ID"
// @Param input body handlers.AttendanceInput true "Updated Attendance Data"
// @Success 200 {object} map[string]string "Returns status"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /attendance/{id} [patch]
func (h *Handler) UpdateAttendance(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	var input AttendanceInput
	if err := c.Bind(&input); err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	parsedDate, err := time.Parse("02.01.2006", input.VisitDay)
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	attendance := &models.Attendance{
		ID:          id,
		SubjectName: input.SubjectName,
		VisitDay:    parsedDate,
		Visited:     input.Visited,
		StudentID:   input.StudentID,
	}

	if err := h.service.UpdateAttendance(c.Request().Context(), attendance); err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "updated"})
}

// DeleteAttendance removes an attendance record
// @Summary Delete attendance
// @Description Remove an attendance record from the database
// @Tags Attendance
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Attendance ID"
// @Success 200 {object} map[string]string "Returns status"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /attendance/{id} [delete]
func (h *Handler) DeleteAttendance(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	if err := h.service.DeleteAttendance(c.Request().Context(), id); err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}
