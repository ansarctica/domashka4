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
	id, _ := strconv.Atoi(c.Param("id"))

	student := h.service.GetStudent(c.Request().Context(), id)

	return c.JSON(http.StatusOK, student)
}

func (h *Handler) GetAllSchedules(c echo.Context) error {
	schedules := formatSchedules(h.service.GetAllSchedules(c.Request().Context()))

	return c.JSON(http.StatusOK, schedules)
}

func (h *Handler) GetGroupSchedule(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	schedule := formatSchedules(h.service.GetGroupSchedule(c.Request().Context(), id))

	return c.JSON(http.StatusOK, schedule)
}

func (h *Handler) NewAttendance(c echo.Context) error {

	var input struct {
		SubjectID int    `json:"subject_id"`
		VisitDay  string `json:"visit_day"`
		Visited   bool   `json:"visited"`
		StudentID int    `json:"Student_id"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	parsedDate, _ := time.Parse("02.01.2006", input.VisitDay)

	attendance := &models.Attendance{
		SubjectID: input.SubjectID,
		VisitDay:  parsedDate,
		Visited:   input.Visited,
		StudentID: input.StudentID,
	}

	id := h.service.NewAttendance(c.Request().Context(), attendance)

	return c.JSON(http.StatusOK, map[string]int{"id": id})
}

func (h *Handler) GetAttendanceBySubject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	attendanceList := h.service.AttendanceBySubject(c.Request().Context(), id)

	return c.JSON(http.StatusOK, attendanceList)
}

func (h *Handler) GetAttendanceByStudent(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	attendanceList := h.service.AttendanceByStudent(c.Request().Context(), id)

	return c.JSON(http.StatusOK, attendanceList)
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
