package repositories

import (
	"context"

	"gorm.io/gorm"

	"realworldapp/internal/models"
)

type ArticleRepository interface {
	Create(ctx context.Context, article *models.Article) error
	FindBySlug(ctx context.Context, slug string) (*models.Article, error)
	List(ctx context.Context) ([]models.Article, error)
	UpdateBySlug(ctx context.Context, slug string, article *models.Article) (*models.Article, error)
	ReplaceTags(ctx context.Context, article *models.Article, tags []models.Tag) error
	DeleteBySlug(ctx context.Context, slug string) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) Create(ctx context.Context, article *models.Article) error {
	return r.db.WithContext(ctx).Create(article).Error
}

func (r *articleRepository) FindBySlug(ctx context.Context, slug string) (*models.Article, error) {
	var article models.Article
	if err := r.db.WithContext(ctx).
		Preload("Tags").
		Where("slug = ?", slug).
		First(&article).Error; err != nil {
		return nil, err
	}

	return &article, nil
}

func (r *articleRepository) List(ctx context.Context) ([]models.Article, error) {
	var articles []models.Article
	if err := r.db.WithContext(ctx).
		Preload("Tags").
		Order("created_at DESC").
		Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *articleRepository) UpdateBySlug(ctx context.Context, slug string, updates *models.Article) (*models.Article, error) {
	article, err := r.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	article.Slug = updates.Slug
	article.Title = updates.Title
	article.Description = updates.Description
	article.Body = updates.Body
	if err := r.db.WithContext(ctx).Save(article).Error; err != nil {
		return nil, err
	}

	return r.FindBySlug(ctx, article.Slug)
}

func (r *articleRepository) ReplaceTags(ctx context.Context, article *models.Article, tags []models.Tag) error {
	return r.db.WithContext(ctx).Model(article).Association("Tags").Replace(tags)
}

func (r *articleRepository) DeleteBySlug(ctx context.Context, slug string) error {
	result := r.db.WithContext(ctx).Where("slug = ?", slug).Delete(&models.Article{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
