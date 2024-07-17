package validator

import (
	"fmt"

	validatorPkg "github.com/go-playground/validator/v10"
)

type (
	ErrorResponse struct {
		FailedField string
		Tag         string
		Param       string
		Value       interface{}
	}

	formValidator struct {
		validator *validatorPkg.Validate
	}

	formValidatorAPI interface {
		validate(data interface{}) []ErrorResponse
	}
)

func (err *ErrorResponse) Error() string {
	return fmt.Sprintf("FailedField: %s | Tag: %s | Param: %s | Value: %v", err.FailedField, err.Tag, err.Param, err.Value)
}
