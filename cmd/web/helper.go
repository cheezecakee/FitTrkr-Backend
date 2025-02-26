package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Helper function to convert string to UUID
func StrToUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id) // Directly return result and error
}

// Helper function to check if the user exists.
func UserExists(data *Data, userID string) bool {
	for _, user := range data.Users {
		if user.ID == userID {
			return true
		}
	}
	return false
}

// Helper function to validate exercise IDs.
func ValidateExercise(data *Data, exerciseID string) bool {
	for _, ex := range data.Exercises {
		if ex.ID == exerciseID {
			return true
		}
	}
	return false
}

// Helper function to get the exercise name by ID.
func GetExerciseNameByID(data *Data, exerciseID string) string {
	for _, ex := range data.Exercises {
		if ex.ID == exerciseID {
			return ex.Name
		}
	}
	return ""
}

func StringToNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{Valid: false}
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return err
	}
	return nil
}

func GetBearerToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")
	if token == "" {
		return "", errors.New("missing authorization header")
	}
	token = token[len("Bearer "):]
	return token, nil
}
