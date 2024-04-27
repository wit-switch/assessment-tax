package middleware

import (
	"context"
	"crypto/subtle"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type key string

const (
	keyBasicAuth key = "basicAuth"
)

func (m *Middleware) Auth() echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte(m.auth.Username)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(m.auth.Password)) == 1 {
			ctx := context.WithValue(c.Request().Context(), keyBasicAuth, username)
			c.SetRequest(c.Request().WithContext(ctx))

			return true, nil
		}

		return false, nil
	})
}
