package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cheezecakee/fitrkr/internal/models"
)

func (api *Api) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ClientError(w, http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := api.Helper.HashPassword(req.PasswordHash)
	if err != nil {
		api.ServerError(w, err)
		return
	}

	// Insert user into DB
	newUser, err := api.UserSvc.Register(ctx, &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
	})
	if err != nil {
		api.ServerError(w, err)
		return
	}

	// Add custom logger and err later
	log.Println("User created succesfully!")
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (api *Api) UpdateUser(w http.ResponseWriter, r *http.Request) {
}

func (api *Api) DeleteUser(w http.ResponseWriter, r *http.Request) {
}
