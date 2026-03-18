package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joyyth/go-boilerplate/internal/config"
	"github.com/joyyth/go-boilerplate/internal/database"
	"github.com/joyyth/go-boilerplate/internal/server"
	"github.com/joyyth/go-boilerplate/pkg/logger"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
        os.Exit(1)
	}

	logger := logger.NewLogger(logger.Options{
		Level:       cfg.Logger.Level,
		Pretty:      cfg.Logger.Pretty,
		ServiceName: "go-boilerplate",
	})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	pool, err := database.NewPostgresPool(ctx, cfg.Database, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database pool")
	}
	defer pool.Close()

	if err := database.RunMigrations(ctx, pool, "./migrations", logger); err != nil {
		logger.Fatal().Err(err).Msg("Failed to apply migrations")
	}

	srv := server.New(*cfg,pool,logger)

	go func() {
		if err := srv.Start(); err != nil {
			logger.Error().Err(err).Msg("Server failed")
			cancel()
		}
	}()
	<-ctx.Done()
	logger.Info().Msg("shutdown signal received, starting graceful shutdown")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("Failed to shutdown server gracefully")
	}
	logger.Info().Msg("Server stopped gracefully")

}
