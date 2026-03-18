package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joyyth/go-boilerplate/internal/container"
	internal_middleware "github.com/joyyth/go-boilerplate/internal/middleware"
	"github.com/joyyth/go-boilerplate/pkg/response"
)

func (s *Server) MountRoutes() {
	c := container.NewContainer(s.cfg, s.db, &s.logger)
	s.router.Get("/health", s.handleHealthCheck())

	s.router.Route("/api/v1", func(r chi.Router) {
		// auth routes (public)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", c.User.Register)
			r.Post("/login", c.User.Login)
			r.Post("/refresh", c.User.RefreshToken)
			r.Post("/logout", c.User.Logout)
		})

		// protected routes
		r.Group(func(r chi.Router) {
			r.Use(internal_middleware.RequireAuth(s.cfg.Auth.AccessSecret))

			// r.Get("/products", c.Product.List)

		})
	})
}
func (s *Server) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.db.Ping(r.Context()); err != nil {
			s.logger.Error().Err(err).Msg("health check failed")
			response.Error(w, http.StatusServiceUnavailable, "database unreachable")
			return
		}

		response.Success(w, http.StatusOK, "Health check successful", nil)
	}
}