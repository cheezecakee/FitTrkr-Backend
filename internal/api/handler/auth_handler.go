package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/cheezecakee/fitrkr/internal/db/user"
)

type AuthHandler struct {
	svc user.UserService
}

func NewAuthHandler(svc user.UserService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Login logs in a user
// @Summary Log in a user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body user.LoginRequest true "Login payload"
// @Success 200 {object} map[string]string "JWT token"
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ClientError(w, http.StatusBadRequest)
		return
	}

	token, err := h.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		ClientError(w, http.StatusUnauthorized)
		return
	}

	log.Println("User logged in successfully!")
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour), // Adjust as needed
	})

	w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ClientError(w, http.StatusUnauthorized)
		return
	}

	log.Println("User logged out successfully!", userID)
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Placeholder: Requires refresh token table/repo/service
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ClientError(w, http.StatusBadRequest)
		return
	}

	// TODO: Validate refresh token via UserSvc, generate new JWT
	NotFound(w) // Temporary
}

func (h *AuthHandler) RevokeToken(w http.ResponseWriter, r *http.Request) {
	// Placeholder: Requires refresh token table/repo/service
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ClientError(w, http.StatusBadRequest)
		return
	}

	// TODO: Revoke refresh token via UserSvc
	NotFound(w) // Temporary
}
