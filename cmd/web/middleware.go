package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note: This is split across multiple lines for readability. You don't need to do this in your own code.
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; fontsrc fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

type contextKey string

const userIDKey contextKey = "userID"

func (apiCfg *ApiConfig) isAuthenticated(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent caching of responses
		w.Header().Add("Cache-Control", "no-store")

		// Get the token from request headers
		refreshToken, err := GetBearerToken(r.Header)
		if err != nil {
			log.Printf("Error getting bearer token: %s", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Printf("RefreshToken: %s\n", refreshToken)
		// Get session from the database
		session, err := apiCfg.DB.GetSession(context.Background(), refreshToken)
		if err != nil {
			log.Printf("Error retrieving refresh token: %s", err)
			http.Redirect(w, r, "/api/user/login", http.StatusNotFound)
			return
		}

		// Check if the token is expired
		if time.Now().After(session.ExpiresAt) {
			log.Printf("Refresh token has expired.")
			http.Redirect(w, r, "/api/user/login", http.StatusNotFound)
			return
		}

		// Check if the token is revoked
		if session.IsRevoked {
			log.Printf("Refresh token has been revoked.")
			http.Redirect(w, r, "/api/user/login", http.StatusNotFound)
			return
		}

		// Attach user ID to the request context
		ctx := context.WithValue(r.Context(), userIDKey, session.UserID)

		// Proceed to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (apiCfg *ApiConfig) ValidateSession(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := GetBearerToken(r.Header)
		if err != nil {
			log.Printf("%s", err)
			http.Error(w, "Failed to get bearer token", http.StatusInternalServerError)
			return
		}

		userID, err := apiCfg.ValidateJWT(token)
		if err != nil {
			log.Printf("%s", err)
			return
		}

		session, err := apiCfg.DB.GetLatestSessionByID(context.Background(), userID)
		if err != nil {
			log.Printf("Error retrieving sessions: %s", err)
			http.Redirect(w, r, "/api/user/login", http.StatusNotFound)
			return
		}

		if session.IsRevoked {
			log.Printf("Refresh token has been revoked.")
			http.Redirect(w, r, "/api/user/login", http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
