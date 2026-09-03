package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrCannotFollowSelf   = errors.New("cannot follow yourself")
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

func APIError(statusCode int, message string) APIErrorResponse {
	return APIErrorResponse{ErrorCode: statusCode, ErrorMessage: message}
}

func ServiceError(err error) error {
	if errors.Is(err, ErrCannotFollowSelf) {
		return APIError(http.StatusBadRequest, err.Error())
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return APIError(http.StatusNotFound, "resource not found")
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return APIError(http.StatusConflict, "resource already exists")
	}

	return APIError(http.StatusInternalServerError, "internal server error")
}
