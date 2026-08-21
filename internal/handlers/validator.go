package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

var requestValidator = validator.New(validator.WithRequiredStructEnabled())

func bindAndValidate(c *echo.Context, request any) error {
	if err := c.Bind(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if err := requestValidator.Struct(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "request validation failed")
	}

	return nil
}
