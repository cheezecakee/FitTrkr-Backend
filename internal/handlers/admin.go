package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cheezecakee/FitLogr/internal/models"
)

func (cfg *Config) GetVersion(w http.ResponseWriter, r *http.Request) {
	version := "Version: " + cfg.APIRoute[5:len(cfg.APIRoute)-1]
	// Create a response structure
	Version := struct {
		Response string `json:"response"`
	}{
		Response: version,
	}

	// Set response headers and encode the response into JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Ensure 200 OK status code is set
	json.NewEncoder(w).Encode(Version)
}

func (cfg *Config) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := cfg.DB.GetUsers(context.Background())
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	var userList []models.User
	for _, user := range users {
		userList = append(userList, models.User{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Age:       user.Age.Int32,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userList)
}
