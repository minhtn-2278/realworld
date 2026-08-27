package models

import "time"

type UserFollow struct {
	FollowerID  uint      `gorm:"primaryKey;not null"`
	FollowingID uint      `gorm:"primaryKey;not null"`
	CreatedAt   time.Time `gorm:"not null"`
}
