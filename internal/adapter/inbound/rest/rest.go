package rest

import (
	"github.com/labstack/echo/v4"

	"skeleton-service/internal/adapter/inbound/registry"
	authmw "skeleton-service/internal/adapter/inbound/rest/middleware"
	hv1 "skeleton-service/internal/adapter/inbound/rest/router/v1"
)

func Apply(e *echo.Echo, basicAuth *echo.Group, serviceRegistry registry.ServiceRegistry) {
	hv1.HealthCheckRouter(e)
	hv1.AuthRouter(basicAuth, &serviceRegistry)

	protected := basicAuth.Group("", authmw.JWTAuth(serviceRegistry.GetJWTSecret()))
	hv1.ProfileRouter(protected, &serviceRegistry)
}
