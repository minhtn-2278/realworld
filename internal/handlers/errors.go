package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type ValidationFieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationError struct {
	ErrorCode    int                    `json:"errorCode"`
	ErrorMessage string                 `json:"errorMessage"`
	Errors       []ValidationFieldError `json:"errors"`
}

func (e ValidationError) Error() string {
	return e.ErrorMessage
}

func (e ValidationError) StatusCode() int {
	return http.StatusBadRequest
}

func (e ValidationError) MarshalJSON() ([]byte, error) {
	type validationError ValidationError
	return json.Marshal(validationError(e))
}

type APIErrorResponse struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func (e APIErrorResponse) Error() string {
	return e.ErrorMessage
}

func (e APIErrorResponse) StatusCode() int {
	return e.ErrorCode
}

func (e APIErrorResponse) MarshalJSON() ([]byte, error) {
	type jsonAPIErrorResponse APIErrorResponse
	return json.Marshal(jsonAPIErrorResponse(e))
}

func apiError(statusCode int, message string) APIErrorResponse {
	return APIErrorResponse{ErrorCode: statusCode, ErrorMessage: message}
}

func serviceError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apiError(http.StatusNotFound, "resource not found")
	}

	return apiError(http.StatusInternalServerError, "internal server error")
}

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

	if writeErr := c.JSON(statusCode, apiError(statusCode, message)); writeErr != nil {
		c.Logger().Error("failed to write API error response", "error", writeErr)
	}
}
