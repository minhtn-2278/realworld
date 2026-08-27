package repositories

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"realworldapp/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id uint) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	AddFollow(ctx context.Context, followerID uint, followingID uint) error
	RemoveFollow(ctx context.Context, followerID uint, followingID uint) error
	IsFollowing(ctx context.Context, followerID uint, followingID uint) (bool, error)
	CountFollowers(ctx context.Context, userID uint) (int64, error)
	CountFollowing(ctx context.Context, userID uint) (int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Preload("Articles", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, title, author_id").Order("created_at DESC")
		}).
		Where("id = ?", id).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Preload("Articles", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, title, author_id").Order("created_at DESC")
		}).
		Where("username = ?", username).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) AddFollow(ctx context.Context, followerID uint, followingID uint) error {
	follow := &models.UserFollow{FollowerID: followerID, FollowingID: followingID}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(follow).Error
}

func (r *userRepository) RemoveFollow(ctx context.Context, followerID uint, followingID uint) error {
	return r.db.WithContext(ctx).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&models.UserFollow{}).Error
}

func (r *userRepository) IsFollowing(ctx context.Context, followerID uint, followingID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.UserFollow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepository) CountFollowers(ctx context.Context, userID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.UserFollow{}).
		Where("following_id = ?", userID).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *userRepository) CountFollowing(ctx context.Context, userID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.UserFollow{}).
		Where("follower_id = ?", userID).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
