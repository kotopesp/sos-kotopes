package http

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/user"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

func (r *Router) AddUserToFavourites(ctx *fiber.Ctx) error {
	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	favouriteUserID, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	addedUser, err := r.userService.AddUserToFavourite(ctx.UserContext(), favouriteUserID, userID)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrNoSuchUser):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		case errors.Is(err, core.ErrCantAddYourselfIntoFavourites):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
		case errors.Is(err, core.ErrUserAlreadyInFavourites):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
		default:
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
		}
	}
	responseUser := user.ToResponseUser(&addedUser)
	return ctx.Status(fiber.StatusOK).JSON(responseUser)
}

func (r *Router) GetFavouriteUsers(ctx *fiber.Ctx) error {
	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	favouriteUsers, err := r.userService.GetFavouriteUsers(ctx.UserContext(), userID)
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
	responseUsers := make([]user.ResponseUser, 0, len(favouriteUsers))
	for i := range favouriteUsers {
		responseUsers = append(responseUsers, user.ToResponseUser(&favouriteUsers[i]))
	}
	return ctx.Status(fiber.StatusOK).JSON(responseUsers)
}

func (r *Router) DeleteUserFromFavourites(ctx *fiber.Ctx) error {
	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	favouriteUserID, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	err = r.userService.DeleteUserFromFavourite(ctx.UserContext(), favouriteUserID, userID)
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

	return ctx.Status(fiber.StatusNoContent).JSON(userID)
}
