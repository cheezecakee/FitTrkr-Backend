package errors

import "errors"

var (
	// User-related errors
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailExists        = errors.New("email alrady exists")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrUsernameTaken      = errors.New("username alrady taken")
	ErrInvalidUsername    = errors.New("invalid username format")
	ErrWeakPassword       = errors.New("password does not meet complexity requirements")

	// Authentication & Authorization errors
	ErrUnauthorized   = errors.New("unauthorized access")
	ErrForbiden       = errors.New("forbiden action")
	ErrInternalServer = errors.New("internal server error")

	// General errors
	ErrBadRequest = errors.New("bad request")
	ErrConflict   = errors.New("conflict detected")
)
