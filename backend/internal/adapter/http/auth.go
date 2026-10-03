package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/akr24/petter-help/backend/internal/domain"
	"github.com/akr24/petter-help/backend/internal/usecase"
)

type authHandler struct {
	auth *usecase.Auth
}

type registerRequest struct {
	Email    string      `json:"email"`
	Password string      `json:"password"`
	Role     domain.Role `json:"role"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h authHandler) register(c echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "malformed JSON")
	}

	u, err := h.auth.Register(c.Request().Context(), req.Email, req.Password, req.Role)
	switch {
	case errors.Is(err, usecase.ErrInvalidInput):
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrEmailTaken):
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	case err != nil:
		return err
	}
	return c.JSON(http.StatusCreated, u)
}

func (h authHandler) login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "malformed JSON")
	}

	u, err := h.auth.Login(c.Request().Context(), req.Email, req.Password)
	if errors.Is(err, domain.ErrInvalidCredentials) {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	if err != nil {
		return err
	}
	// No session or token yet: that's a separate step.
	return c.JSON(http.StatusOK, u)
}
