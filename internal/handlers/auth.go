package handlers

import (
	"errors"
	"net/http"
	"net/mail"

	"github.com/ansarctica/domashka4/internal/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(c echo.Context) error {
	var input models.User

	if err := c.Bind(&input); err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	if _, err := mail.ParseAddress(input.Email); err != nil {
		return JSON(c, http.StatusBadRequest, errors.New("wrong email format"))
	}

	id, err := h.service.CreateUser(c.Request().Context(), &input)
	if err != nil {
		return JSON(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) Login(c echo.Context) error {
	var input models.User

	if err := c.Bind(&input); err != nil {
		return JSON(c, http.StatusBadRequest, err)
	}

	token, err := h.service.GenerateToken(c.Request().Context(), input.Email, input.Password)
	if err != nil {
		return JSON(c, http.StatusUnauthorized, err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

func (h *Handler) GetMe(c echo.Context) error {
	id, ok := c.Get("userId").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "no user id"})
	}

	user, err := h.service.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return JSON(c, http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
	})
}
