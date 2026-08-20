package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func migrationCreateComments() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "202608200004",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS comments (
					id BIGSERIAL PRIMARY KEY,
					body TEXT NOT NULL,
					article_id BIGINT NOT NULL,
					author_id BIGINT NOT NULL,
					created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
					CONSTRAINT fk_comments_article FOREIGN KEY (article_id) REFERENCES articles (id) ON UPDATE CASCADE ON DELETE CASCADE,
					CONSTRAINT fk_comments_author FOREIGN KEY (author_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
				)
			`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_comments_article_id ON comments (article_id)`).Error; err != nil {
				return err
			}
			return tx.Exec(`CREATE INDEX IF NOT EXISTS idx_comments_author_id ON comments (author_id)`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`DROP TABLE IF EXISTS comments`).Error
		},
	}
}
