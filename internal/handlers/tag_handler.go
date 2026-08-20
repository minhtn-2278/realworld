package handlers

import "github.com/labstack/echo/v5"

type TagHandler struct{}

func NewTagHandler() *TagHandler {
	return &TagHandler{}
}

func (h *TagHandler) List(c *echo.Context) error {
	return nil
}
