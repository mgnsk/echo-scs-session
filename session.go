package session

import (
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	Skipper        middleware.Skipper
	SessionManager *scs.SessionManager
}

var DefaultSessionConfig = Config{
	Skipper: middleware.DefaultSkipper,
}

func LoadAndSave(sessionManager *scs.SessionManager) echo.MiddlewareFunc {
	c := DefaultSessionConfig
	c.SessionManager = sessionManager

	return LoadAndSaveWithConfig(c)
}

func LoadAndSaveWithConfig(config Config) echo.MiddlewareFunc {

	if config.Skipper == nil {
		config.Skipper = DefaultSessionConfig.Skipper
	}

	if config.SessionManager == nil {
		panic("Session middleware requires a session manager")
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			c.Response().Header().Add("Vary", "Cookie")

			ctx := c.Request().Context()
			var token string
			cookie, err := c.Cookie(config.SessionManager.Cookie.Name)
			if err == nil {
				token = cookie.Value
			}

			ctx, err = config.SessionManager.Load(ctx, token)
			if err != nil {
				return err
			}

			c.SetRequest(c.Request().WithContext(ctx))

			c.Response().Before(func() {
				switch config.SessionManager.Status(ctx) {

				case scs.Modified:
					token, expiry, err := config.SessionManager.Commit(ctx)
					if err != nil {
						panic(err)
					}

					config.SessionManager.WriteSessionCookie(ctx, c.Response().Writer, token, expiry)

				case scs.Destroyed:
					config.SessionManager.WriteSessionCookie(ctx, c.Response().Writer, "", time.Time{})
				default:
					// session might not exist yet
				}
			})
			return next(c)
		}
	}
}
