package handlers

import (
	"context"
	"net/http"
	"testing"

	"github.com/labstack/echo/v5"

	"realworldapp/internal/models"
	"realworldapp/internal/services"
	"realworldapp/internal/utils"
)

func TestArticleHandlers(t *testing.T) {
	handler := NewArticleHandler(articleServiceStub{})

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
		{"create", http.MethodPost, "/api/articles", `{"title":"Post","description":"Description","body":"Body","tagList":["go"]}`, nil, true, http.StatusCreated, handler.Create},
		{"list", http.MethodGet, "/api/articles?tag=go&author=alice", "", nil, false, http.StatusOK, handler.List},
		{"list feed", http.MethodGet, "/api/articles/feed?page=2", "", nil, true, http.StatusOK, handler.ListFeed},
		{"get by slug", http.MethodGet, "/api/articles/post", "", []string{"slug", "post"}, false, http.StatusOK, handler.GetBySlug},
		{"favorite", http.MethodPost, "/api/articles/post/favorite", `{"favorite":true}`, []string{"slug", "post"}, true, http.StatusOK, handler.Favorite},
		{"update", http.MethodPut, "/api/articles/post", `{"title":"Post","description":"Description","body":"Body","tagList":["go"]}`, []string{"slug", "post"}, true, http.StatusOK, handler.Update},
		{"delete", http.MethodDelete, "/api/articles/post", "", []string{"slug", "post"}, true, http.StatusOK, handler.Delete},
		{"list comments", http.MethodGet, "/api/articles/post/comments", "", []string{"slug", "post"}, false, http.StatusOK, handler.ListComments},
		{"create comment", http.MethodPost, "/api/articles/post/comments", `{"body":"Nice article"}`, []string{"slug", "post"}, true, http.StatusCreated, handler.CreateComment},
		{"delete comment", http.MethodDelete, "/api/articles/post/comments/1", "", []string{"slug", "post", "id", "1"}, true, http.StatusOK, handler.DeleteComment},
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

type articleServiceStub struct{}

func (articleServiceStub) Create(_ context.Context, _ services.CreateArticleInput) (*models.Article, error) {
	return articleTestData(), nil
}

func (articleServiceStub) GetBySlug(_ context.Context, _ string) (*models.Article, error) {
	return articleTestData(), nil
}

func (articleServiceStub) GetDetail(_ context.Context, _ string) (*services.ArticleDetailResult, error) {
	return &services.ArticleDetailResult{Article: articleTestData(), Comments: []models.Comment{articleCommentTestData()}, FavoritesCount: 1}, nil
}

func (articleServiceStub) Favorite(_ context.Context, _ string, _ uint, _ bool) error { return nil }

func (articleServiceStub) List(_ context.Context, _ services.ListArticlesInput) (*services.ListArticlesResult, error) {
	return articleListTestData(), nil
}

func (articleServiceStub) ListFeed(_ context.Context, _ uint, _ utils.Pagination) (*services.ListArticlesResult, error) {
	return articleListTestData(), nil
}

func (articleServiceStub) Update(_ context.Context, _ string, _ uint, _ services.UpdateArticleInput) (*models.Article, error) {
	return articleTestData(), nil
}

func (articleServiceStub) Delete(_ context.Context, _ string, _ uint) error { return nil }

func (articleServiceStub) ListComments(_ context.Context, _ string) ([]models.Comment, error) {
	return []models.Comment{articleCommentTestData()}, nil
}

func (articleServiceStub) CreateComment(_ context.Context, _ string, _ string, _ uint) (*models.Comment, error) {
	comment := articleCommentTestData()
	return &comment, nil
}

func (articleServiceStub) GetCommentOwnerID(_ context.Context, _ string, _ uint) (uint, error) {
	return 1, nil
}

func (articleServiceStub) DeleteComment(_ context.Context, _ string, _ uint, _ uint) error {
	return nil
}

func articleTestData() *models.Article {
	return &models.Article{
		ID:          1,
		Slug:        "post",
		Title:       "Post",
		Description: "Description",
		Body:        "Body",
		AuthorID:    1,
		Tags:        []models.Tag{{ID: 1, Name: "go"}},
	}
}

func articleListTestData() *services.ListArticlesResult {
	article := articleTestData()
	return &services.ListArticlesResult{
		Articles:       []models.Article{*article},
		FavoriteCounts: map[uint]int64{article.ID: 1},
		Total:          1,
	}
}

func articleCommentTestData() models.Comment {
	return models.Comment{ID: 1, Body: "Nice article", AuthorID: 1, Author: models.User{ID: 1, Username: "alice"}}
}
