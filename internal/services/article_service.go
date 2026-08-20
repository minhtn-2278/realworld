package services

import (
	"context"

	"realworldapp/internal/models"
	"realworldapp/internal/repositories"
)

type ArticleService interface {
	Create(ctx context.Context, article *models.Article) error
}

type articleService struct {
	articleRepository repositories.ArticleRepository
}

func NewArticleService(articleRepository repositories.ArticleRepository) ArticleService {
	return &articleService{
		articleRepository: articleRepository,
	}
}

func (s *articleService) Create(ctx context.Context, article *models.Article) error {
	return nil
}
