package handlers

import (
	"context"
	"net/http"
	"testing"

	"realworldapp/internal/models"
)

func TestTagHandler(t *testing.T) {
	handler := NewTagHandler(tagServiceStub{})
	context, recorder := newHandlerTestContext(http.MethodGet, "/api/tags", "")

	assertHandlerStatus(t, handler.List(context), recorder, http.StatusOK)
}

type tagServiceStub struct{}

func (tagServiceStub) List(_ context.Context) ([]models.Tag, error) {
	return []models.Tag{{ID: 1, Name: "go"}}, nil
}
