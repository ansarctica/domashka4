package handlers

import (
	"net/http"
	"strconv"

	"github.com/ansarctica/domashka4/internal/models"
	"github.com/labstack/echo/v4"
)

// GetSchedules retrieves class schedules
// @Summary Get schedules
// @Description Get schedules, optionally filtered by group ID
// @Tags Schedules
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param group_id query int false "Filter by Group ID"
// @Success 200 {array} map[string]interface{}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /schedules [get]
func (h *Handler) GetSchedules(c echo.Context) error {
	var params struct {
		GroupID *int `query:"group_id"`
	}

	if err := c.Bind(&params); err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	schedules, err := h.service.GetSchedules(c.Request().Context(), params.GroupID)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, formatSchedules(schedules))
}

// CreateSchedule adds a new schedule entry
// @Summary Create schedule
// @Description Add a new class schedule entry
// @Tags Schedules
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body models.Schedule true "Schedule Data"
// @Success 201 {object} map[string]int "Returns created ID"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /schedules [post]
func (h *Handler) CreateSchedule(c echo.Context) error {
	var schedule models.Schedule
	if err := c.Bind(&schedule); err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	schedule.ID = 0

	id, err := h.service.CreateSchedule(c.Request().Context(), &schedule)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, map[string]int{"id": id})
}

// UpdateSchedule modifies a schedule entry
// @Summary Update schedule
// @Description Update an existing schedule entry
// @Tags Schedules
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param input body models.Schedule true "Updated Schedule Data"
// @Success 200 {object} map[string]string "Returns status"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /schedules/{id} [patch]
func (h *Handler) UpdateSchedule(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	var schedule models.Schedule
	if err := c.Bind(&schedule); err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}
	schedule.ID = id

	if err := h.service.UpdateSchedule(c.Request().Context(), &schedule); err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "updated"})
}

// DeleteSchedule removes a schedule entry
// @Summary Delete schedule
// @Description Remove a schedule entry from the database
// @Tags Schedules
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} map[string]string "Returns status"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /schedules/{id} [delete]
func (h *Handler) DeleteSchedule(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	if err := h.service.DeleteSchedule(c.Request().Context(), id); err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
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
