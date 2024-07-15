package http

import (
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (r *Router) ChangeName(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	if idStr == "" {
		logger.Log().Debug(ctx.UserContext(), "Error: id is required")
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "id is required",
		})
	}
	nameStr := ctx.Params("username")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	err = r.userService.ChangeName(ctx.UserContext(), id, nameStr)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "success",
	})
}

func (r *Router) ChangeDescription(ctx *fiber.Ctx) error {
	return nil
}
