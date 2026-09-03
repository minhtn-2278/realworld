package repositories

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"realworldapp/internal/models"
	"realworldapp/internal/utils"
)

type ArticleRepository interface {
	Create(ctx context.Context, article *models.Article) error
	FindBySlug(ctx context.Context, slug string) (*models.Article, error)
	List(ctx context.Context, filter ArticleListFilter) ([]models.Article, int64, error)
	UpdateBySlug(ctx context.Context, slug string, authorID uint, article *models.Article) (*models.Article, error)
	ReplaceTags(ctx context.Context, article *models.Article, tags []models.Tag) error
	DeleteBySlug(ctx context.Context, slug string, authorID uint) error
	AddFavorite(ctx context.Context, articleID uint, userID uint) error
	RemoveFavorite(ctx context.Context, articleID uint, userID uint) error
	CountFavorites(ctx context.Context, articleID uint) (int64, error)
	CountFavoritesByArticleIDs(ctx context.Context, articleIDs []uint) (map[uint]int64, error)
}

type ArticleListFilter struct {
	Tag        string
	Author     string
	FollowerID uint
	Pagination utils.Pagination
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

func (r *articleRepository) List(ctx context.Context, filter ArticleListFilter) ([]models.Article, int64, error) {
	countQuery := r.db.WithContext(ctx).Model(&models.Article{})
	countQuery = applyArticleListFilters(countQuery, filter)

	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var articles []models.Article
	query := r.db.WithContext(ctx).
		Model(&models.Article{}).
		Preload("Tags").
		Order("articles.created_at DESC, articles.id DESC")
	query = applyArticleListFilters(query, filter)

	if filter.Pagination.Limit > 0 {
		query = query.Limit(filter.Pagination.Limit)
	}
	if filter.Pagination.Offset() > 0 {
		query = query.Offset(filter.Pagination.Offset())
	}

	if err := query.Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

func applyArticleListFilters(query *gorm.DB, filter ArticleListFilter) *gorm.DB {
	if filter.Tag != "" {
		query = query.Where(`EXISTS (
			SELECT 1
			FROM article_tags
			JOIN tags ON tags.id = article_tags.tag_id
			WHERE article_tags.article_id = articles.id
			  AND tags.name = ?
		)`, filter.Tag)
	}

	if filter.Author != "" {
		query = query.Joins("JOIN users ON users.id = articles.author_id").
			Where("users.username = ?", filter.Author)
	}

	if filter.FollowerID != 0 {
		query = query.Where(`EXISTS (
			SELECT 1
			FROM user_follows
			WHERE user_follows.follower_id = ?
			  AND user_follows.following_id = articles.author_id
		)`, filter.FollowerID)
	}

	return query
}

func (r *articleRepository) UpdateBySlug(ctx context.Context, slug string, authorID uint, updates *models.Article) (*models.Article, error) {
	var article models.Article
	if err := r.db.WithContext(ctx).
		Where("slug = ? AND author_id = ?", slug, authorID).
		First(&article).Error; err != nil {
		return nil, err
	}

	if article.Title != updates.Title {
		article.Slug = updates.Slug
	}
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

func (r *articleRepository) DeleteBySlug(ctx context.Context, slug string, authorID uint) error {
	result := r.db.WithContext(ctx).
		Where("slug = ? AND author_id = ?", slug, authorID).
		Delete(&models.Article{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *articleRepository) AddFavorite(ctx context.Context, articleID uint, userID uint) error {
	favorite := &models.ArticleFavorite{ArticleID: articleID, UserID: userID}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(favorite).Error
}

func (r *articleRepository) RemoveFavorite(ctx context.Context, articleID uint, userID uint) error {
	return r.db.WithContext(ctx).
		Where("article_id = ? AND user_id = ?", articleID, userID).
		Delete(&models.ArticleFavorite{}).Error
}

func (r *articleRepository) CountFavorites(ctx context.Context, articleID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.ArticleFavorite{}).
		Where("article_id = ?", articleID).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *articleRepository) CountFavoritesByArticleIDs(ctx context.Context, articleIDs []uint) (map[uint]int64, error) {
	countsByArticleID := make(map[uint]int64, len(articleIDs))
	if len(articleIDs) == 0 {
		return countsByArticleID, nil
	}

	var counts []struct {
		ArticleID      uint
		FavoritesCount int64
	}
	if err := r.db.WithContext(ctx).
		Model(&models.ArticleFavorite{}).
		Select("article_id, COUNT(*) AS favorites_count").
		Where("article_id IN ?", articleIDs).
		Group("article_id").
		Scan(&counts).Error; err != nil {
		return nil, err
	}

	for _, count := range counts {
		countsByArticleID[count.ArticleID] = count.FavoritesCount
	}

	return countsByArticleID, nil
}
