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

// GetStudent godoc
// @Summary      Get student details
// @Tags         students
// @Produce      json
// @Param        id   path      int  true  "Student ID"
// @Success      200  {object}  models.Student
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /student/{id} [get]
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

// GetAllSchedules godoc
// @Summary      Get all schedules
// @Tags         schedule
// @Produce      json
// @Success      200  {array}   map[string]interface{}
// @Failure      500  {object}  models.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /all_class_schedule [get]
func (h *Handler) GetAllSchedules(c echo.Context) error {
	schedules, err := h.service.GetAllSchedules(c.Request().Context())
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, formatSchedules(schedules))
}

// GetGroupSchedule godoc
// @Summary      Get schedule by group
// @Tags         schedule
// @Produce      json
// @Param        id   path      int  true  "Group ID"
// @Success      200  {array}   map[string]interface{}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /schedule/group/{id} [get]
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

// NewAttendance godoc
// @Summary      Record attendance
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Param        input  body      attendanceInput  true  "Attendance data"
// @Success      200    {object}  map[string]int
// @Failure      400    {object}  models.ErrorResponse
// @Failure      500    {object}  models.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /attendance/subject [post]
func (h *Handler) NewAttendance(c echo.Context) error {
	var input attendanceInput

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

// GetAttendanceBySubject godoc
// @Summary      Get attendance by subject
// @Tags         attendance
// @Produce      json
// @Param        id   path      int  true  "Subject ID"
// @Success      200  {array}   models.Attendance
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /attendanceBySubjectId/{id} [get]
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

// GetAttendanceByStudent godoc
// @Summary      Get attendance by student
// @Tags         attendance
// @Produce      json
// @Param        id   path      int  true  "Student ID"
// @Success      200  {array}   models.Attendance
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /attendanceByStudentId/{id} [get]
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

// NewAssignment godoc
// @Summary      Create assignment
// @Tags         assignments
// @Accept       json
// @Produce      json
// @Param        input  body      assignmentInput  true  "Assignment data"
// @Success      200    {object}  map[string]int
// @Failure      400    {object}  models.ErrorResponse
// @Failure      500    {object}  models.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /assignments [post]
func (h *Handler) NewAssignment(c echo.Context) error {
	var input assignmentInput

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

// NewGrade godoc
// @Summary      Add a grade
// @Tags         grades
// @Accept       json
// @Produce      json
// @Param        input  body      gradeInput  true  "Grade data"
// @Success      200    {object}  map[string]int
// @Failure      400    {object}  models.ErrorResponse
// @Failure      500    {object}  models.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /grades [post]
func (h *Handler) NewGrade(c echo.Context) error {
	var input gradeInput

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

// GetGPA godoc
// @Summary      Get Student GPA
// @Tags         gpa
// @Produce      json
// @Param        id   path      int  true  "Student ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /students/{id}/gpa [get]
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

// GetSubjectGPA godoc
// @Summary      Get Student GPA for a specific subject
// @Tags         gpa
// @Produce      json
// @Param        studentId  path      int  true  "Student ID"
// @Param        subjectId  path      int  true  "Subject ID"
// @Success      200        {object}  map[string]interface{}
// @Failure      400        {object}  models.ErrorResponse
// @Failure      500        {object}  models.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /students/{studentId}/subjects/{subjectId}/gpa [get]
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

// GetGPARankingByGroup godoc
// @Summary      Get GPA Ranking by Group
// @Tags         ranking
// @Produce      json
// @Param        id   path      int  true  "Group ID"
// @Success      200  {array}   map[string]interface{}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /groups/{id}/ranking [get]
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

// GetGPARankingBySubject godoc
// @Summary      Get GPA Ranking by Subject
// @Tags         ranking
// @Produce      json
// @Param        id   path      int  true  "Subject ID"
// @Success      200  {array}   map[string]interface{}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /subjects/{id}/ranking [get]
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

// GetSubjectGPARankingByGroup godoc
// @Summary      Get Subject GPA Ranking by Group
// @Tags         ranking
// @Produce      json
// @Param        groupId    path      int  true  "Group ID"
// @Param        subjectId  path      int  true  "Subject ID"
// @Success      200        {array}   map[string]interface{}
// @Failure      400        {object}  models.ErrorResponse
// @Failure      500        {object}  models.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /groups/{groupId}/subjects/{subjectId}/ranking [get]
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

type attendanceInput struct {
	SubjectID int    `json:"subject_id"`
	VisitDay  string `json:"visit_day"`
	Visited   bool   `json:"visited"`
	StudentID int    `json:"Student_id"`
}

type assignmentInput struct {
	Name      string `json:"name"`
	SubjectID int    `json:"subject_id"`
	Weight    int    `json:"weight"`
	Date      string `json:"date"`
}

type gradeInput struct {
	StudentID    int `json:"student_id"`
	AssignmentID int `json:"assignment_id"`
	Mark         int `json:"mark"`
}
