package handlers

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/ansarctica/domashka4/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (h *Handler) UserIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		if header == "" {
			return JSON(c, http.StatusUnauthorized, errors.New("no auth header"))
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			return JSON(c, http.StatusUnauthorized, errors.New("wrong auth header format"))
		}

		tokenString := headerParts[1]

		token, err := jwt.ParseWithClaims(tokenString, &service.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "wrong signing method")
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil {
			return JSON(c, http.StatusUnauthorized, err)
		}

		claims, ok := token.Claims.(*service.TokenClaims)
		if ok && token.Valid {
			c.Set("userId", claims.UserID)
			return next(c)
		}

		return JSON(c, http.StatusUnauthorized, errors.New("token is not valid"))
	}
}
