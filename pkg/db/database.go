package db

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("database DSN is empty")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open PostgreSQL connection: %w", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		return nil, fmt.Errorf("get PostgreSQL connection pool: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := ping(sqlDB); err != nil {
		return nil, err
	}

	return database, nil
}

func ping(database *sql.DB) error {
	if err := database.Ping(); err != nil {
		return fmt.Errorf("ping PostgreSQL: %w", err)
	}
	return nil
}
