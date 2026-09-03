package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"realworldapp/internal/dto"
	"realworldapp/internal/services"
	"realworldapp/internal/utils"
)

type TagHandler struct {
	tagService services.TagService
}

func NewTagHandler(tagService services.TagService) *TagHandler {
	return &TagHandler{tagService: tagService}
}

// List lists all tags.
// @Summary List tags
// @Tags tags
// @Produce json
// @Success 200 {object} dto.TagListResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Router /tags [get]
func (h *TagHandler) List(c *echo.Context) error {
	tags, err := h.tagService.List(c.Request().Context())
	if err != nil {
		return utils.APIError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, dto.NewTagListResponse(tags))
}
