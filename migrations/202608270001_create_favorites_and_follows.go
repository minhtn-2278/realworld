package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func migrationCreateFavoritesAndFollows() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "202608270001",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS article_favorites (
					article_id BIGINT NOT NULL,
					user_id BIGINT NOT NULL,
					created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
					PRIMARY KEY (article_id, user_id),
					CONSTRAINT fk_article_favorites_article FOREIGN KEY (article_id) REFERENCES articles (id) ON UPDATE CASCADE ON DELETE CASCADE,
					CONSTRAINT fk_article_favorites_user FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
				)
			`).Error; err != nil {
				return err
			}

			return tx.Exec(`
				CREATE TABLE IF NOT EXISTS user_follows (
					follower_id BIGINT NOT NULL,
					following_id BIGINT NOT NULL,
					created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
					PRIMARY KEY (follower_id, following_id),
					CONSTRAINT fk_user_follows_follower FOREIGN KEY (follower_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
					CONSTRAINT fk_user_follows_following FOREIGN KEY (following_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
				)
			`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec(`DROP TABLE IF EXISTS user_follows`).Error; err != nil {
				return err
			}

			return tx.Exec(`DROP TABLE IF EXISTS article_favorites`).Error
		},
	}
}
