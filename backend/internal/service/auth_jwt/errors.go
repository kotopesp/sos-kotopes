package auth

import "errors"

var (
	ErrUserIsNil         = errors.New("user is nil")
	ErrNotUniqueUsername = errors.New("username must be unique")
)
