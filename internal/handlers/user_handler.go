package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"realworldapp/internal/dto"
	"realworldapp/internal/services"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Create(c *echo.Context) error {
	return nil
}

func (h *UserHandler) Profile(c *echo.Context) error {
	user, err := h.userService.GetProfile(c.Request().Context(), c.Param("username"))
	if err != nil {
		return serviceError(err)
	}

	return c.JSON(http.StatusOK, dto.NewUserResponse(user))
}
