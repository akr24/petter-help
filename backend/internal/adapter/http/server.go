// Package http is the Echo-based HTTP adapter. It translates requests into
// use-case calls and use-case errors into status codes; no business logic
// lives here.
package http

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/akr24/petter-help/backend/internal/domain"
	"github.com/akr24/petter-help/backend/internal/token"
	"github.com/akr24/petter-help/backend/internal/usecase"
)

// Pinger is the slice of the database pool the health check needs.
type Pinger interface {
	Ping(ctx context.Context) error
}

// Deps are the collaborators the router wires handlers to.
type Deps struct {
	DB         Pinger
	Dogs       *usecase.Dogs
	Auth       *usecase.Auth
	Profiles   *usecase.SeekerProfiles
	Tokens     *token.Issuer
	CORSOrigin string
}

// New builds the Echo instance with middleware and routes registered.
func New(d Deps) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{d.CORSOrigin},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	}))

	e.GET("/healthz", healthHandler(d.DB))

	api := e.Group("/api")
	api.GET("/dogs", dogsHandler{d.Dogs}.list)

	auth := api.Group("/auth")
	h := authHandler{d.Auth}
	auth.POST("/register", h.register)
	auth.POST("/login", h.login)
	auth.GET("/me", h.me, requireAuth(d.Tokens))

	seeker := api.Group("/seeker", requireAuth(d.Tokens), requireRole(domain.RoleSeeker))
	p := seekerProfileHandler{d.Profiles}
	seeker.GET("/profile", p.get)
	seeker.PUT("/profile", p.put)

	return e
}

func healthHandler(db Pinger) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := db.Ping(c.Request().Context()); err != nil {
			return echo.NewHTTPError(503, "database unreachable")
		}
		return c.JSON(200, map[string]string{"status": "ok"})
	}
}
