package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Email        string    `gorm:"size:255;not null;uniqueIndex" json:"email"`
	PasswordHash string    `gorm:"column:password_hash;not null" json:"-"`
	Bio          *string   `gorm:"type:text" json:"bio,omitempty"`
	Image        *string   `gorm:"size:500" json:"image,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
