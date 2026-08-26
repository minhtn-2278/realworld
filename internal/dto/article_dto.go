package dto

import (
	"time"

	"realworldapp/internal/models"
)

type CreateArticleRequest struct {
	Title       string   `json:"title" validate:"required,max=255"`
	Description string   `json:"description" validate:"required"`
	Body        string   `json:"body" validate:"required"`
	TagList     []string `json:"tagList,omitempty" validate:"dive,max=100"`
}

type UpdateArticleRequest struct {
	Title       string   `json:"title" validate:"required,max=255"`
	Description string   `json:"description" validate:"required"`
	Body        string   `json:"body" validate:"required"`
	TagList     []string `json:"tagList,omitempty" validate:"dive,max=100"`
}

type ArticleResponse struct {
	ID          uint      `json:"id"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	TagList     []string  `json:"tagList"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ArticleResponseEnvelope struct {
	Article ArticleResponse `json:"article"`
}

type ArticleListResponse struct {
	Articles []ArticleResponse `json:"articles"`
}

func NewArticleResponse(article *models.Article) ArticleResponse {
	tagList := make([]string, 0, len(article.Tags))
	for _, tag := range article.Tags {
		tagList = append(tagList, tag.Name)
	}

	return ArticleResponse{
		ID:          article.ID,
		Slug:        article.Slug,
		Title:       article.Title,
		Description: article.Description,
		Body:        article.Body,
		TagList:     tagList,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
	}
}

func NewArticleListResponse(articles []models.Article) ArticleListResponse {
	items := make([]ArticleResponse, 0, len(articles))
	for index := range articles {
		items = append(items, NewArticleResponse(&articles[index]))
	}

	return ArticleListResponse{Articles: items}
}
