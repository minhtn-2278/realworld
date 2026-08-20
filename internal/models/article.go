package models

import "time"

type Article struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Slug        string    `gorm:"size:255;not null;uniqueIndex" json:"slug"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Body        string    `gorm:"type:text;not null" json:"body"`
	AuthorID    uint      `gorm:"not null;index" json:"-"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
