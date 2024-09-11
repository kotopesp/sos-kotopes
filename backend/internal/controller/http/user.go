package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/user"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

// @Summary		Get user by id
// @Tags			user
// @Description	Get user by id
// @ID				get-user
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"User ID"
// @Success		200	{object}	model.Response{data=user.ResponseUser}
// @Failure		400	{object}	model.Response
// @Failure		404	{object}	model.Response
// @Failure		500	{object}	model.Response
// @Security		ApiKeyAuthBasic
// @Router			/users/{id} [get]
func (r *Router) getUser(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}
	currentUser, err := r.userService.GetUser(ctx.UserContext(), id)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrNoSuchUser):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		default:
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
		}
	}
	responseUser := user.ToResponseUser(&currentUser)
	return ctx.Status(fiber.StatusOK).JSON(responseUser)
}

// @Summary		Update user
// @Tags			user
// @Description	Update user
// @ID				update-user
// @Accept			json
// @Produce		json
// @Param			username	formData	string	false	"Username"
// @Param			firstname	formData	string	false	"First name"
// @Param			lastname	formData	string	false	"Last name"
// @Param			description	formData	string	false	"Description"
// @Param			photo		formData	file	false	"Photo"
// @Param			password	formData	string	false	"Password"
// @Success		200			{object}	model.Response{data=user.ResponseUser}
// @Failure		400			{object}	model.Response
// @Failure		404			{object}	model.Response
// @Failure		422			{object}	model.Response{data=validator.Response}
// @Failure		500			{object}	model.Response
// @Router			/users [patch]
func (r *Router) updateUser(ctx *fiber.Ctx) error {
	id, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	var update user.UpdateUser
	update.Photo, err = openAndValidatePhoto(ctx)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidPhotoSize):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
		case errors.Is(err, model.ErrInvalidExtension):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
		default:
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
		}
	}

	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &update)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}
	coreUpdate := update.ToCoreUpdateUser()

	updatedUser, err := r.userService.UpdateUser(ctx.UserContext(), id, coreUpdate)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrNoSuchUser):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		case errors.Is(err, core.ErrEmptyUpdateRequest):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
		default:
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
		}
	}
	responseUser := user.ToResponseUser(&updatedUser)

	return ctx.Status(fiber.StatusOK).JSON(responseUser)
}
