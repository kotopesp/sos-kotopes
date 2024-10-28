package validator

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/kotopesp/sos-kotopes/pkg/logger"

	validatorPkg "github.com/go-playground/validator/v10"
)

var (
	uppercase = regexp.MustCompile(`[A-Z]`).MatchString
	digit     = regexp.MustCompile(`\d`).MatchString
	alphaNum  = regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString
)

// custom validator tags
func customValidationOptions(ctx context.Context, validator *validatorPkg.Validate) {
	err := validator.RegisterValidation("contains_digit", func(fl validatorPkg.FieldLevel) bool {
		return digit(fl.Field().String())
	})
	if err != nil {
		logger.Log().Fatal(ctx, err.Error())
	}
	err = validator.RegisterValidation("contains_uppercase", func(fl validatorPkg.FieldLevel) bool {
		return uppercase(fl.Field().String())
	})
	if err != nil {
		logger.Log().Fatal(ctx, err.Error())
	}
	err = validator.RegisterValidation("no_specials", func(fl validatorPkg.FieldLevel) bool {
		return alphaNum(fl.Field().String())
	})
	if err != nil {
		logger.Log().Fatal(ctx, err.Error())
	}
}

func New(ctx context.Context, validator *validatorPkg.Validate) FormValidatorService {
	logger.Log().Info(ctx, "validator created")

	customValidationOptions(ctx, validator)
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

func (err *ResponseError) Error() string {
	return fmt.Sprintf("FailedField: %s | Tag: %s | Param: %s | Value: %v", err.FailedField, err.Tag, err.Param, err.Value)
}

func NewResponse(validationErrors []ResponseError, message *string) Response {
	return Response{
		ValidationErrors: validationErrors,
		Message:          message,
	}
}
