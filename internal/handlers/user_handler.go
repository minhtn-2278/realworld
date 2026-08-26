package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"

	"realworldapp/internal/dto"
	appmiddleware "realworldapp/internal/middleware"
	"realworldapp/internal/services"
	"realworldapp/internal/utils"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Create(c *echo.Context) error {
	var request dto.RegisterUserRequest
	if err := utils.BindAndValidate(c, &request); err != nil {
		return err
	}

	user, err := h.userService.Register(c.Request().Context(), services.RegisterUserInput{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		return utils.ServiceError(err)
	}

	return c.JSON(http.StatusCreated, dto.NewRegisterUserResponse(user))
}

func (h *UserHandler) Login(c *echo.Context) error {
	var request dto.LoginUserRequest
	if err := utils.BindAndValidate(c, &request); err != nil {
		return err
	}
	result, err := h.userService.Login(c.Request().Context(), services.LoginUserInput{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		if errors.Is(err, utils.ErrInvalidCredentials) {
			return utils.APIError(http.StatusUnauthorized, err.Error())
		}
		return utils.ServiceError(err)
	}

	return c.JSON(http.StatusOK, dto.NewLoginUserResponse(result.AccessToken, result.RefreshToken))
}

func (h *UserHandler) Profile(c *echo.Context) error {
	user, err := h.userService.GetProfile(c.Request().Context(), c.Param("username"))
	if err != nil {
		return utils.ServiceError(err)
	}

	return c.JSON(http.StatusOK, dto.NewUserResponse(user))
}

func (h *UserHandler) Current(c *echo.Context) error {
	username, err := appmiddleware.CurrentUsername(c)
	if err != nil {
		return err
	}

	user, err := h.userService.GetProfile(c.Request().Context(), username)
	if err != nil {
		return utils.ServiceError(err)
	}

	return c.JSON(http.StatusOK, dto.NewUserResponse(user))
}
