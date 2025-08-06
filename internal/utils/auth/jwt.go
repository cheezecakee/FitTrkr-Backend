// Package auth provides JWT authentication utilities.
package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWT interface {
	MakeJWT(userID uuid.UUID, roles []string) (string, error)
	ValidateJWT(tokenString string) (uuid.UUID, error)
}

type UserClaims struct {
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	SecretKey []byte
	ExpiresIn time.Duration
	Claims    UserClaims
}

func NewJWTManager(secretKey []byte, expiresIn time.Duration) JWT {
	return &JWTManager{SecretKey: secretKey, ExpiresIn: expiresIn}
}

func (j *JWTManager) MakeJWT(userID uuid.UUID, roles []string) (string, error) {
	claims := &UserClaims{
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "fitrkr",
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		// Add custom logger and err later
		return "", err
	}
	return ss, nil
}

func (j *JWTManager) ValidateJWT(tokenString string) (uuid.UUID, error) {
	var userClaims UserClaims

	token, err := jwt.ParseWithClaims(tokenString, &userClaims, func(token *jwt.Token) (any, error) {
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, nil
	}

	userID, err := userClaims.GetSubject()
	if err != nil {
		// Add custom logger and err later
		return uuid.Nil, err
	}

	// Add custom logger and err later
	return uuid.MustParse(userID), err
}
