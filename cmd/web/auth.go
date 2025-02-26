package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	jwt.RegisteredClaims
}

func (apiCfg *ApiConfig) MakeJWT(userID uuid.UUID, expiresIn time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer:    "fittrackr",
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
