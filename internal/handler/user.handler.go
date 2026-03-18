package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/joyyth/go-boilerplate/internal/dto"
	"github.com/joyyth/go-boilerplate/internal/service"
	"github.com/joyyth/go-boilerplate/pkg/response"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	userService *service.UserService
	logger      *zerolog.Logger
}

func NewUserHandler(userService *service.UserService, logger *zerolog.Logger) *UserHandler {
	return &UserHandler{userService: userService, logger: logger}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}
	if errs := ValidateRequest(&req); errs != nil {
		response.ValidationError(w, http.StatusBadRequest, errs)
		return
	}
	res, err := h.userService.Register(r.Context(), req)
	if errors.Is(err, service.ErrUserAlreadyExists) {
		response.Error(w, http.StatusConflict, "User already exists")
		return
	}
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to register user")
		response.InternalServerError(w)
		return
	}
	response.Success(w, http.StatusCreated, "User registered successfully", res)
}
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}
	if errs := ValidateRequest(&req); errs != nil {
		response.ValidationError(w, http.StatusBadRequest, errs)
		return
	}
	res, err := h.userService.Login(r.Context(), req)
	if errors.Is(err, service.ErrInvalidCredentials) {
		response.Error(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to login user")
		response.InternalServerError(w)
		return
	}
	response.Success(w, http.StatusOK, "Login successful", res)
}
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}
	if err := h.userService.Logout(r.Context(), req); err != nil {
		h.logger.Error().Err(err).Msg("failed to logout user")
		response.InternalServerError(w)
		return
	}

	response.Success(w, http.StatusOK, "Logout successful", nil)
}
func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}
	if errs := ValidateRequest(&req); errs != nil {
		response.ValidationError(w, http.StatusBadRequest, errs)
		return
	}
	res, err := h.userService.RefreshToken(r.Context(), req)
	if errors.Is(err, service.ErrInvalidToken) {
		response.Error(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to generate new token")
		response.InternalServerError(w)
		return
	}

	response.Success(w, http.StatusOK, "New token generated", res)
}
