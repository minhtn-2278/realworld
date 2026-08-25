package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"realworldapp/internal/dto"
	"realworldapp/internal/services"
)

type TagHandler struct {
	tagService services.TagService
}

func NewTagHandler(tagService services.TagService) *TagHandler {
	return &TagHandler{tagService: tagService}
}

func (h *TagHandler) List(c *echo.Context) error {
	tags, err := h.tagService.List(c.Request().Context())
	if err != nil {
		return apiError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, dto.NewTagListResponse(tags))
}
