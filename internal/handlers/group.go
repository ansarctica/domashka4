package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetGroups retrieves all groups
// @Summary Get groups
// @Description Get a list of all student groups
// @Tags Groups
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} models.Group
// @Failure 500 {object} models.ErrorResponse
// @Router /groups [get]
func (h *Handler) GetGroups(c echo.Context) error {
	groups, err := h.service.GetAllGroups(c.Request().Context())
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, groups)
}
