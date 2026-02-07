package handlers

import (
	"net/http"
	"time"

	"github.com/ansarctica/domashka4/internal/models"
	"github.com/labstack/echo/v4"
)

type AssignmentInput struct {
	Name        string `json:"name"`
	SubjectName string `json:"subject_name"`
	Weight      int    `json:"weight"`
	Date        string `json:"date"`
}

type GradeInput struct {
	StudentID    int `json:"student_id"`
	AssignmentID int `json:"assignment_id"`
	Mark         int `json:"mark"`
}

// GetAssignments retrieves assignments
// @Summary Get assignments
// @Description Get a list of assignments, optionally filtered by subject
// @Tags Assignments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param subject_name query string false "Filter by Subject Name"
// @Success 200 {array} models.Assignment
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /assignments [get]
func (h *Handler) GetAssignments(c echo.Context) error {
	var params struct {
		SubjectName *string `query:"subject_name"`
	}

	if err := c.Bind(&params); err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	assignments, err := h.service.GetAssignments(c.Request().Context(), params.SubjectName)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, assignments)
}

// CreateAssignment creates a new assignment
// @Summary Create assignment
// @Description Add a new assignment for a subject
// @Tags Assignments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body handlers.AssignmentInput true "Assignment Data"
// @Success 201 {object} map[string]int "Returns created ID"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /assignments [post]
func (h *Handler) CreateAssignment(c echo.Context) error {
	var input AssignmentInput

	if err := c.Bind(&input); err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	parsedDate, err := time.Parse("02.01.2006", input.Date)
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	assignment := &models.Assignment{
		Name:        input.Name,
		SubjectName: input.SubjectName,
		Weight:      input.Weight,
		Date:        parsedDate,
	}

	id, err := h.service.NewAssignment(c.Request().Context(), assignment)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, map[string]int{"id": id})
}

// CreateGrade records a grade for a student
// @Summary Create grade
// @Description Record a grade for a specific assignment
// @Tags Grades
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body handlers.GradeInput true "Grade Data"
// @Success 201 {object} map[string]int "Returns created ID"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /grades [post]
func (h *Handler) CreateGrade(c echo.Context) error {
	var input GradeInput

	if err := c.Bind(&input); err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	grade := &models.Grade{
		StudentID:    input.StudentID,
		AssignmentID: input.AssignmentID,
		Mark:         input.Mark,
	}

	id, err := h.service.NewGrade(c.Request().Context(), grade)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, map[string]int{"id": id})
}

// GetRankings retrieves student rankings
// @Summary Get rankings
// @Description Get student rankings by GPA, filtered by group or subject
// @Tags Grades
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param group_id query int false "Filter by Group ID"
// @Param subject_name query string false "Filter by Subject Name"
// @Success 200 {array} models.StudentGPA
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /rankings [get]
func (h *Handler) GetRankings(c echo.Context) error {
	var params struct {
		GroupID     *int    `query:"group_id"`
		SubjectName *string `query:"subject_name"`
	}

	if err := c.Bind(&params); err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	rankings, err := h.service.GetRankings(c.Request().Context(), params.GroupID, params.SubjectName)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, rankings)
}
