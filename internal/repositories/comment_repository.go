package repositories

import (
	"context"

	"gorm.io/gorm"

	"realworldapp/internal/models"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *models.Comment) error
	ListByArticleID(ctx context.Context, articleID uint) ([]models.Comment, error)
	FindByIDAndArticleID(ctx context.Context, commentID uint, articleID uint) (*models.Comment, error)
	DeleteByIDAndArticleIDAndAuthorID(ctx context.Context, commentID uint, articleID uint, authorID uint) error
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

func (r *commentRepository) FindByIDAndArticleID(ctx context.Context, commentID uint, articleID uint) (*models.Comment, error) {
	var comment models.Comment
	if err := r.db.WithContext(ctx).
		Where("id = ? AND article_id = ?", commentID, articleID).
		First(&comment).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *commentRepository) DeleteByIDAndArticleIDAndAuthorID(ctx context.Context, commentID uint, articleID uint, authorID uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND article_id = ? AND author_id = ?", commentID, articleID, authorID).
		Delete(&models.Comment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
