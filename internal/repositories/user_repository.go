package repositories

import (
	"context"

	"gorm.io/gorm"

	"realworldapp/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return nil
}
