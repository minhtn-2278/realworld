package models

import "time"

type ArticleFavorite struct {
	ArticleID uint      `gorm:"primaryKey;not null"`
	UserID    uint      `gorm:"primaryKey;not null"`
	CreatedAt time.Time `gorm:"not null"`
}
