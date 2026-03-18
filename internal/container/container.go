package container

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joyyth/go-boilerplate/internal/config"
	"github.com/joyyth/go-boilerplate/internal/handler"
	"github.com/joyyth/go-boilerplate/internal/repository"
	"github.com/joyyth/go-boilerplate/internal/service"
	"github.com/rs/zerolog"
)

type Container struct {
	User *handler.UserHandler
	// as you add more domains, register them here:
	// Product *handler.ProductHandler
}

func NewContainer(cfg config.Config, db *pgxpool.Pool, logger *zerolog.Logger) *Container {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, &cfg.Auth)
	userHandler := handler.NewUserHandler(userService, logger)

	return &Container{
		User: userHandler,
	}
}
