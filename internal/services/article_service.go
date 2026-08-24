package services

import (
	"context"
	"strings"
	"unicode"

	"gorm.io/gorm"

	"realworldapp/internal/models"
	"realworldapp/internal/repositories"
)

type ArticleService interface {
	Create(ctx context.Context, input CreateArticleInput) (*models.Article, error)
	GetBySlug(ctx context.Context, slug string) (*models.Article, error)
	List(ctx context.Context) ([]models.Article, error)
	Update(ctx context.Context, slug string, input UpdateArticleInput) (*models.Article, error)
	Delete(ctx context.Context, slug string) error
	ListComments(ctx context.Context, slug string) ([]models.Comment, error)
	CreateComment(ctx context.Context, slug string, body string, authorID uint) (*models.Comment, error)
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

type articleService struct {
	articleRepository repositories.ArticleRepository
	commentRepository repositories.CommentRepository
	tagRepository     repositories.TagRepository
	userRepository    repositories.UserRepository
	db                *gorm.DB
}

func NewArticleService(
	articleRepository repositories.ArticleRepository,
	commentRepository repositories.CommentRepository,
	tagRepository repositories.TagRepository,
	userRepository repositories.UserRepository,
	db *gorm.DB,
) ArticleService {
	return &articleService{
		articleRepository: articleRepository,
		commentRepository: commentRepository,
		tagRepository:     tagRepository,
		userRepository:    userRepository,
		db:                db,
	}
}

func (s *articleService) Create(ctx context.Context, input CreateArticleInput) (*models.Article, error) {
	tags, err := s.tagRepository.FindOrCreateByNames(ctx, input.TagList)
	if err != nil {
		return nil, err
	}

	article := &models.Article{
		Slug:        slugify(input.Title),
		Title:       input.Title,
		Description: input.Description,
		Body:        input.Body,
		AuthorID:    input.AuthorID,
		Tags:        tags,
	}
	if err := s.articleRepository.Create(ctx, article); err != nil {
		return nil, err
	}

	return article, nil
}

func (s *articleService) GetBySlug(ctx context.Context, slug string) (*models.Article, error) {
	return s.articleRepository.FindBySlug(ctx, slug)
}

func (s *articleService) List(ctx context.Context) ([]models.Article, error) {
	return s.articleRepository.List(ctx)
}

func (s *articleService) Update(ctx context.Context, slug string, input UpdateArticleInput) (updated *models.Article, err error) {
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		articleRepository := repositories.NewArticleRepository(tx)
		tagRepository := repositories.NewTagRepository(tx)

		article, err := articleRepository.UpdateBySlug(ctx, slug, &models.Article{
			Slug:        slugify(input.Title),
			Title:       input.Title,
			Description: input.Description,
			Body:        input.Body,
		})
		if err != nil {
			return err
		}

		if input.TagList != nil {
			tags, err := tagRepository.FindOrCreateByNames(ctx, input.TagList)
			if err != nil {
				return err
			}
			if err := articleRepository.ReplaceTags(ctx, article, tags); err != nil {
				return err
			}
			article.Tags = tags
		}

		updated = article
		return nil
	})

	return updated, err
}

func (s *articleService) Delete(ctx context.Context, slug string) error {
	return s.articleRepository.DeleteBySlug(ctx, slug)
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

	author, err := s.userRepository.FindByID(ctx, authorID)
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
