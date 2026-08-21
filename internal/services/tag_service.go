package services

import (
	"context"

	"realworldapp/internal/models"
	"realworldapp/internal/repositories"
)

type TagService interface {
	List(ctx context.Context) ([]models.Tag, error)
}

type tagService struct {
	tagRepository repositories.TagRepository
}

func NewTagService(tagRepository repositories.TagRepository) TagService {
	return &tagService{tagRepository: tagRepository}
}

func (s *tagService) List(ctx context.Context) ([]models.Tag, error) {
	return s.tagRepository.List(ctx)
}
