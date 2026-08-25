package handlers

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

var requestValidator = validator.New(validator.WithRequiredStructEnabled())

func bindAndValidate(c *echo.Context, request any) error {
	if err := c.Bind(request); err != nil {
		return apiError(http.StatusBadRequest, "invalid request body")
	}

	if err := requestValidator.Struct(request); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return newValidationError(request, validationErrors)
		}

		return apiError(http.StatusBadRequest, "request validation failed")
	}

	return nil
}

func newValidationError(request any, validationErrors validator.ValidationErrors) ValidationError {
	errors := make([]ValidationFieldError, 0, len(validationErrors))
	for _, validationError := range validationErrors {
		field := jsonFieldName(request, validationError)
		errors = append(errors, ValidationFieldError{
			Field:   field,
			Message: validationMessage(field, validationError),
		})
	}

	errorMessage := "request validation failed"
	if len(errors) > 0 {
		errorMessage = errors[0].Message
	}

	return ValidationError{
		ErrorCode:    http.StatusBadRequest,
		ErrorMessage: errorMessage,
		Errors:       errors,
	}
}

func jsonFieldName(request any, validationError validator.FieldError) string {
	fieldName := validationError.Field()
	baseFieldName := validationError.StructField()
	suffix := ""
	if index := strings.Index(fieldName, "["); index >= 0 {
		baseFieldName = fieldName[:index]
		suffix = fieldName[index:]
	}

	typeOfRequest := reflect.TypeOf(request)
	for typeOfRequest.Kind() == reflect.Pointer {
		typeOfRequest = typeOfRequest.Elem()
	}
	if typeOfRequest.Kind() != reflect.Struct {
		return fieldName
	}

	field, ok := typeOfRequest.FieldByName(baseFieldName)
	if !ok {
		return fieldName
	}

	jsonName := strings.Split(field.Tag.Get("json"), ",")[0]
	if jsonName == "" || jsonName == "-" {
		jsonName = baseFieldName
	}

	return jsonName + suffix
}

func validationMessage(field string, validationError validator.FieldError) string {
	switch validationError.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email"
	case "max":
		return field + " must be at most " + validationError.Param() + " characters"
	case "min":
		return field + " must be at least " + validationError.Param() + " characters"
	case "len":
		return field + " must be exactly " + validationError.Param() + " characters"
	case "oneof":
		return field + " must be one of: " + validationError.Param()
	default:
		return field + " failed validation rule " + validationError.Tag()
	}
}
