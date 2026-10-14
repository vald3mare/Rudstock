package models

import "errors"

var (
	ErrInvalidEmail = errors.New("invalid email")
	ErrWeakPassword = errors.New("weak password")
	ErrEmailTaken   = errors.New("email taken")
	// ErrInvalidCredentials (401)
	// ErrUserInactive (403)
)
