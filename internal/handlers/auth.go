package handlers

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

func JWTErrorHandler(_ *echo.Context, err error) error {
	switch {
	case errors.Is(err, echojwt.ErrJWTMissing):
		return apiError(http.StatusUnauthorized, "access token is required")
	case errors.Is(err, echojwt.ErrJWTInvalid):
		return apiError(http.StatusUnauthorized, "invalid or expired access token")
	default:
		return apiError(http.StatusUnauthorized, "invalid access token")
	}
}

func currentUsername(c *echo.Context) (string, error) {
	token, err := echo.ContextGet[*jwt.Token](c, "user")
	if err != nil {
		return "", apiError(http.StatusUnauthorized, "invalid access token")
	}
	if token == nil {
		return "", apiError(http.StatusUnauthorized, "invalid access token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", apiError(http.StatusUnauthorized, "invalid access token")
	}

	tokenType, ok := claims["token_type"].(string)
	if !ok || tokenType != "access" {
		return "", apiError(http.StatusUnauthorized, "access token required")
	}

	username, ok := claims["username"].(string)
	if !ok || username == "" {
		return "", apiError(http.StatusUnauthorized, "invalid access token")
	}

	return username, nil
}
