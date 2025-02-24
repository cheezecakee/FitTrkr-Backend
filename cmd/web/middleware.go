package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	jwt.RegisteredClaims
}

func (apiCfg *ApiConfig) MakeJWT(userID uuid.UUID, expiresIn time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer:    "fitlogr",
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(apiCfg.JWTSecret))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func (apiCfg *ApiConfig) ValidateJWT(tokenString string) (uuid.UUID, error) {
	var userClaims UserClaims

	token, err := jwt.ParseWithClaims(tokenString, &userClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(apiCfg.JWTSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	// fmt.Printf("Signature: %v\n", token.Signature)
	// fmt.Printf("Method: %v\n", token.Method.Alg())

	if !token.Valid {
		return uuid.Nil, nil
	}
	userID, err := userClaims.GetSubject()
	if err != nil {
		return uuid.Nil, errors.New("invalid token or claims")
	}

	return uuid.MustParse(userID), err
}

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", errors.New("error generating a new token")
	}

	encodedToken := hex.EncodeToString(token)
	return encodedToken, nil
}

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note: This is split across multiple lines for readability. You don't need to do this in your own code.
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; fontsrc fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}
