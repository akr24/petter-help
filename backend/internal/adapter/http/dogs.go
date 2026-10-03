package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/akr24/petter-help/backend/internal/usecase"
)

type dogsHandler struct {
	dogs *usecase.Dogs
}

func (h dogsHandler) list(c echo.Context) error {
	dogs, err := h.dogs.List(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dogs)
}
