package handler

import (
	"context"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/cheezecakee/fitrkr/internal/db/user"
	"github.com/cheezecakee/fitrkr/internal/utils/auth"
)

type AuthMiddleware struct {
	JWTManager auth.JWT
	UserSvc    user.UserService
}

func NewAuthMiddleware(jwtMgr auth.JWT, userSvc user.UserService) *AuthMiddleware {
	return &AuthMiddleware{
		JWTManager: jwtMgr,
		UserSvc:    userSvc,
	}
}

func (m *AuthMiddleware) extractToken(r *http.Request) (string, error) {
	var token string

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		token = strings.TrimPrefix(authHeader, "Bearer ")
		log.Printf("Found Bearer token in Authorization header")
		return token, nil
	}

	cookieNames := []string{"session"}

	for _, cookieName := range cookieNames {
		if cookie, err := r.Cookie(cookieName); err == nil && cookie.Value != "" {
			token = cookie.Value
			log.Printf("Found JWT token in cookie: %s", cookieName)
			return token, nil
		}
	}

	log.Printf("No JWT token found in Authorization header or cookies")
	return "", &AuthError{Message: "no authentication token found"}
}

// AuthError represents authentication errors
type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

func (m *AuthMiddleware) IsAuthenticated() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Log request details for debugging remove in production
			log.Printf("Auth check for %s %s", r.Method, r.URL.Path)
			log.Printf("Request headers: %v", r.Header)

			// Extract token from either Bearer header or cookies
			token, err := m.extractToken(r)
			if err != nil {
				log.Printf("Token extraction failed: %v", err)
				http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
				return
			}

			log.Printf("Extracted token (first 20 chars): %s...",
				func() string {
					if len(token) > 20 {
						return token[:20]
					}
					return token
				}())

			userID, err := m.JWTManager.ValidateJWT(token)
			if err != nil {
				log.Printf("JWT validation failed: %v", err)
				http.Error(w, "unauthorized: invalid token", http.StatusUnauthorized)
				return
			}

			log.Printf("JWT validation successful for user ID: %s", userID)

			user, err := m.UserSvc.GetUserByID(r.Context(), userID)
			if err != nil {
				log.Printf("Failed to get user by ID %s: %v", userID, err)
				http.Error(w, "unauthorized: user not found", http.StatusUnauthorized)
				return
			}

			log.Printf("User found: %s with roles: %v", user.Username, user.Roles)

			ctx := context.WithValue(r.Context(), UserKey, &user)
			ctx = context.WithValue(ctx, UserIDKey, user.ID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (m *AuthMiddleware) RequireAdmin() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("RequireAdmin: Context keys available: %+v", r.Context())

			currentUser, ok := r.Context().Value(UserKey).(*user.User)
			log.Printf("RequireAdmin: userCtx extraction - ok: %v, userCtx: %+v", ok, currentUser)

			if !ok || currentUser == nil {
				log.Printf("RequireAdmin: Failed to get user from context")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			log.Printf("RequireAdmin: User found: %s, Roles: %v, Role type: %T", currentUser.Username, currentUser.Roles, currentUser.Roles)

			if len(currentUser.Roles) == 0 {
				log.Printf("RequireAdmin: User has no roles assigned")
			}

			if slices.Contains(currentUser.Roles, "admin") {
				log.Println("RequireAdmin: Admin access granted")
				next.ServeHTTP(w, r)
				return
			}

			log.Printf("RequireAdmin: User is not admin. Expected 'admin', got roles: %v", currentUser.Roles)
			http.Error(w, "forbidden: admin only", http.StatusForbidden)
		})
	}
}

func (m *AuthMiddleware) RequireAnyRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			currentUser, ok := r.Context().Value(UserKey).(*user.User)
			if !ok || currentUser == nil {
				log.Printf("RequireAnyRole: Failed to get user from context")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			log.Printf("RequireAnyRole: Checking if user %s has any of roles: %v (user has: %v)",
				currentUser.Username, roles, currentUser.Roles)

			for _, requiredRole := range roles {
				if slices.Contains(currentUser.Roles, requiredRole) {
					log.Printf("RequireAnyRole: Access granted - user has role: %s", requiredRole)
					next.ServeHTTP(w, r)
					return
				}
			}

			log.Printf("RequireAnyRole: Access denied - user lacks required roles: %v", roles)
			http.Error(w, "forbidden: insufficient privileges", http.StatusForbidden)
		})
	}
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow multiple origins - add your actual Flutter app URLs
		allowedOrigins := []string{
			"http://localhost:5173", // Your Svelte app
			"http://localhost:3000", // Common Flutter web port
			"http://localhost:8080", // Another common Flutter web port
			"http://10.0.2.2:3000",  // Android emulator
			"http://127.0.0.1:3000", // Alternative localhost
			"https://receiver-consistently-exchange-women.trycloudflare.com", // Your Cloudflare tunnel
			"http://localhost",      // For mobile apps
			"capacitor://localhost", // For Capacitor apps
			"ionic://localhost",     // For Ionic apps
		}

		origin := r.Header.Get("Origin")

		// Debug logging - remove in production
		log.Printf("Request Origin: %s", origin)
		log.Printf("Request Method: %s", r.Method)
		log.Printf("Request Headers: %v", r.Header)

		if slices.Contains(allowedOrigins, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Printf("Origin allowed: %s", origin)
		} else if origin == "" {
			// For mobile apps that don't send Origin header, allow them
			w.Header().Set("Access-Control-Allow-Origin", "*")
			log.Printf("No origin header, allowing all")
		} else {
			log.Printf("Origin not allowed: %s", origin)
			// Might want to still set some CORS headers for the error response
			w.Header().Set("Access-Control-Allow-Origin", "null")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400") // Cache preflight for 24 hours

		if r.Method == http.MethodOptions {
			log.Printf("Handling OPTIONS preflight request from: %s", origin)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
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
