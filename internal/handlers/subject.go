package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetSubjects retrieves all subjects
// @Summary Get subjects
// @Description Get a list of all available subjects
// @Tags Subjects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} models.Subject
// @Failure 500 {object} models.ErrorResponse
// @Router /subjects [get]
func (h *Handler) GetSubjects(c echo.Context) error {
	subjects, err := h.service.GetAllSubjects(c.Request().Context())
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, subjects)
}
