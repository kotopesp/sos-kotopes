package model

import "gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/validator"

var (
	ErrNotUniqueUsername = validator.ErrorResponse{
		FailedField: "Username",
		Tag:         "unique",
	}
)
