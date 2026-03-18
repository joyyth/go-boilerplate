package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/joyyth/go-boilerplate/internal/auth"
	"github.com/joyyth/go-boilerplate/pkg/response"
)

type contextKey string

const UserClaimKey contextKey = "user_claims"

func RequireAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				response.Unauthorized(w, "No token found")
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := auth.VerifyToken(tokenStr, secret)
			if err != nil {
				response.Unauthorized(w, "Invalid token")
				return
			}
			ctx := context.WithValue(r.Context(), UserClaimKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}

func GetUserClaims(r *http.Request) (*auth.Claims, error) {
	claims, ok := r.Context().Value(UserClaimKey).(*auth.Claims)
	if !ok {
		return nil, fmt.Errorf("no user claims found in context")
	}
	return claims, nil
}
