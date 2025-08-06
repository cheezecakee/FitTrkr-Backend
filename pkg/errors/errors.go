// Package errors provides custom error types and helpers for FitTrkr.
package errors

import "errors"

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	// User-related errors
	ErrUserNotFound       = errors.New("user not found")
	ErrDuplicateUsername  = errors.New("username already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailExists        = errors.New("email already exists")
	ErrDuplicateEmail     = errors.New("email already exists")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrInvalidUsername    = errors.New("invalid username format")
	ErrWeakPassword       = errors.New("password does not meet complexity requirements")

	// Authentication & Authorization errors
	ErrUnauthorized   = errors.New("unauthorized access")
	ErrForbidden      = errors.New("forbidden action")
	ErrInternalServer = errors.New("internal server error")

	// General errors
	ErrBadRequest = errors.New("bad request")
	ErrConflict   = errors.New("conflict detected")
)
