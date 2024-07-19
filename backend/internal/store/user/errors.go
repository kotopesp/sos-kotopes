package user

import "errors"

var (
	ErrNotUniqueUsername = errors.New("username must be unique")
	ErrNoSuchUser        = errors.New("user does not exist")
)
