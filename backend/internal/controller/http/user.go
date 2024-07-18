package http

import (
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/user"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (r *Router) UpdateUser(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")

	if idStr == "" {
		logger.Log().Debug(ctx.UserContext(), "Error: id is required")
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "id is required",
		})
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	var update user.UpdateUser
	err = ctx.BodyParser(&update)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "invalid request body",
		})
	}
	err = r.userService.UpdateUser(ctx.UserContext(), id, update)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "success",
	})
}

func (r *Router) GetUser(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	if idStr == "" {
		logger.Log().Debug(ctx.UserContext(), "Error: id is required")
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "id is required",
		})
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	currentUser, err := r.userService.GetUser(ctx.UserContext(), id)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(currentUser)
}
