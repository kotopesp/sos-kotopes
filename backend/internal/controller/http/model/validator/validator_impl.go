package validator

import (
	"regexp"

	validatorPkg "github.com/go-playground/validator/v10"
)

var (
	upperase = regexp.MustCompile(`[A-Z]`).MatchString
	digit    = regexp.MustCompile(`\d`).MatchString
)

// custom validation tags
func customValidationOptions(validator *validatorPkg.Validate) {
	_ = validator.RegisterValidation("contains_digit", func(fl validatorPkg.FieldLevel) bool {
		return digit(fl.Field().String())
	})
	_ = validator.RegisterValidation("contains_uppercase", func(fl validatorPkg.FieldLevel) bool {
		return upperase(fl.Field().String())
	})
}

func newValidator(validator *validatorPkg.Validate) formValidatorAPI {
	customValidationOptions(validator)
	return &formValidator{
		validator: validator,
	}
}

func (v *formValidator) validate(data interface{}) []ResponseError {
	validationErrors := []ResponseError{}

	errs := v.validator.Struct(data)

	if errs != nil {
		pkgValidationErrors, ok := errs.(validatorPkg.ValidationErrors)
		if !ok {
			return []ResponseError{}
		}
		for _, err := range pkgValidationErrors {
			var elem ResponseError

			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Param = err.Param()
			elem.Value = err.Value()

			validationErrors = append(validationErrors, elem)
		}
	}
	return validationErrors
}

func Validate(data interface{}) []ResponseError {
	return v.validate(data)
}
