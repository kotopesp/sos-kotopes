package model

import "gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/validator"

var (
	ErrNotUniqueUsername = func(username string) validator.ResponseError {
		return validator.ResponseError{
			FailedField: "Username",
			Tag:         "unique",
			Value:       username,
		}
	}
)
