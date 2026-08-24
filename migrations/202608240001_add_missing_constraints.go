package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func migrationAddMissingConstraints() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "202608240001",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec(`
				ALTER TABLE users
					ADD CONSTRAINT uq_users_username UNIQUE (username);
			`).Error; err != nil {
				return err
			}

			if err := tx.Exec(`
				ALTER TABLE users
					ADD CONSTRAINT uq_users_email UNIQUE (email);
			`).Error; err != nil {
				return err
			}

			if err := tx.Exec(`
				ALTER TABLE articles
					ADD CONSTRAINT fk_articles_author
					FOREIGN KEY (author_id) REFERENCES users (id)
					ON UPDATE CASCADE ON DELETE CASCADE;
			`).Error; err != nil {
				return err
			}

			if err := tx.Exec(`
				CREATE INDEX IF NOT EXISTS idx_articles_author_id
				ON articles (author_id)
			`).Error; err != nil {
				return err
			}

			return tx.Exec(`
				CREATE INDEX IF NOT EXISTS idx_article_tags_tag_id
				ON article_tags (tag_id)
			`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec(`DROP INDEX IF EXISTS idx_article_tags_tag_id`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`DROP INDEX IF EXISTS idx_articles_author_id`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`ALTER TABLE articles DROP CONSTRAINT IF EXISTS fk_articles_author`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`ALTER TABLE users DROP CONSTRAINT IF EXISTS uq_users_email`).Error; err != nil {
				return err
			}
			return tx.Exec(`ALTER TABLE users DROP CONSTRAINT IF EXISTS uq_users_username`).Error
		},
	}
}
