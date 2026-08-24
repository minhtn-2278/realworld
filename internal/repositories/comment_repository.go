package repositories

import (
	"context"

	"gorm.io/gorm"

	"realworldapp/internal/models"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *models.Comment) error
	ListByArticleID(ctx context.Context, articleID uint) ([]models.Comment, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *models.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *commentRepository) ListByArticleID(ctx context.Context, articleID uint) ([]models.Comment, error) {
	var comments []models.Comment
	if err := r.db.WithContext(ctx).
		Preload("Author").
		Where("article_id = ?", articleID).
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}
