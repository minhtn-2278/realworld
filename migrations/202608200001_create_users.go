package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func migrationCreateUsers() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "202608200001",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`
				CREATE TABLE IF NOT EXISTS users (
					id BIGSERIAL PRIMARY KEY,
					username VARCHAR(50) NOT NULL,
					email VARCHAR(255) NOT NULL,
					password_hash TEXT NOT NULL,
					bio TEXT,
					image VARCHAR(500),
					created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
				)
			`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`DROP TABLE IF EXISTS users`).Error
		},
	}
}
