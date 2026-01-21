package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ansarctica/domashka4/internal/models"
	"github.com/ansarctica/domashka4/internal/service"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

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

func (h *Handler) GetAllSchedules(c echo.Context) error {
	schedules, err := h.service.GetAllSchedules(c.Request().Context())
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, formatSchedules(schedules))
}

func (h *Handler) GetGroupSchedule(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	schedule, err := h.service.GetGroupSchedule(c.Request().Context(), id)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, formatSchedules(schedule))
}

func (h *Handler) NewAttendance(c echo.Context) error {
	var input struct {
		SubjectID int    `json:"subject_id"`
		VisitDay  string `json:"visit_day"`
		Visited   bool   `json:"visited"`
		StudentID int    `json:"Student_id"`
	}

	if err := c.Bind(&input); err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	parsedDate, err := time.Parse("02.01.2006", input.VisitDay)
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	attendance := &models.Attendance{
		SubjectID: input.SubjectID,
		VisitDay:  parsedDate,
		Visited:   input.Visited,
		StudentID: input.StudentID,
	}

	id, err := h.service.NewAttendance(c.Request().Context(), attendance)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]int{"id": id})
}

func (h *Handler) GetAttendanceBySubject(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	attendanceList, err := h.service.AttendanceBySubject(c.Request().Context(), id)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, attendanceList)
}

func (h *Handler) GetAttendanceByStudent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	attendanceList, err := h.service.AttendanceByStudent(c.Request().Context(), id)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, attendanceList)
}

func (h *Handler) NewAssignment(c echo.Context) error {
	var input struct {
		Name      string `json:"name"`
		SubjectID int    `json:"subject_id"`
		Weight    int    `json:"weight"`
		Date      string `json:"date"`
	}

	if err := c.Bind(&input); err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	parsedDate, err := time.Parse("02.01.2006", input.Date)
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	assignment := &models.Assignment{
		Name:      input.Name,
		SubjectID: input.SubjectID,
		Weight:    input.Weight,
		Date:      parsedDate,
	}

	id, err := h.service.NewAssignment(c.Request().Context(), assignment)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]int{"id": id})
}

func (h *Handler) NewGrade(c echo.Context) error {
	var input struct {
		StudentID    int `json:"student_id"`
		AssignmentID int `json:"assignment_id"`
		Mark         int `json:"mark"`
	}

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

	return c.JSON(http.StatusOK, map[string]int{"id": id})
}

func (h *Handler) GetGPA(c echo.Context) error {
	studentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	gpa, err := h.service.GetGPA(c.Request().Context(), studentID)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"student_id": studentID,
		"gpa":        gpa,
	})
}

func (h *Handler) GetSubjectGPA(c echo.Context) error {
	studentID, err := strconv.Atoi(c.Param("studentId"))
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	subjectID, err := strconv.Atoi(c.Param("subjectId"))
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	gpa, err := h.service.GetSubjectGPA(c.Request().Context(), studentID, subjectID)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"student_id": studentID,
		"subject_id": subjectID,
		"gpa":        gpa,
	})
}

func (h *Handler) GetGPARankingByGroup(c echo.Context) error {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	ranking, err := h.service.RankingByGroup(c.Request().Context(), groupID)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, ranking)
}

func (h *Handler) GetGPARankingBySubject(c echo.Context) error {
	subjectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}
	ranking, err := h.service.RankingBySubject(c.Request().Context(), subjectID)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, ranking)
}

func (h *Handler) GetSubjectGPARankingByGroup(c echo.Context) error {
	groupID, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	subjectID, err := strconv.Atoi(c.Param("subjectId"))
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	ranking, err := h.service.RankingByGroupSubject(c.Request().Context(), groupID, subjectID)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, ranking)
}

func formatSchedules(schedules []models.Schedule) []map[string]interface{} {
	result := make([]map[string]interface{}, len(schedules))
	for i, s := range schedules {
		result[i] = map[string]interface{}{
			"id":         s.ID,
			"group_id":   s.GroupID,
			"subject":    s.Subject,
			"start_time": s.StartTime.Format("15:04"),
			"end_time":   s.EndTime.Format("15:04"),
		}
	}
	return result
}

func JSON(c echo.Context, status int, err error) error {
	response := models.ErrorResponse{
		Error: err.Error(),
	}
	return c.JSON(status, response)
}
