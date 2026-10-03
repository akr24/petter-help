package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/akr24/petter-help/backend/internal/domain"
	"github.com/akr24/petter-help/backend/internal/usecase"
)

type seekerProfileHandler struct {
	profiles *usecase.SeekerProfiles
}

func (h seekerProfileHandler) get(c echo.Context) error {
	p, err := h.profiles.Get(c.Request().Context(), currentClaims(c).UserID)
	if errors.Is(err, domain.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "no profile yet")
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, p)
}

func (h seekerProfileHandler) put(c echo.Context) error {
	var in domain.SeekerProfile
	if err := c.Bind(&in); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "malformed JSON")
	}

	p, err := h.profiles.Save(c.Request().Context(), currentClaims(c).UserID, in)
	if errors.Is(err, usecase.ErrInvalidInput) {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, p)
}
