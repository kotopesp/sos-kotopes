package validator

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/keeper"
	"github.com/kotopesp/sos-kotopes/pkg/logger"

	validatorPkg "github.com/go-playground/validator/v10"
)

var (
	uppercase      = regexp.MustCompile(`[A-Z]`).MatchString
	digit          = regexp.MustCompile(`\d`).MatchString
	alphaNum       = regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString
	numericNatural = regexp.MustCompile(`^[1-9]\d*$`).MatchString
	notBlank       = regexp.MustCompile(`.*\S+.*`).MatchString
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
	err = validator.RegisterValidation("numeric_natural", func(fl validatorPkg.FieldLevel) bool {
		return numericNatural(fl.Field().String())
	})
	if err != nil {
		logger.Log().Fatal(ctx, err.Error())
	}
	err = validator.RegisterValidation("notblank", func(fl validatorPkg.FieldLevel) bool {
		return notBlank(fl.Field().String())
	})
	if err != nil {
		logger.Log().Fatal(ctx, err.Error())
	}
	err = validator.RegisterValidation("sort_keeper", func(fl validatorPkg.FieldLevel) bool {
		sort := fl.Field().String()
		if sort == "" {
			return true
		}

		sortBy, sortOrder := keeper.ParseSort(sort)
		if sortBy == "" && sortOrder == "" {
			return false
		}

		validFields := map[string]bool{
			"avg_grade":  true,
			"price":      true,
			"created_at": true,
		}
		if _, ok := validFields[sortBy]; !ok {
			return false
		}

		validOrders := map[string]bool{
			"asc":  true,
			"desc": true,
		}
		if _, ok := validOrders[strings.ToLower(sortOrder)]; !ok {
			return false
		}

		return true
	})
	if err != nil {
		logger.Log().Fatal(ctx, err.Error())
	}
}

func New(ctx context.Context, validator *validatorPkg.Validate) FormValidatorService {
	logger.Log().Info(context.Background(), "validator created")
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
