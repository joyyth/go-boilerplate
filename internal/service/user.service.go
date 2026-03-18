package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/joyyth/go-boilerplate/internal/auth"
	"github.com/joyyth/go-boilerplate/internal/config"
	"github.com/joyyth/go-boilerplate/internal/dto"
	"github.com/joyyth/go-boilerplate/internal/model"
	"github.com/joyyth/go-boilerplate/internal/repository"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

type UserService struct {
	userRepo *repository.UserRepository
	authcfg  *config.AuthConfig
}

func NewUserService(userRepo *repository.UserRepository, authcfg *config.AuthConfig) *UserService {
	return &UserService{userRepo: userRepo, authcfg: authcfg}
}

func (s *UserService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	_, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("failed to check if user exists: %w", err)
	}
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user, err := s.userRepo.CreateUser(ctx, req.Email, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return s.generateAuthResponse(ctx, user)
}

func (s *UserService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}
	return s.generateAuthResponse(ctx, user)
}

func (s *UserService) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	tokenHash := hashToken(req.RefreshToken)

	storedToken, err := s.userRepo.GetRefreshToken(ctx, tokenHash)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrInvalidToken
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}
	if err := s.userRepo.DeleteRefreshToken(ctx, tokenHash); err != nil {
		return nil, fmt.Errorf("failed to delete old refresh token: %w", err)
	}
	user, err := s.userRepo.GetUserByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return s.generateAuthResponse(ctx, user)
}

func (s *UserService) Logout(ctx context.Context, req dto.RefreshTokenRequest) error {
	tokenHash := hashToken(req.RefreshToken)
	if err := s.userRepo.DeleteRefreshToken(ctx, tokenHash); err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	return nil
}

// helper function to generate auth response
func (s *UserService) generateAuthResponse(ctx context.Context, user *model.User) (*dto.AuthResponse, error) {
	tokenPair, err := auth.GenerateTokenPair(user.ID, user.Email, auth.JWTConfig{
		AccessTokenSecret:      s.authcfg.AccessSecret,
		RefreshTokenSecret:     s.authcfg.RefreshSecret,
		AccessTokenExpiration:  s.authcfg.AccessTokenExpiry,
		RefreshTokenExpiration: s.authcfg.RefreshTokenExpiry,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}
	tokenHash := hashToken(tokenPair.RefreshToken)
	expiresAt := time.Now().Add(s.authcfg.RefreshTokenExpiry).Unix()

	if _, err := s.userRepo.CreateRefreshToken(ctx, user.ID, tokenHash, expiresAt); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}
	return &dto.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}
