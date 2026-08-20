package models

type Tag struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:100;not null;uniqueIndex" json:"name"`
}
