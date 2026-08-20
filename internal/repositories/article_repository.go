package repositories

import (
	"context"

	"gorm.io/gorm"

	"realworldapp/internal/models"
)

type ArticleRepository interface {
	Create(ctx context.Context, article *models.Article) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) Create(ctx context.Context, article *models.Article) error {
	return nil
}
