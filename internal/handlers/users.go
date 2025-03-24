package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/cheezecakee/FitLogr/internal/models"
	"github.com/cheezecakee/FitLogr/internal/services"
)

func (cfg *Config) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.Helper.ClientError(w, http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := cfg.Helper.HashPassword(req.Password)
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	// Insert user into DB
	newUser, err := cfg.DB.RegisterUser(ctx, services.RegisterUserParams{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Age:          sql.NullInt32{Int32: req.Age, Valid: true},
	})
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	cfg.Logger.InfoLog.Println("User created succesfully!")

	// Return the created user (excluding password)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUser)
}

// Login User
func (cfg *Config) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.Helper.ClientError(w, http.StatusBadRequest)
		return
	}

	// Fetch user from DB
	user, err := cfg.DB.GetUserByEmail(ctx, req.Email)
	if err != nil {
		cfg.Helper.ClientError(w, http.StatusUnauthorized)
		return
	}

	// Check password
	if err := cfg.Helper.CheckPasswordHash(user.PasswordHash, req.Password); err != nil {
		cfg.Helper.ClientError(w, http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	accessToken, err := cfg.JWTManager.MakeJWT(user.ID)
	log.Printf("Generated Token: %s", accessToken)
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	refreshToken, err := cfg.JWTManager.MakeRefreshToken()
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	// Expires in 30 days
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	params := services.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}

	_, err = cfg.DB.CreateRefreshToken(context.Background(), params)
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	cfg.Logger.InfoLog.Println("User login successful!")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.User{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (cfg *Config) LogoutUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(cfg.ContextKey["userID"]).(uuid.UUID)
	if !ok {
		cfg.Helper.ClientError(w, http.StatusUnauthorized)
		return
	}

	err := cfg.DB.DeleteSession(context.Background(), userID)
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	cfg.Logger.InfoLog.Println("User logged out successfully!")
	w.WriteHeader(http.StatusNoContent)
}

func (cfg *Config) EditUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(cfg.ContextKey["userID"]).(uuid.UUID)
	if !ok {
		cfg.Helper.ClientError(w, http.StatusUnauthorized)
		return
	}

	var req models.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.Helper.ClientError(w, http.StatusBadRequest)
		return
	}

	// Handle password hashing only if provided
	var passwordHash string
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			cfg.Helper.ServerError(w, err)
			return
		}

		hashedPasswordStr := string(hashedPassword)
		passwordHash = hashedPasswordStr
	}

	// Update user
	updatedUser, err := cfg.DB.EditUser(ctx, services.EditUserParams{
		ID:           userID,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Age:          sql.NullInt32{Int32: req.Age, Valid: true},
		PasswordHash: passwordHash,
	})
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	cfg.Logger.InfoLog.Println("User account updated successfully!")

	// Return updated user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

func (cfg *Config) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(cfg.ContextKey["userID"]).(uuid.UUID)
	if !ok {
		cfg.Helper.ClientError(w, http.StatusUnauthorized)
		return
	}

	err := cfg.DB.DeleteUser(ctx, userID)
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	cfg.Logger.InfoLog.Println("User deleted successfully!")
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
