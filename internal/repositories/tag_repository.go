package repositories

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"realworldapp/internal/models"
)

type TagRepository interface {
	List(ctx context.Context) ([]models.Tag, error)
	FindOrCreateByNames(ctx context.Context, names []string) ([]models.Tag, error)
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) List(ctx context.Context) ([]models.Tag, error) {
	var tags []models.Tag
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *tagRepository) FindOrCreateByNames(ctx context.Context, names []string) ([]models.Tag, error) {
	tags := make([]models.Tag, 0, len(names))
	seen := make(map[string]struct{}, len(names))

	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}

		var tag models.Tag
		if err := r.db.WithContext(ctx).Where("name = ?", name).FirstOrCreate(&tag, models.Tag{Name: name}).Error; err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
