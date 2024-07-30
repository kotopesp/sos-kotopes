package http

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/keeper"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/user"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/validator"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

// token helpers: getting info from token
func getIDFromToken(ctx *fiber.Ctx) (id int, err error) {
	idItem := getPayloadItem(ctx, "id")
	idFloat, ok := idItem.(float64)
	if !ok {
		return 0, model.ErrInvalidTokenID
	}
	return int(idFloat), nil
}

func getUsernameFromToken(ctx *fiber.Ctx) (username string, err error) {
	usernameItem := getPayloadItem(ctx, "username")
	username, ok := usernameItem.(string)
	if !ok {
		return "", model.ErrInvalidTokenUsername
	}
	return username, nil
}

// validation helpers
func parseAndValidate(ctx *fiber.Ctx, formValidator validator.FormValidatorService, apiUser *user.User) (fiberError, parseOrValidationError error) {
	if err := ctx.BodyParser(apiUser); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error())), err
	}
	errs := formValidator.Validate(apiUser)
	if len(errs) > 0 {
		logger.Log().Info(ctx.UserContext(), fmt.Sprintf("%v", errs))
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ErrorResponse(fiber.Map{
			"validation_errors": errs,
		})), model.ErrValidationFailed
	}
	return nil, nil
}

func parseAndValidateNewKeeper(ctx *fiber.Ctx, formValidator validator.FormValidatorService, entity *keeper.KeepersCreate) (fiberError, parseOrValidationError error) {
	if err := ctx.BodyParser(entity); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error())), err
	}
	errs := formValidator.Validate(entity)
	if len(errs) > 0 {
		logger.Log().Info(ctx.UserContext(), fmt.Sprintf("%v", errs))
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ErrorResponse(fiber.Map{
			"validation_errors": errs,
		})), model.ErrValidationFailed
	}
	return nil, nil
}

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}
