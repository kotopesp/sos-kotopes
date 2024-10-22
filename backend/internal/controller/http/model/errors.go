package model

import (
	"errors"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/validator"
)

const MaxFileSize = 10 * 1024 * 1024

var (
	ErrNotUniqueUsername = func(username string) validator.ResponseError {
		return validator.ResponseError{
			FailedField: "Username",
			Tag:         "unique",
			Value:       username,
		}
	}

	ErrInvalidTokenID       = errors.New("invalid token id")
	ErrInvalidTokenUsername = errors.New("invalid token username")
	ErrFailedToParseToken   = errors.New("failed to parse token")
	ErrOAuthStateMismatch   = errors.New("states do not match")
	ErrValidationFailed     = errors.New("validation failed")
	ErrInvalidBody          = errors.New("invalid body")
	ErrInvalidPhotoSize     = errors.New("photo is too large")
	ErrInvalidExtension     = errors.New("invalid photo extension")
	ErrPhotoNotFound        = errors.New("photo not found")
)
