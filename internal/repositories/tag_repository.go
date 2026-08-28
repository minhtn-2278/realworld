package repositories

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"realworldapp/internal/models"
)

type TagRepository interface {
	List(ctx context.Context) ([]models.Tag, error)
	FindOrCreateByNames(ctx context.Context, names []string) ([]models.Tag, bool, error)
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

func (r *tagRepository) FindOrCreateByNames(ctx context.Context, names []string) ([]models.Tag, bool, error) {
	unique := normalizeTagNames(names)
	if len(unique) == 0 {
		return []models.Tag{}, false, nil
	}

	toInsert := make([]models.Tag, 0, len(unique))
	for _, name := range unique {
		toInsert = append(toInsert, models.Tag{Name: name})
	}

	result := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&toInsert)
	if result.Error != nil {
		return nil, false, result.Error
	}

	var found []models.Tag
	if err := r.db.WithContext(ctx).Where("name IN ?", unique).Find(&found).Error; err != nil {
		return nil, false, err
	}

	byName := make(map[string]models.Tag, len(found))
	for _, tag := range found {
		byName[tag.Name] = tag
	}
	tags := make([]models.Tag, 0, len(unique))
	for _, name := range unique {
		if tag, ok := byName[name]; ok {
			tags = append(tags, tag)
		}
	}

	return tags, result.RowsAffected > 0, nil
}

func normalizeTagNames(names []string) []string {
	seen := make(map[string]struct{}, len(names))
	unique := make([]string, 0, len(names))
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}
		unique = append(unique, name)
	}

	return unique
}
