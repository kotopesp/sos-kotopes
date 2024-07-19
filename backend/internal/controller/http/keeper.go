package http

import (
	"strconv"

	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/keeper"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func (r *Router) getKeepers(ctx *fiber.Ctx) error {
	var params keeper.GetAllKeepersParams
	if err := ctx.QueryParser(&params); err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	panic("implement me")
}

func (r *Router) getKeeperByID(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	data, err := r.keeperService.GetByID(ctx.UserContext(), id)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(data))
}

func (r *Router) createKeeper(ctx *fiber.Ctx) error {
	panic("impl")
}

func (r *Router) updateKeeperById(ctx *fiber.Ctx) error {
	panic("impl")
}

func (r *Router) deleteKeeperById(ctx *fiber.Ctx) error {
	panic("impl")
}
