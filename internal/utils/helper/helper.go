// Package helper provides utility functions for FitTrkr.
package helper

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"

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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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

func CleanString(input string) string {
	return strings.TrimSpace(input)
}

func IsValidLength(input string, maxLength int) bool {
	return len(input) <= maxLength
}

func CleanAndCheckLength(input string, maxLength int) (string, error) {
	cleaned := strings.TrimSpace(input)
	if len(cleaned) > maxLength {
		return "", fmt.Errorf("input exceeds maximum length of %d", maxLength)
	}
	return cleaned, nil
}
