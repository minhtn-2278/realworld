package handlers

import (
	"testing"

	"github.com/labstack/echo/v5"
)

func TestRegisterRoutes(t *testing.T) {
	e := echo.New()
	authMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return next
	}

	RegisterRoutes(e, articleServiceStub{}, tagServiceStub{}, userServiceStub{}, authMiddleware)
	if len(e.Router().Routes()) == 0 {
		t.Fatal("expected routes to be registered")
	}
}
