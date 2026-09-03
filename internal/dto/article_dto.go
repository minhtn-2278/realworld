package dto

import (
	"time"

	"realworldapp/internal/models"
	"realworldapp/internal/utils"
)

type CreateArticleRequest struct {
	Title       string   `json:"title" validate:"required,max=255"`
	Description string   `json:"description" validate:"required,max=300"`
	Body        string   `json:"body" validate:"required,max=2000"`
	TagList     []string `json:"tagList,omitempty" validate:"dive,max=100"`
}

type UpdateArticleRequest struct {
	Title       string   `json:"title" validate:"required,max=255"`
	Description string   `json:"description" validate:"required,max=300"`
	Body        string   `json:"body" validate:"required,max=2000"`
	TagList     []string `json:"tagList,omitempty" validate:"dive,max=100"`
}

type FavoriteArticleRequest struct {
	Favorite *bool `json:"favorite" validate:"required"`
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

type ArticleDetailResponse struct {
	ArticleResponse
	FavoritesCount int64             `json:"favoritesCount"`
	Comments       []CommentResponse `json:"comments"`
}

type ArticleDetailResponseEnvelope struct {
	Article ArticleDetailResponse `json:"article"`
}

type ArticleListResponse struct {
	Articles []ArticleListItemResponse `json:"articles"`
	Meta     utils.PaginationMeta      `json:"meta"`
}

type ArticleListItemResponse struct {
	ArticleResponse
	FavoritesCount int64 `json:"favoritesCount"`
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

func NewArticleDetailResponse(
	article *models.Article,
	comments []models.Comment,
	favoritesCount int64,
) ArticleDetailResponse {
	return ArticleDetailResponse{
		ArticleResponse: NewArticleResponse(article),
		FavoritesCount:  favoritesCount,
		Comments:        newCommentResponses(comments),
	}
}

func NewArticleListResponse(
	articles []models.Article,
	favoriteCounts map[uint]int64,
	meta utils.PaginationMeta,
) ArticleListResponse {
	items := make([]ArticleListItemResponse, 0, len(articles))
	for index := range articles {
		items = append(items, ArticleListItemResponse{
			ArticleResponse: NewArticleResponse(&articles[index]),
			FavoritesCount:  favoriteCounts[articles[index].ID],
		})
	}

	return ArticleListResponse{Articles: items, Meta: meta}
}
