package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTAuth(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if header == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"code":    http.StatusUnauthorized,
					"message": "missing authorization header",
				})
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"code":    http.StatusUnauthorized,
					"message": "invalid authorization header",
				})
			}

			tokenString := strings.TrimSpace(parts[1])
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"code":    http.StatusUnauthorized,
					"message": "invalid token",
				})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"code":    http.StatusUnauthorized,
					"message": "invalid token claims",
				})
			}

			userID, err := userIDFromClaims(claims)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"code":    http.StatusUnauthorized,
					"message": "invalid token subject",
				})
			}

			c.Set("user_id", userID)
			return next(c)
		}
	}
}

func userIDFromClaims(claims jwt.MapClaims) (int64, error) {
	sub, ok := claims["sub"]
	if !ok {
		return 0, errors.New("missing sub")
	}

	switch v := sub.(type) {
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	case string:
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, err
		}
		return id, nil
	default:
		return 0, errors.New("unsupported sub type")
	}
}
