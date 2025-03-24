package helper

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Helper struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

func NewHelper(infoLog, errorLog *log.Logger) *Helper {
	return &Helper{
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}
}

// Strings and UUIDs
func (h *Helper) StrToUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id) // Directly return result and error
}

func (h *Helper) StringToNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{Valid: false}
}

// Password
func (h *Helper) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (h *Helper) CheckPasswordHash(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return err
	}
	return nil
}

// Tokens
func (h *Helper) GetBearerToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")
	if token == "" {
		return "", errors.New("missing authorization header")
	}
	token = token[len("Bearer "):]
	return token, nil
}

// Middlware
func (h *Helper) MiddlewareWrapper(next http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		next = middleware(next)
	}
	return next
}

// Errors
func (h *Helper) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	h.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (h *Helper) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (h *Helper) NotFound(w http.ResponseWriter) {
	h.ClientError(w, http.StatusNotFound)
}
