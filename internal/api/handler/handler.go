package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

type ContextKey string

const UserIDKey ContextKey = "userID"

// Errors
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	log.Println(trace)
	// logger.Log.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func ClientError(w http.ResponseWriter, status int) {
	// logger.Log.InfoLog.Printf("Client error: %d", status)
	http.Error(w, http.StatusText(status), status)
}

func NotFound(w http.ResponseWriter) {
	ClientError(w, http.StatusNotFound)
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// Response sends a JSON response with the given status code and message.
func Response(w http.ResponseWriter, status int, message any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}
