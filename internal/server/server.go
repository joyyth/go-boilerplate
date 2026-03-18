package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joyyth/go-boilerplate/internal/config"
	internal_middleware "github.com/joyyth/go-boilerplate/internal/middleware"
	"github.com/rs/zerolog"
)

type Server struct {
	httpServer *http.Server
	cfg        config.Config
	router     *chi.Mux
	db         *pgxpool.Pool
	logger     zerolog.Logger
}

func New(cfg config.Config, db *pgxpool.Pool, logger zerolog.Logger) *Server {
	s := &Server{
		router: chi.NewRouter(),
		db:     db,
		logger: logger,
		cfg:    cfg,
	}
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      s.router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}
	s.MountMiddleware()
	s.MountRoutes()
	return s
}

func (s *Server) MountMiddleware() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Compress(5))
	// rate limiting — 100 requests per minute per IP
	s.router.Use(httprate.LimitByIP(100, time.Minute))

	s.router.Use(internal_middleware.LoggerMiddleware(s.logger))
	s.router.Use(internal_middleware.CorsMiddleware(s.cfg.Server))
}

func (s *Server) Start() error {
	s.logger.Info().Str("addr", s.httpServer.Addr).Msg("starting server")
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info().Msg("shutting down server")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}
	s.logger.Info().Msg("server stopped cleanly")
	return nil
}
