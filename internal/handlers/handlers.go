package handlers

import (
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

func JSON(c echo.Context, status int, err error) error {
	response := models.ErrorResponse{
		Error: err.Error(),
	}
	return c.JSON(status, response)
}
