package http

import (
	"strconv"

	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func (r *Router) getKeepers(ctx *fiber.Ctx) error {
	var params core.GetAllKeepersParams
	if err := ctx.QueryParser(&params); err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	usrCtx := ctx.UserContext()
	var keepers []core.Keepers
	keepers, err := r.keeperService.GetAll(&usrCtx, params)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(keepers))
}

func (r *Router) getKeeperByID(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	usrCtx := ctx.UserContext()
	data, err := r.keeperService.GetByID(&usrCtx, id)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(data))
}

func (r *Router) createKeeper(ctx *fiber.Ctx) error {
	usrCtx := ctx.UserContext()
	err := r.keeperService.Create(&usrCtx, core.Keepers{
		ID:          0,
		UserID:      0,
		Description: "hello keeper",
		Location:    "asdasd",
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *Router) updateKeeperById(ctx *fiber.Ctx) error {
	panic("impl")
}

func (r *Router) deleteKeeperById(ctx *fiber.Ctx) error {
	panic("impl")
}
