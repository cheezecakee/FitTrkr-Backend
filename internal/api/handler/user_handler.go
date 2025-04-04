package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/cheezecakee/fitrkr/internal/models"
	"github.com/cheezecakee/fitrkr/internal/service"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Username  string `json:"username"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ClientError(w, http.StatusBadRequest)
		return
	}

	// Insert user into DB
	newUser, err := h.svc.Register(ctx, &models.User{
		Username:     req.Username,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: req.Password,
	})
	if err != nil {
		switch err {
		case service.ErrDuplicateEmail, service.ErrDuplicateUsername:
			log.Println("client error: ", err)
			ClientError(w, http.StatusConflict)
		default:
			ServerError(w, err)
			log.Println("server error:", err)
		}
		return
	}

	// Add custom logger and err later
	log.Println("User created successfully!")
	response := models.UserResponse{
		ID:        newUser.ID,
		Username:  newUser.Username,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		IsPremium: newUser.IsPremium,
	}

	// Return the created user (excluding password)
	Response(w, http.StatusCreated, response)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ClientError(w, http.StatusUnauthorized)
		return
	}

	var req models.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ClientError(w, http.StatusBadRequest)
		return
	}

	// Handle password hashing only if provided
	// Move this to the service?
	var passwordHash string
	if req.PasswordHash != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			ServerError(w, err)
			return
		}

		hashedPasswordStr := string(hashedPassword)
		passwordHash = hashedPasswordStr
	}

	// Update user
	updatedUser, err := h.svc.Update(ctx, &models.User{
		ID:           userID,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PasswordHash: passwordHash,
	})
	if err != nil {
		ServerError(w, err)
		return
	}

	log.Println("User account updated successfully!")
	response := models.UserResponse{
		ID:        updatedUser.ID,
		Username:  updatedUser.Username,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		IsPremium: updatedUser.IsPremium,
	}

	// Return updated user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ClientError(w, http.StatusUnauthorized)
		return
	}

	err := h.svc.Delete(ctx, userID)
	if err != nil {
		ServerError(w, err)
		return
	}

	log.Println("User deleted successfully!")
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
