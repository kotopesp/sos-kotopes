package auth

import "errors"

var (
	ErrNotUniqueUsername = errors.New("username must be unique")
)
