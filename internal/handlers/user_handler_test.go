package handlers

import (
	"context"
	"net/http"
	"testing"

	"github.com/labstack/echo/v5"

	"realworldapp/internal/models"
	"realworldapp/internal/services"
)

func TestUserHandlers(t *testing.T) {
	handler := NewUserHandler(userServiceStub{})

	tests := []struct {
		name       string
		method     string
		target     string
		body       string
		params     []string
		auth       bool
		statusCode int
		handler    func(*echo.Context) error
	}{
		{"create", http.MethodPost, "/api/users", `{"username":"alice","email":"alice@example.com","password":"password"}`, nil, false, http.StatusCreated, handler.Create},
		{"login", http.MethodPost, "/api/users/login", `{"username":"alice","password":"password"}`, nil, false, http.StatusOK, handler.Login},
		{"profile", http.MethodGet, "/api/profiles/alice", "", []string{"username", "alice"}, true, http.StatusOK, handler.Profile},
		{"follow", http.MethodPost, "/api/profiles/alice/follow", `{"follow":true}`, []string{"username", "alice"}, true, http.StatusOK, handler.Follow},
		{"current user", http.MethodGet, "/api/user", "", nil, true, http.StatusOK, handler.Current},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, recorder := newHandlerTestContext(test.method, test.target, test.body)
			setPathParams(context, test.params...)
			if test.auth {
				setAuthenticatedUser(context)
			}

			assertHandlerStatus(t, test.handler(context), recorder, test.statusCode)
		})
	}
}

type userServiceStub struct{}

func (userServiceStub) Register(_ context.Context, _ services.RegisterUserInput) (*models.User, error) {
	return userTestData(), nil
}

func (userServiceStub) Login(_ context.Context, _ services.LoginUserInput) (*services.LoginUserResult, error) {
	return &services.LoginUserResult{AccessToken: "access-token", RefreshToken: "refresh-token"}, nil
}

func (userServiceStub) GetByID(_ context.Context, _ uint) (*models.User, error) {
	return userTestData(), nil
}

func (userServiceStub) GetProfile(_ context.Context, _ string) (*models.User, error) {
	return userTestData(), nil
}

func (userServiceStub) GetProfileForUser(_ context.Context, _ string, _ uint) (*models.User, error) {
	user := userTestData()
	user.Following = true
	return user, nil
}

func (userServiceStub) Follow(_ context.Context, _ string, _ uint, _ bool) error { return nil }

func userTestData() *models.User {
	return &models.User{
		ID:       1,
		Username: "alice",
		Email:    "alice@example.com",
		Articles: []models.Article{{ID: 1, Title: "Post", AuthorID: 1}},
	}
}
