package handlers

import "github.com/labstack/echo/v5"

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) Create(c *echo.Context) error {
	return nil
}
