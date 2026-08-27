package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"

	"realworldapp/internal/utils"
)

func JWTErrorHandler(_ *echo.Context, err error) error {
	switch {
	case errors.Is(err, echojwt.ErrJWTMissing):
		return utils.APIError(http.StatusUnauthorized, "access token is required")
	case errors.Is(err, echojwt.ErrJWTInvalid):
		return utils.APIError(http.StatusUnauthorized, "invalid or expired access token")
	default:
		return utils.APIError(http.StatusUnauthorized, "invalid access token")
	}
}

func CurrentUsername(c *echo.Context) (string, error) {
	claims, err := currentAccessTokenClaims(c)
	if err != nil {
		return "", err
	}

	username, ok := claims["username"].(string)
	if !ok || username == "" {
		return "", utils.APIError(http.StatusUnauthorized, "invalid access token")
	}

	return username, nil
}

func CurrentUserID(c *echo.Context) (uint, error) {
	claims, err := currentAccessTokenClaims(c)
	if err != nil {
		return 0, err
	}

	subject, ok := claims["sub"].(string)
	if !ok || subject == "" {
		return 0, utils.APIError(http.StatusUnauthorized, "invalid access token")
	}

	userID, err := strconv.ParseUint(subject, 10, 0)
	if err != nil || userID == 0 {
		return 0, utils.APIError(http.StatusUnauthorized, "invalid access token")
	}

	return uint(userID), nil
}

type OwnerResolver func(ctx context.Context, c *echo.Context) (uint, error)

func RequireOwner(authMiddleware echo.MiddlewareFunc, resolveOwner OwnerResolver) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		ownerCheck := func(c *echo.Context) error {
			userID, err := CurrentUserID(c)
			if err != nil {
				return err
			}

			ownerID, err := resolveOwner(c.Request().Context(), c)
			if err != nil {
				var statusCoder echo.HTTPStatusCoder
				if errors.As(err, &statusCoder) {
					return err
				}
				return utils.ServiceError(err)
			}

			if ownerID != userID {
				return utils.APIError(http.StatusForbidden, "only the resource owner can modify this resource")
			}

			return next(c)
		}

		return authMiddleware(ownerCheck)
	}
}

func currentAccessTokenClaims(c *echo.Context) (jwt.MapClaims, error) {
	token, err := echo.ContextGet[*jwt.Token](c, "user")
	if err != nil || token == nil {
		return nil, utils.APIError(http.StatusUnauthorized, "invalid access token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, utils.APIError(http.StatusUnauthorized, "invalid access token")
	}

	tokenType, ok := claims["token_type"].(string)
	if !ok || tokenType != "access" {
		return nil, utils.APIError(http.StatusUnauthorized, "access token required")
	}

	return claims, nil
}
