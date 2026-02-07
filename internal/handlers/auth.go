package handlers

import (
	"errors"
	"net/http"
	"net/mail"

	"github.com/ansarctica/domashka4/internal/models"
	"github.com/labstack/echo/v4"
)

// Register creates a new user
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body models.User true "User Registration Details"
// @Success 200 {object} map[string]int "Returns user ID"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/register [post]
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

// Login authenticates a user
// @Summary Login
// @Description Authenticate using email and password to receive a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body models.User true "User Credentials"
// @Success 200 {object} map[string]string "Returns JWT token"
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/login [post]
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

// GetMe gets the current user's profile
// @Summary Get Current User
// @Description Get the profile details of the currently logged-in user
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/me [get]
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
