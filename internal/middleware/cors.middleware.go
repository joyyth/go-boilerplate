package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
	"github.com/joyyth/go-boilerplate/internal/config"
)

func CorsMiddleware(cfg config.ServerConfig) func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-Id"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
