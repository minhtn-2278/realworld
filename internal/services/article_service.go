package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"unicode"

	"gorm.io/gorm"

	"realworldapp/internal/cache"
	"realworldapp/internal/models"
	"realworldapp/internal/repositories"
	"realworldapp/internal/utils"
)

type ArticleService interface {
	Create(ctx context.Context, input CreateArticleInput) (*models.Article, error)
	GetBySlug(ctx context.Context, slug string) (*models.Article, error)
	GetDetail(ctx context.Context, slug string) (*ArticleDetailResult, error)
	Favorite(ctx context.Context, slug string, userID uint, favorite bool) error
	List(ctx context.Context, input ListArticlesInput) (*ListArticlesResult, error)
	ListFeed(ctx context.Context, userID uint, pagination utils.Pagination) (*ListArticlesResult, error)
	Update(ctx context.Context, slug string, authorID uint, input UpdateArticleInput) (*models.Article, error)
	Delete(ctx context.Context, slug string, authorID uint) error
	ListComments(ctx context.Context, slug string) ([]models.Comment, error)
	CreateComment(ctx context.Context, slug string, body string, authorID uint) (*models.Comment, error)
	GetCommentOwnerID(ctx context.Context, slug string, commentID uint) (uint, error)
	DeleteComment(ctx context.Context, slug string, commentID uint, authorID uint) error
}

type CreateArticleInput struct {
	Title       string
	Description string
	Body        string
	AuthorID    uint
	TagList     []string
}

type UpdateArticleInput struct {
	Title       string
	Description string
	Body        string
	TagList     []string
}

type ListArticlesInput struct {
	Tag        string
	Author     string
	Pagination utils.Pagination
}

type ListArticlesResult struct {
	Articles       []models.Article
	FavoriteCounts map[uint]int64
	Total          int64
}

const feedCacheTTL = 5 * time.Minute

type ArticleDetailResult struct {
	Article        *models.Article
	Comments       []models.Comment
	FavoritesCount int64
}

type articleService struct {
	articleRepository repositories.ArticleRepository
	commentRepository repositories.CommentRepository
	userRepository    repositories.UserRepository
	db                *gorm.DB
	cache             cache.Store
}

func NewArticleService(
	articleRepository repositories.ArticleRepository,
	commentRepository repositories.CommentRepository,
	userRepository repositories.UserRepository,
	db *gorm.DB,
	cacheStore cache.Store,
) ArticleService {
	return &articleService{
		articleRepository: articleRepository,
		commentRepository: commentRepository,
		userRepository:    userRepository,
		db:                db,
		cache:             cacheStore,
	}
}

func (s *articleService) Create(ctx context.Context, input CreateArticleInput) (created *models.Article, err error) {
	slug, err := newArticleSlug(input.Title)
	if err != nil {
		return nil, fmt.Errorf("generate article slug: %w", err)
	}

	var hasNewTags bool
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		articleRepository := repositories.NewArticleRepository(tx)
		tagRepository := repositories.NewTagRepository(tx)

		tags, newTagsCreated, err := tagRepository.FindOrCreateByNames(ctx, input.TagList)
		if err != nil {
			return err
		}
		hasNewTags = newTagsCreated

		article := &models.Article{
			Slug:        slug,
			Title:       input.Title,
			Description: input.Description,
			Body:        input.Body,
			AuthorID:    input.AuthorID,
			Tags:        tags,
		}
		if err := articleRepository.Create(ctx, article); err != nil {
			return err
		}

		created = article
		return nil
	})

	if err != nil {
		return nil, err
	}
	if hasNewTags {
		if err := s.cache.Delete(ctx, cache.TagsListKey); err != nil {
			slog.WarnContext(ctx, "invalidate tags cache failed", "key", cache.TagsListKey, "error", err)
		}
	}
	return created, nil
}

func (s *articleService) GetBySlug(ctx context.Context, slug string) (*models.Article, error) {
	return s.articleRepository.FindBySlug(ctx, slug)
}

func (s *articleService) GetDetail(ctx context.Context, slug string) (*ArticleDetailResult, error) {
	article, err := s.articleRepository.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	favoritesCount, err := s.articleRepository.CountFavorites(ctx, article.ID)
	if err != nil {
		return nil, err
	}
	comments, err := s.commentRepository.ListByArticleID(ctx, article.ID)
	if err != nil {
		return nil, err
	}

	return &ArticleDetailResult{
		Article:        article,
		Comments:       comments,
		FavoritesCount: favoritesCount,
	}, nil
}

func (s *articleService) Favorite(ctx context.Context, slug string, userID uint, favorite bool) error {
	article, err := s.articleRepository.FindBySlug(ctx, slug)
	if err != nil {
		return err
	}

	if favorite {
		if err := s.articleRepository.AddFavorite(ctx, article.ID, userID); err != nil {
			return err
		}
	} else {
		if err := s.articleRepository.RemoveFavorite(ctx, article.ID, userID); err != nil {
			return err
		}
	}

	return nil
}

func (s *articleService) List(ctx context.Context, input ListArticlesInput) (*ListArticlesResult, error) {
	return s.list(ctx, repositories.ArticleListFilter{
		Tag:        input.Tag,
		Author:     input.Author,
		Pagination: input.Pagination,
	})
}

