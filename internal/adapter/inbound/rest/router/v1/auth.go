package v1

import (
	"github.com/labstack/echo/v4"

	"skeleton-service/internal/adapter/inbound/registry"
	h "skeleton-service/internal/adapter/inbound/rest/router/v1/handler"
)

func AuthRouter(eg *echo.Group, reg *registry.ServiceRegistry) {
	handler := h.NewAuthHandler(reg)

	eg.POST("/api/v1/auth/register", handler.Register)
	eg.POST("/api/v1/auth/login", handler.Login)
}

func ProfileRouter(eg *echo.Group, reg *registry.ServiceRegistry) {
	handler := h.NewAuthHandler(reg)
	eg.GET("/api/v1/profile", handler.Profile)
}
