package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func migrationCreateArticles() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "202608200003",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS articles (
					id BIGSERIAL PRIMARY KEY,
					slug VARCHAR(255) NOT NULL,
					title VARCHAR(255) NOT NULL,
					description TEXT NOT NULL,
					body TEXT NOT NULL,
					author_id BIGINT NOT NULL,
					created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
					CONSTRAINT uq_articles_slug UNIQUE (slug),
					CONSTRAINT fk_articles_author FOREIGN KEY (author_id)
						REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
				)
			`).Error; err != nil {
				return err
			}
			return tx.Exec(`CREATE INDEX IF NOT EXISTS idx_articles_author_id ON articles (author_id)`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`DROP TABLE IF EXISTS articles`).Error
		},
	}
}
