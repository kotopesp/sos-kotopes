package validator

import (
	validatorPkg "github.com/go-playground/validator/v10"
)

type (
	Response struct {
		ValidationErrors []ResponseError `json:"validation_errors,omitempty"`
		Message          *string         `json:"message,omitempty"`
	}

	ResponseError struct {
		FailedField string      `json:"failed_field" example:"username"`
		Tag         string      `json:"tag" example:"required"`
		Param       string      `json:"param" example:""`
		Value       interface{} `json:"value"`
	}

	formValidator struct {
		validator *validatorPkg.Validate
	}

	FormValidatorService interface {
		Validate(data interface{}) []ResponseError
	}
)
