package http

import (
	"errors"
	"fmt"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/user"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/validator"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

// token helpers: getting info from token
func getIDFromToken(ctx *fiber.Ctx) (int, error) {
	idItem := getPayloadItem(ctx, "id")
	idFloat, ok := idItem.(float64)
	if !ok {
		return 0, errors.New("invalid id")
	}
	return int(idFloat), nil
}

func getUsernameFromToken(ctx *fiber.Ctx) (string, error) {
	usernameItem := getPayloadItem(ctx, "username")
	username, ok := usernameItem.(string)
	if !ok {
		return "", errors.New("invalid username")
	}
	return username, nil
}

// validation helpers
func parseAndValidate(ctx *fiber.Ctx, formValidator validator.FormValidatorService, apiUser user.User) error {
	if err := ctx.BodyParser(&apiUser); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}
	errs := formValidator.Validate(apiUser)
	if len(errs) > 0 {
		logger.Log().Error(ctx.UserContext(), fmt.Sprintf("%v", errs))
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ErrorResponse(fiber.Map{
			"validation_errors": errs,
		}))
	}
	return nil
}
