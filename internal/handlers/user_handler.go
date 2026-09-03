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

// Create registers a new user.
// @Summary Register a user
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.RegisterUserRequest true "Registration payload"
// @Success 201 {object} dto.RegisterUserResponse
// @Failure 400 {object} utils.ValidationError
// @Failure 409 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Router /users [post]
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

// Login authenticates a user and returns access tokens.
// @Summary Log in
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.LoginUserRequest true "Credentials"
// @Success 200 {object} dto.LoginUserResponse
// @Failure 400 {object} utils.ValidationError
// @Failure 401 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Router /users/login [post]
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

// Profile returns a user profile for the authenticated viewer.
// @Summary Get a user profile
// @Tags profiles
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Security BearerAuth
// @Router /profiles/{username} [get]
func (h *UserHandler) Profile(c *echo.Context) error {
	viewerID, err := appmiddleware.CurrentUserID(c)
	if err != nil {
		return err
	}

	user, err := h.userService.GetProfileForUser(c.Request().Context(), c.Param("username"), viewerID)
	if err != nil {
		return utils.ServiceError(err)
	}

	return c.JSON(http.StatusOK, dto.NewUserResponse(user))
}

// Follow follows or unfollows a user for the authenticated user.
// @Summary Follow or unfollow a user
// @Tags profiles
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Param request body dto.FollowUserRequest true "Follow payload"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} utils.ValidationError
// @Failure 401 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Security BearerAuth
// @Router /profiles/{username}/follow [post]
func (h *UserHandler) Follow(c *echo.Context) error {
	var request dto.FollowUserRequest
	if err := utils.BindAndValidate(c, &request); err != nil {
		return err
	}
	followerID, err := appmiddleware.CurrentUserID(c)
	if err != nil {
		return err
	}

	if err := h.userService.Follow(c.Request().Context(), c.Param("username"), followerID, *request.Follow); err != nil {
		return utils.ServiceError(err)
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{Message: "ok"})
}

// Current returns the authenticated user's profile.
// @Summary Get the current user
// @Tags users
// @Produce json
// @Success 200 {object} dto.CurrentUserResponse
// @Failure 401 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Security BearerAuth
// @Router /user [get]
func (h *UserHandler) Current(c *echo.Context) error {
	username, err := appmiddleware.CurrentUsername(c)
	if err != nil {
		return err
	}

	user, err := h.userService.GetProfile(c.Request().Context(), username)
	if err != nil {
		return utils.ServiceError(err)
	}

	return c.JSON(http.StatusOK, dto.NewCurrentUserResponse(user))
}
