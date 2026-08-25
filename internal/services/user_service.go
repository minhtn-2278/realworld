package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"realworldapp/internal/models"
	"realworldapp/internal/repositories"
)

type RegisterUserInput struct {
	Username string
	Email    string
	Password string
}

type LoginUserInput struct {
	Username string
	Password string
}

type LoginUserResult struct {
	AccessToken  string
	RefreshToken string
}

type UserClaims struct {
	jwt.RegisteredClaims
	Username  string `json:"username"`
	TokenType string `json:"token_type"`
}

const (
	jwtIssuer       = "realworldapp"
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 30 * 24 * time.Hour
)

type UserService interface {
	Register(ctx context.Context, input RegisterUserInput) (*models.User, error)
	Login(ctx context.Context, input LoginUserInput) (*LoginUserResult, error)
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetProfile(ctx context.Context, username string) (*models.User, error)
}

type userService struct {
	userRepository repositories.UserRepository
	jwtSecret      string
}

func NewUserService(userRepository repositories.UserRepository, jwtSecret string) UserService {
	return &userService{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
	}
}

func (s *userService) Register(ctx context.Context, input RegisterUserInput) (*models.User, error) {
	passwordHash, err := HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash user password: %w", err)
	}

	user := &models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: passwordHash,
	}
	if err := s.userRepository.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (s *userService) GetProfile(ctx context.Context, username string) (*models.User, error) {
	return s.userRepository.FindByUsername(ctx, username)
}

func (s *userService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	return s.userRepository.FindByID(ctx, id)
}

func (s *userService) Login(ctx context.Context, input LoginUserInput) (*LoginUserResult, error) {
	if input.Username == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := s.userRepository.FindByUsername(ctx, input.Username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user, err = s.userRepository.FindByEmail(ctx, input.Username)
	}
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !VerifyPassword(input.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := s.signToken(user, "access", accessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}
	refreshToken, err := s.signToken(user, "refresh", refreshTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("sign refresh token: %w", err)
	}

	return &LoginUserResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) signToken(user *models.User, tokenType string, lifetime time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			Subject:   strconv.FormatUint(uint64(user.ID), 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(lifetime)),
		},
		Username:  user.Username,
		TokenType: tokenType,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.jwtSecret))
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generate bcrypt hash: %w", err)
	}

	return string(hashedPassword), nil
}

func VerifyPassword(password string, passwordHash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) == nil
}
