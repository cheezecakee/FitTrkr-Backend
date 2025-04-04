package helper

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Clamp(value, min, max int) int {
	switch {
	case value < min:
		return min
	case value > max:
		return max
	default:
		return value
	}
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		// Add custom logger and err later
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		// Add custom logger and err later
		return err
	}
	return nil
}

func GetBearerToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")
	if token == "" {
		// Add custom logger and err later
		return "", errors.New("error getting Bearer Token")
	}
	token = token[len("Bearer "):]
	return token, nil
}

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		// Add custom logger and err later
		return "", err
	}

	encodedToken := hex.EncodeToString(token)
	return encodedToken, nil
}