func (s *articleService) ListFeed(ctx context.Context, userID uint, pagination utils.Pagination) (*ListArticlesResult, error) {
	pagination.Limit = utils.DefaultPaginationLimit

	cacheKey := cache.FeedKey(userID, pagination.Limit, pagination.Page)
	var cachedResult ListArticlesResult
	found, err := s.cache.Get(ctx, cacheKey, &cachedResult)
	if err != nil {
		slog.WarnContext(ctx, "read article feed cache failed; using database result", "key", cacheKey, "error", err)
	}
	if err == nil && found {
		return &cachedResult, nil
	}

	result, err := s.list(ctx, repositories.ArticleListFilter{
		FollowerID: userID,
		Pagination: pagination,
	})
	if err != nil {
		return nil, err
	}
	if err := s.cache.Set(ctx, cacheKey, result, feedCacheTTL); err != nil {
		slog.WarnContext(ctx, "write article feed cache failed", "key", cacheKey, "error", err)
	}

	return result, nil
}

func (s *articleService) list(ctx context.Context, filter repositories.ArticleListFilter) (*ListArticlesResult, error) {
	articles, total, err := s.articleRepository.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	articleIDs := make([]uint, 0, len(articles))
	for _, article := range articles {
		articleIDs = append(articleIDs, article.ID)
	}
	favoriteCounts, err := s.articleRepository.CountFavoritesByArticleIDs(ctx, articleIDs)
	if err != nil {
		return nil, err
	}

	return &ListArticlesResult{
		Articles:       articles,
		FavoriteCounts: favoriteCounts,
		Total:          total,
	}, nil
}

func (s *articleService) Update(ctx context.Context, slug string, authorID uint, input UpdateArticleInput) (updated *models.Article, err error) {
	updatedSlug, err := newArticleSlug(input.Title)
	if err != nil {
		return nil, fmt.Errorf("generate article slug: %w", err)
	}

	var hasNewTags bool
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		articleRepository := repositories.NewArticleRepository(tx)
		tagRepository := repositories.NewTagRepository(tx)

		article, err := articleRepository.UpdateBySlug(ctx, slug, authorID, &models.Article{
			Slug:        updatedSlug,
			Title:       input.Title,
			Description: input.Description,
			Body:        input.Body,
		})
		if err != nil {
			return err
		}

		if input.TagList != nil {
			tags, newTagsCreated, err := tagRepository.FindOrCreateByNames(ctx, input.TagList)
			if err != nil {
				return err
			}
			hasNewTags = newTagsCreated
			if err := articleRepository.ReplaceTags(ctx, article, tags); err != nil {
				return err
			}
			article.Tags = tags
		}

		updated = article
		return nil
	})

	if err != nil {
		return nil, err
	}
	if hasNewTags {
		if err := s.cache.Delete(ctx, cache.TagsListKey); err != nil {
			slog.WarnContext(ctx, "invalidate tags cache failed", "key", cache.TagsListKey, "error", err)
		}
	}
	return updated, nil
}

func (s *articleService) Delete(ctx context.Context, slug string, authorID uint) error {
	if err := s.articleRepository.DeleteBySlug(ctx, slug, authorID); err != nil {
		return err
	}

	return nil
}

func (s *articleService) ListComments(ctx context.Context, slug string) ([]models.Comment, error) {
	article, err := s.articleRepository.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return s.commentRepository.ListByArticleID(ctx, article.ID)
}

func (s *articleService) CreateComment(ctx context.Context, slug string, body string, authorID uint) (*models.Comment, error) {
	article, err := s.articleRepository.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	author, err := s.userRepository.FindByID(ctx, authorID, false)
	if err != nil {
		return nil, err
	}

	comment := &models.Comment{
		Body:      body,
		ArticleID: article.ID,
		AuthorID:  author.ID,
	}
	if err := s.commentRepository.Create(ctx, comment); err != nil {
		return nil, err
	}
	comment.Author = *author

	return comment, nil
}

func (s *articleService) GetCommentOwnerID(ctx context.Context, slug string, commentID uint) (uint, error) {
	article, err := s.articleRepository.FindBySlug(ctx, slug)
	if err != nil {
		return 0, err
	}

	comment, err := s.commentRepository.FindByIDAndArticleID(ctx, commentID, article.ID)
	if err != nil {
		return 0, err
	}

	return comment.AuthorID, nil
}

func (s *articleService) DeleteComment(ctx context.Context, slug string, commentID uint, authorID uint) error {
	article, err := s.articleRepository.FindBySlug(ctx, slug)
	if err != nil {
		return err
	}

	return s.commentRepository.DeleteByIDAndArticleIDAndAuthorID(ctx, commentID, article.ID, authorID)
}

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var builder strings.Builder
	previousDash := false

	for _, char := range value {
		switch {
		case unicode.IsLetter(char) || unicode.IsDigit(char):
			builder.WriteRune(char)
			previousDash = false
		case !previousDash && builder.Len() > 0:
			builder.WriteByte('-')
			previousDash = true
		}
	}

	return strings.Trim(builder.String(), "-")
}

func newArticleSlug(title string) (string, error) {
	randomSuffix := make([]byte, 4)
	if _, err := rand.Read(randomSuffix); err != nil {
		return "", fmt.Errorf("read random slug suffix: %w", err)
	}

	baseSlug := slugify(title)
	suffix := hex.EncodeToString(randomSuffix)
	if baseSlug == "" {
		return suffix, nil
	}

	return baseSlug + "-" + suffix, nil
}
