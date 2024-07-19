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

	if sortBy := ctx.Query("sortBy"); sortBy != "" {
		params.SortBy = &sortBy
	}
	if sortOrder := ctx.Query("sortOrder"); sortOrder != "" {
		params.SortOrder = &sortOrder
	}
	if location := ctx.Query("location"); location != "" {
		params.Location = &location
	}
	if minRating := ctx.Query("minRating"); minRating != "" {
		rating, err := strconv.ParseFloat(minRating, 64)
		if err == nil {
			params.MinRating = &rating
		}
	}
	if maxRating := ctx.Query("maxRating"); maxRating != "" {
		rating, err := strconv.ParseFloat(maxRating, 64)
		if err == nil {
			params.MaxRating = &rating
		}
	}
	if limit := ctx.Query("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err == nil {
			params.Limit = &l
		}
	}
	if offset := ctx.Query("offset"); offset != "" {
		o, err := strconv.Atoi(offset)
		if err == nil {
			params.Offset = &o
		}
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
