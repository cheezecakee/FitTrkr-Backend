package handler

import (
	"encoding/json"
	"net/http"
)

// Main sessions page
// Directs the user either to sessions/workouts or sessions/start/workouts/{id}
func (cfg *Config) GetLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Workout Summary"})
}
