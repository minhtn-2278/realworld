package services

import (
	"context"

	"realworldapp/internal/models"
	"realworldapp/internal/repositories"
)

type UserService interface {
	Create(ctx context.Context, user *models.User) error
	GetProfile(ctx context.Context, username string) (*models.User, error)
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) Create(ctx context.Context, user *models.User) error {
	return nil
}

func (s *userService) GetProfile(ctx context.Context, username string) (*models.User, error) {
	return s.userRepository.FindByUsername(ctx, username)
}
