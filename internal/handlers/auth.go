package handlers

import (
	"errors"
	"net/http"
	"net/mail"

	"github.com/ansarctica/domashka4/internal/models"
	"github.com/labstack/echo/v4"
)

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      models.User  true  "User credentials"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  models.ErrorResponse
// @Failure      500    {object}  models.ErrorResponse
// @Router       /api/auth/register [post]
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

// Login godoc
// @Summary      Login
// @Description  Login and get a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      models.User  true  "User credentials"
// @Success      200    {object}  map[string]string
// @Failure      400    {object}  models.ErrorResponse
// @Failure      401    {object}  models.ErrorResponse
// @Router       /api/auth/login [post]
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

// GetMe godoc
// @Summary      Get current user info
// @Description  Get ID and Email of the logged-in user
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /api/users/me [get]
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
