package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func migrationCreateTags() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "202608200002",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`
				CREATE TABLE IF NOT EXISTS tags (
					id BIGSERIAL PRIMARY KEY,
					name VARCHAR(100) NOT NULL,
					CONSTRAINT uq_tags_name UNIQUE (name)
				)
			`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`DROP TABLE IF EXISTS tags`).Error
		},
	}
}
