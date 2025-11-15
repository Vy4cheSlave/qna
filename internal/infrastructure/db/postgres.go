package db

import (
	"context"
	"fmt"
	"github.com/Vy4cheSlave/qna/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewRepository(ctx context.Context, cfg config.PostgreSQL) (*Repository, error) {
	connString := fmt.Sprintf(
		`user=%s password=%s host=%s port=%d dbname=%s sslmode=%s`,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	gormDB, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open gorm connection: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.PoolMaxConns)
	sqlDB.SetConnMaxLifetime(cfg.PoolMaxConnLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.PoolMaxConnIdleTime)

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return &Repository{db: gormDB}, nil

}

type Repository struct {
	db *gorm.DB
}
