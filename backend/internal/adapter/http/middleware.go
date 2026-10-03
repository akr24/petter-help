package http

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/akr24/petter-help/backend/internal/domain"
	"github.com/akr24/petter-help/backend/internal/token"
)

const claimsKey = "auth.claims"

// requireAuth parses the Bearer token and stores its claims on the context.
// Handlers read them with currentClaims.
func requireAuth(tokens *token.Issuer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get(echo.HeaderAuthorization)
			raw, ok := strings.CutPrefix(header, "Bearer ")
			if !ok || raw == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing bearer token")
			}
			claims, err := tokens.Parse(raw)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}
			c.Set(claimsKey, claims)
			return next(c)
		}
	}
}

// requireRole runs after requireAuth and rejects callers of any other role.
func requireRole(role domain.Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if currentClaims(c).Role != role {
				return echo.NewHTTPError(http.StatusForbidden, string(role)+" account required")
			}
			return next(c)
		}
	}
}

func currentClaims(c echo.Context) token.Claims {
	claims, _ := c.Get(claimsKey).(token.Claims)
	return claims
}
