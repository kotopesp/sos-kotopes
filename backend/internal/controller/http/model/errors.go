package model

import (
	"errors"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/validator"
)

var (
	ErrNotUniqueUsername = func(username string) validator.ResponseError {
		return validator.ResponseError{
			FailedField: "Username",
			Tag:         "unique",
			Value:       username,
		}
	}

	ErrInvalidTokenID = errors.New("invalid token id")

	ErrInvalidTokenUsername = errors.New("invalid token username")

	ErrOAuthStateMismatch = errors.New("states do not match")
)
