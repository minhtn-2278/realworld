package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func migrationCreateArticleTags() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "202608200005",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`
				CREATE TABLE IF NOT EXISTS article_tags (
					article_id BIGINT NOT NULL,
					tag_id BIGINT NOT NULL,
					PRIMARY KEY (article_id, tag_id),
					CONSTRAINT fk_article_tags_article FOREIGN KEY (article_id) REFERENCES articles (id) ON UPDATE CASCADE ON DELETE CASCADE,
					CONSTRAINT fk_article_tags_tag FOREIGN KEY (tag_id) REFERENCES tags (id) ON UPDATE CASCADE ON DELETE CASCADE
				)
			`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`DROP TABLE IF EXISTS article_tags`).Error
		},
	}
}
