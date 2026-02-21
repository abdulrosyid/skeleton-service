package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"skeleton-service/internal/adapter/inbound/registry"
	cservice "skeleton-service/internal/core/service"
	iservice "skeleton-service/internal/port/inbound/service"
	"skeleton-service/shared/util"
)

type AuthHandler struct {
	reg *registry.ServiceRegistry
}

func NewAuthHandler(reg *registry.ServiceRegistry) *AuthHandler {
	return &AuthHandler{reg: reg}
}

func (h *AuthHandler) Register(c echo.Context) error {
	ctx := c.Request().Context()
	svc := h.reg.GetAuthService()

	var req iservice.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return util.SetResponseError(c, http.StatusBadRequest, "invalid request body")
	}

	data, err := svc.Register(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, cservice.ErrInvalidInput):
			return util.SetResponseError(c, http.StatusBadRequest, "full_name/email wajib diisi, password minimal 6 karakter")
		case errors.Is(err, cservice.ErrEmailExists):
			return util.SetResponseError(c, http.StatusConflict, "email sudah terdaftar")
		default:
			return util.SetResponseError(c, http.StatusInternalServerError, "internal server error")
		}
	}

	return util.SetResponse(c, http.StatusCreated, "registrasi berhasil", data)
}

func (h *AuthHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	svc := h.reg.GetAuthService()

	var req iservice.LoginRequest
	if err := c.Bind(&req); err != nil {
		return util.SetResponseError(c, http.StatusBadRequest, "invalid request body")
	}

	data, err := svc.Login(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, cservice.ErrInvalidInput):
			return util.SetResponseError(c, http.StatusBadRequest, "email dan password wajib diisi")
		case errors.Is(err, cservice.ErrInvalidCreds):
			return util.SetResponseError(c, http.StatusUnauthorized, "email atau password salah")
		default:
			return util.SetResponseError(c, http.StatusInternalServerError, "internal server error")
		}
	}

	return util.SetResponse(c, http.StatusOK, "login berhasil", data)
}

func (h *AuthHandler) Profile(c echo.Context) error {
	ctx := c.Request().Context()
	svc := h.reg.GetAuthService()

	userID, err := contextUserID(c)
	if err != nil {
		return util.SetResponseError(c, http.StatusUnauthorized, "unauthorized")
	}

	data, err := svc.Profile(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, cservice.ErrUserNotFound):
			return util.SetResponseError(c, http.StatusNotFound, "user tidak ditemukan")
		default:
			return util.SetResponseError(c, http.StatusInternalServerError, "internal server error")
		}
	}

	return util.SetResponse(c, http.StatusOK, "success", data)
}

func contextUserID(c echo.Context) (int64, error) {
	raw := c.Get("user_id")
	if raw == nil {
		return 0, errors.New("missing user id")
	}

	switch v := raw.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, errors.New("invalid user id type")
	}
}
