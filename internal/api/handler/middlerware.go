package handler

import (
	"context"
	"net/http"

	"github.com/cheezecakee/fitrkr/internal/utils/auth"
	"github.com/cheezecakee/fitrkr/internal/utils/helper"
)

type AuthMiddleware struct {
	JWTManager auth.JWT
}

func NewAuthMiddleware(jwtMgr auth.JWT) *AuthMiddleware {
	return &AuthMiddleware{
		JWTManager: jwtMgr,
	}
}

func (m *AuthMiddleware) IsAuthenticated() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := helper.GetBearerToken(r.Header)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			userID, err := m.JWTManager.ValidateJWT(token)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; fontsrc fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}
