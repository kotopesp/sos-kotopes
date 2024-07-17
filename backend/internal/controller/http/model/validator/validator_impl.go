package validator

import (
	"regexp"

	validatorPkg "github.com/go-playground/validator/v10"
)

var (
	upperase = regexp.MustCompile(`[A-Z]`).MatchString
	digit    = regexp.MustCompile(`[0-9]`).MatchString
)

// if you need to create custom validation tags
func customValidationOptions(validator *validatorPkg.Validate) {
	validator.RegisterValidation("containsDigit", func(fl validatorPkg.FieldLevel) bool {
		return digit(fl.Field().String())
	})
	validator.RegisterValidation("containsUppercase", func(fl validatorPkg.FieldLevel) bool {
		return upperase(fl.Field().String())
	})
}

func new(validator *validatorPkg.Validate) formValidatorAPI {
	customValidationOptions(validator)
	return &formValidator{
		validator: validator,
	}
}

func (v *formValidator) validate(data interface{}) []ErrorResponse {
	validationErrors := []ErrorResponse{}

	errs := v.validator.Struct(data)

	if errs != nil {
		for _, err := range errs.(validatorPkg.ValidationErrors) {
			var elem ErrorResponse

			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Param = err.Param()
			elem.Value = err.Value()

			validationErrors = append(validationErrors, elem)
		}
	}
	return validationErrors
}

func Validate(data interface{}) []ErrorResponse {
	return v.validate(data)
}
