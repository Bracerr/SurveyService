package apperrors

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenUsed          = errors.New("token already used")
	ErrValidationFailed   = errors.New("validation failed")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidEmail       = errors.New("invalid email")
)
