package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

func (j *JWTManager) MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", errors.New("error generating a new token")
	}

	encodedToken := hex.EncodeToString(token)
	return encodedToken, nil
}
