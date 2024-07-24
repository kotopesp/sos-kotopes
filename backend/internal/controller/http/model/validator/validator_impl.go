package validator

import (
	"context"
	"errors"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/logger"
	"regexp"

	validatorPkg "github.com/go-playground/validator/v10"
)

var (
	uppercase = regexp.MustCompile(`[A-Z]`).MatchString
	digit     = regexp.MustCompile(`\d`).MatchString
)

// custom validator tags
func customValidationOptions(validator *validatorPkg.Validate) {
	_ = validator.RegisterValidation("contains_digit", func(fl validatorPkg.FieldLevel) bool {
		return digit(fl.Field().String())
	})
	_ = validator.RegisterValidation("contains_uppercase", func(fl validatorPkg.FieldLevel) bool {
		return uppercase(fl.Field().String())
	})
}

func New(validator *validatorPkg.Validate) FormValidatorService {
	logger.Log().Info(context.Background(), "validator created")
	customValidationOptions(validator)
	return &formValidator{
		validator: validator,
	}
}

func (v *formValidator) Validate(data interface{}) []ResponseError {
	var validationErrors []ResponseError

	errs := v.validator.Struct(data)

	if errs != nil {
		var pkgValidationErrors validatorPkg.ValidationErrors
		ok := errors.As(errs, &pkgValidationErrors)
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
