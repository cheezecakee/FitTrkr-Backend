package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/cheezecakee/FitLogr/internal/services"
)

func (cfg *Config) IsAuthenticated() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Prevent caching of responses
			w.Header().Add("Cache-Control", "no-store")

			token, err := cfg.Helper.GetBearerToken(r.Header)
			if err != nil {
				cfg.Logger.ErrorLog.Printf("Failed to get berear token: %v", err)
				cfg.RedirectToLogin(w, r)
				return
			}

			userID, err := cfg.JWTManager.ValidateJWT(token)
			if err == nil {
				ctx := context.WithValue(r.Context(), cfg.ContextKey["userID"], userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			cfg.Logger.ErrorLog.Printf("Invalid JWT token: %v", err)

			cfg.RefreshToken(w, r, userID, next)
		})
	}
}

func (cfg *Config) RefreshToken(w http.ResponseWriter, r *http.Request, userID uuid.UUID, next http.Handler) {
	cfg.Logger.InfoLog.Println("Access token expired, attempting refresh...")

	session, err := cfg.DB.GetLatestSessionByID(context.Background(), userID)
	if err != nil || session.IsRevoked || session.ExpiresAt.Before(time.Now()) {
		cfg.Logger.ErrorLog.Println("Invalid or expired refresh token")
		cfg.RedirectToLogin(w, r)
		return
	}

	// Generate JWT token
	newAccessToken, err := cfg.JWTManager.MakeJWT(userID)
	if err != nil {
		cfg.Logger.ErrorLog.Println("Failed to generate new jwt token")
		cfg.RedirectToLogin(w, r)
		return
	}

	newRefreshToken, err := cfg.JWTManager.MakeRefreshToken()
	if err != nil {
		cfg.Logger.ErrorLog.Println("Failed to generate new refresh tokens")
		cfg.RedirectToLogin(w, r)
		return
	}

	params := services.ReplaceRefreshTokenParams{
		Token:     newRefreshToken,
		UserID:    userID,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	_, err = cfg.DB.ReplaceRefreshToken(context.Background(), params)
	if err != nil {
		cfg.Logger.ErrorLog.Println("Failed to store new refresh token")
		cfg.Helper.ServerError(w, err)
		return
	}

	w.Header().Set("Authorization", "Bearer "+newAccessToken)

	next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), cfg.ContextKey["userID"], userID)))
}

func (cfg *Config) RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, cfg.APIRoute+"users/login", http.StatusUnauthorized)
}

func (cfg *Config) LoggingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cfg.Logger.RequestLog.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
			next.ServeHTTP(w, r)
		})
	}
}

func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note: This is split across multiple lines for readability. You don't need to do this in your own code.
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; fontsrc fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

func (cfg *Config) RequireActiveWorkoutSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := ctx.Value(cfg.ContextKey["userID"]).(uuid.UUID)
		fmt.Println("userID: ", userID)
		if !ok {
			cfg.Helper.InfoLog.Println("From RequireActiveWorkoutSession")
			cfg.Helper.ClientError(w, http.StatusUnauthorized)
			return
		}

		sessionData := cfg.SessionSvc.GetWorkoutSession(ctx, userID.String())
		if sessionData == nil {
			cfg.Helper.NotFound(w)
			return
		}

		completed, _ := strconv.ParseBool(sessionData["completed"])
		ctx = context.WithValue(ctx, cfg.ContextKey["sessionCompleted"], completed)
		ctx = context.WithValue(ctx, cfg.ContextKey["userID"], userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *Config) RequireActiveExerciseSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := ctx.Value(cfg.ContextKey["userID"]).(uuid.UUID)
		if !ok {
			cfg.Helper.InfoLog.Println("From RequireActiveExerciseSession")
			cfg.Helper.ClientError(w, http.StatusUnauthorized)
			return
		}

		session := cfg.SessionSvc.GetExerciseSession(ctx, userID.String())
		if session == nil {
			cfg.Helper.NotFound(w)
			return
		}

		// Store only the session ID (or relevant details) in context
		ctx = context.WithValue(ctx, cfg.ContextKey["exerciseSessionID"], session["id"])
		ctx = context.WithValue(ctx, cfg.ContextKey["userID"], userID) // 🔥 Keep userID

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
