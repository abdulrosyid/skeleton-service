package util

import "github.com/labstack/echo/v4"

type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SetResponse(c echo.Context, code int, message string, data interface{}) error {
	res := BaseResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
	return c.JSON(code, res)
}

func SetResponseError(c echo.Context, code int, message string) error {
	res := BaseResponse{
		Code:    code,
		Message: message,
	}
	return c.JSON(code, res)
}
