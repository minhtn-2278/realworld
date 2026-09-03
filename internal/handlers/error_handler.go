package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"

	"realworldapp/internal/utils"
)

func HTTPErrorHandler(c *echo.Context, err error) {
	if response, unwrapErr := echo.UnwrapResponse(c.Response()); unwrapErr == nil && response.Committed {
		return
	}

	statusCode := http.StatusInternalServerError
	var statusCoder echo.HTTPStatusCoder
	if errors.As(err, &statusCoder) && statusCoder.StatusCode() != 0 {
		statusCode = statusCoder.StatusCode()
	}

	var marshaler json.Marshaler
	if errors.As(err, &marshaler) {
		if writeErr := c.JSON(statusCode, marshaler); writeErr != nil {
			c.Logger().Error("failed to write API error response", "error", writeErr)
		}
		return
	}

	message := http.StatusText(statusCode)
	var httpError *echo.HTTPError
	if errors.As(err, &httpError) && httpError.Message != "" {
		message = httpError.Message
	}

	if writeErr := c.JSON(statusCode, utils.APIError(statusCode, message)); writeErr != nil {
		c.Logger().Error("failed to write API error response", "error", writeErr)
	}
}
