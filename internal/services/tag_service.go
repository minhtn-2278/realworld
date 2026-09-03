package services

import (
	"context"
	"log/slog"
	"time"

	"realworldapp/internal/cache"
	"realworldapp/internal/models"
	"realworldapp/internal/repositories"
)

const tagsCacheTTL = 5 * time.Minute

type TagService interface {
	List(ctx context.Context) ([]models.Tag, error)
}

type tagService struct {
	tagRepository repositories.TagRepository
	cache         cache.Store
}

func NewTagService(tagRepository repositories.TagRepository, cacheStore cache.Store) TagService {
	return &tagService{tagRepository: tagRepository, cache: cacheStore}
}

func (s *tagService) List(ctx context.Context) ([]models.Tag, error) {
	var tags []models.Tag
	found, err := s.cache.Get(ctx, cache.TagsListKey, &tags)
	if err != nil {
		slog.WarnContext(ctx, "read tags cache failed; using database result", "key", cache.TagsListKey, "error", err)
	}
	if err == nil && found {
		return tags, nil
	}

	tags, err = s.tagRepository.List(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.cache.Set(ctx, cache.TagsListKey, tags, tagsCacheTTL); err != nil {
		slog.WarnContext(ctx, "write tags cache failed", "key", cache.TagsListKey, "error", err)
	}

	return tags, nil
}
