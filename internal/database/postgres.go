package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joyyth/go-boilerplate/internal/config"
	"github.com/rs/zerolog"
)

func DatabaseURL(cfg config.DatabaseConfig) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
}

func NewPostgresPool(ctx context.Context, cfg config.DatabaseConfig, logger zerolog.Logger) (*pgxpool.Pool, error) {
	poolconfig, err := pgxpool.ParseConfig(DatabaseURL(cfg))
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	poolconfig.MaxConns = int32(cfg.MaxOpenConns)
	poolconfig.MinConns = int32(cfg.MaxIdleConns)
	poolconfig.MaxConnLifetime = cfg.ConnMaxLifetime
	poolconfig.MaxConnIdleTime = cfg.ConnMaxIdleTime
	poolconfig.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, poolconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	logger.Info().
		Str("host", cfg.Host).
		Int("port", cfg.Port).
		Str("database", cfg.Name).
		Msg("connected to postgres")

	return pool, nil
}
