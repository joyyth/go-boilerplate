package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joyyth/go-boilerplate/migrations"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
)

func RunMigrations(ctx context.Context, pool *pgxpool.Pool, migrationsDir string, logger zerolog.Logger) error {

	goose.SetLogger(goose.NopLogger())

	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}
	goose.SetBaseFS(migrations.FS)

	if err := goose.UpContext(ctx, db, "."); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	logger.Info().Str("dir", migrationsDir).Msg("migrations applied successfully")
	return nil
}
