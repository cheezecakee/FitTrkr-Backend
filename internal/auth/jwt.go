package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	jwt.RegisteredClaims
}

type JWTManager struct {
	SecretKey []byte
	ExpiresIn time.Duration
}

func NewJWTManager(secretKey []byte, expiresIn time.Duration) *JWTManager {
	return &JWTManager{SecretKey: secretKey, ExpiresIn: expiresIn}
}

func (j *JWTManager) MakeJWT(userID uuid.UUID) (string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer:    "fittrackr",
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
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
