package v1

import "github.com/labstack/echo/v4"

func HealthCheckRouter(e *echo.Echo) {
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})
}
