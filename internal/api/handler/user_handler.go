package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/cheezecakee/fitrkr/internal/db/user"
	"github.com/cheezecakee/fitrkr/pkg/errors"
)

type UserHandler struct {
	svc user.UserService
}

func NewUserHandler(svc user.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// CreateUser creates a new user
// @Summary Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param request body user.CreateUserRequest true "User creation payload"
// @Success 201 {object} user.UserResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
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
	newUser, err := h.svc.Register(r.Context(), user.User{
		Username:     req.Username,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: req.Password,
	})
	if err != nil {
		switch err {
		case errors.ErrDuplicateEmail, errors.ErrDuplicateUsername:
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
	response := user.UserResponse{
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

// UpdateUser updates an authenticated user
// @Summary Update user account
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body user.UserRequest true "User update payload"
// @Success 200 {object} user.UserRequest
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /api/v1/users [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ClientError(w, http.StatusUnauthorized)
		return
	}

	var req user.User

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
	updatedUser, err := h.svc.Update(ctx, user.User{
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
	response := user.UserResponse{
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

// GetCurrentUser returns the current authenticated user's information
// @Summary Get current user
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} user.UserResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/v1/users/me [get]
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ClientError(w, http.StatusUnauthorized)
		return
	}

	userData, err := h.svc.GetUserByID(r.Context(), userID)
	if err != nil {
		switch err {
		case errors.ErrUserNotFound:
			ClientError(w, http.StatusNotFound)
		default:
			ServerError(w, err)
		}
		return
	}

	response := user.UserResponse{
		ID:        userData.ID,
		Username:  userData.Username,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		Email:     userData.Email,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
		IsPremium: userData.IsPremium,
		Roles:     userData.Roles,
	}

	Response(w, http.StatusOK, response)
}

// DeleteUser deletes an authenticated user
// @Summary Delete user account
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 204 {string} string "No Content"
// @Failure 401 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/v1/users [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ClientError(w, http.StatusUnauthorized)
		return
	}

	err := h.svc.Delete(ctx, userID)
	if err != nil {
		switch err {
		case errors.ErrUserNotFound:
			ClientError(w, http.StatusNotFound)
		default:
			ServerError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
