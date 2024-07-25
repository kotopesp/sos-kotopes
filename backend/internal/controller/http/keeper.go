package http

import (
	"strconv"
	"time"

	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
	// if minRating := ctx.Query("minRating"); minRating != "" {
	// 	rating, err := strconv.ParseFloat(minRating, 64)
	// 	if err == nil {
	// 		params.MinRating = &rating
	// 	}
	// }
	// if maxRating := ctx.Query("maxRating"); maxRating != "" {
	// 	rating, err := strconv.ParseFloat(maxRating, 64)
	// 	if err == nil {
	// 		params.MaxRating = &rating
	// 	}
	// }
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
	var keeper core.Keepers

	// parse keeper
	if err := ctx.BodyParser(&keeper); err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	// make sure the create time is set
	if keeper.CreatedAt.IsZero() {
		keeper.CreatedAt = time.Now()
	}

	// create keeper
	usrCtx := ctx.UserContext()
	if err := r.keeperService.Create(&usrCtx, keeper); err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(keeper))
}

func (r *Router) updateKeeperById(ctx *fiber.Ctx) error {
	// get id
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	var keeper core.Keepers

	// parse keeper
	if err := ctx.BodyParser(&keeper); err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	keeper.ID = id

	// update
	var usrCtx = ctx.UserContext()
	if err := r.keeperService.UpdateById(&usrCtx, keeper); err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(keeper))
}

func (r *Router) deleteKeeperById(ctx *fiber.Ctx) error {
	// get id
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	// delete
	var usrCtx = ctx.UserContext()
	if err := r.keeperService.DeleteById(&usrCtx, id); err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}