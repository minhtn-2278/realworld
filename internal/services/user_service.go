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
	"realworldapp/internal/utils"
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
	GetProfileForUser(ctx context.Context, username string, viewerID uint) (*models.User, error)
	Follow(ctx context.Context, username string, followerID uint, follow bool) error
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
	user, err := s.userRepository.FindByUsername(ctx, username, true)
	if err != nil {
		return nil, err
	}

	if err := s.populateFollowersCount(ctx, user); err != nil {
		return nil, err
	}
	if err := s.populateFollowingCount(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetProfileForUser(ctx context.Context, username string, viewerID uint) (*models.User, error) {
	user, err := s.GetProfile(ctx, username)
	if err != nil {
		return nil, err
	}
	following, err := s.userRepository.IsFollowing(ctx, viewerID, user.ID)
	if err != nil {
		return nil, err
	}
	user.Following = following

	return user, nil
}

func (s *userService) Follow(ctx context.Context, username string, followerID uint, follow bool) error {
	user, err := s.userRepository.FindByUsername(ctx, username, false)
	if err != nil {
		return err
	}
	if follow && user.ID == followerID {
		return utils.ErrCannotFollowSelf
	}
	if follow {
		if err := s.userRepository.AddFollow(ctx, followerID, user.ID); err != nil {
			return err
		}
	} else {
		if err := s.userRepository.RemoveFollow(ctx, followerID, user.ID); err != nil {
			return err
		}
	}

	return nil
}

func (s *userService) populateFollowersCount(ctx context.Context, user *models.User) error {
	count, err := s.userRepository.CountFollowers(ctx, user.ID)
	if err != nil {
		return err
	}

	user.FollowersCount = count
	return nil
}

func (s *userService) populateFollowingCount(ctx context.Context, user *models.User) error {
	count, err := s.userRepository.CountFollowing(ctx, user.ID)
	if err != nil {
		return err
	}

	user.FollowingCount = count
	return nil
}

func (s *userService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	return s.userRepository.FindByID(ctx, id, false)
}

func (s *userService) Login(ctx context.Context, input LoginUserInput) (*LoginUserResult, error) {
	if input.Username == "" {
		return nil, utils.ErrInvalidCredentials
	}

	user, err := s.userRepository.FindByUsername(ctx, input.Username, false)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user, err = s.userRepository.FindByEmail(ctx, input.Username)
	}
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find user for login: %w", err)
		}
		return nil, utils.ErrInvalidCredentials
	}

	if !VerifyPassword(input.Password, user.PasswordHash) {
		return nil, utils.ErrInvalidCredentials
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
