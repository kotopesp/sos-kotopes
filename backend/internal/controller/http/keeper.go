package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/keeper"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"gorm.io/gorm"
)

func (r *Router) getKeepers(ctx *fiber.Ctx) error {
	var params keeper.GetAllKeepersParams

	if err := ctx.QueryParser(params); err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	// if sortBy := ctx.Query("sortBy"); sortBy != "" {
	// 	params.SortBy = &sortBy
	// }
	// if sortOrder := ctx.Query("sortOrder"); sortOrder != "" {
	// 	params.SortOrder = &sortOrder
	// }
	// if location := ctx.Query("location"); location != "" {
	// 	params.Location = &location
	// }

	// if minPriceStr := ctx.Query("min_price"); minPriceStr != "" {
	// 	minPrice, err := strconv.ParseFloat(minPriceStr, 64)
	// 	if err != nil {
	// 		logger.Log().Debug(ctx.UserContext(), err.Error())
	// 		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	// 	}
	// 	params.MinPrice = &minPrice
	// }
	// if maxPriceStr := ctx.Query("max_price"); maxPriceStr != "" {
	// 	maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
	// 	if err != nil {
	// 		logger.Log().Debug(ctx.UserContext(), err.Error())
	// 		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	// 	}
	// 	params.MaxPrice = &maxPrice
	// }
	// if minRatingStr := ctx.Query("min_rating"); minRatingStr != "" {
	// 	minRating, err := strconv.ParseFloat(minRatingStr, 64)
	// 	if err != nil {
	// 		logger.Log().Debug(ctx.UserContext(), err.Error())
	// 		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	// 	}
	// 	params.MinRating = &minRating
	// }
	// if maxRatingStr := ctx.Query("max_rating"); maxRatingStr != "" {
	// 	maxRating, err := strconv.ParseFloat(maxRatingStr, 64)
	// 	if err != nil {
	// 		logger.Log().Debug(ctx.UserContext(), err.Error())
	// 		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	// 	}
	// 	params.MaxRating = &maxRating
	// }

	// if limit := ctx.Query("limit"); limit != "" {
	// 	l, err := strconv.Atoi(limit)
	// 	if err == nil {
	// 		params.Limit = &l
	// 	}
	// }
	// if offset := ctx.Query("offset"); offset != "" {
	// 	o, err := strconv.Atoi(offset)
	// 	if err == nil {
	// 		params.Offset = &o
	// 	}
	// }

	usrCtx := ctx.UserContext()
	coreParams := params.FromKeeperRequest()
	coreKeepers, err := r.keeperService.GetAll(&usrCtx, coreParams)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	responseKeepers := Map(coreKeepers, func(k core.Keepers) keeper.Keepers {
		return keeper.Keepers{
			ID:          k.ID,
			UserID:      k.UserID,
			Description: k.Description,
			Price:       k.Price,
			Location:    k.Location,
			CreatedAt:   k.CreatedAt,
			UpdatedAt:   k.UpdatedAt,
		}
	})

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(responseKeepers))
}

func (r *Router) getKeeperByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 0)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	usrCtx := ctx.UserContext()
	k, err := r.keeperService.GetByID(&usrCtx, int(id))
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(keeper.Keepers{
		ID:          k.ID,
		UserID:      k.UserID,
		Description: k.Description,
		Price:       k.Price,
		Location:    k.Location,
		CreatedAt:   k.CreatedAt,
		UpdatedAt:   k.UpdatedAt,
	}))
}

func (r *Router) createKeeper(ctx *fiber.Ctx) error {
	var newKeeper keeper.KeepersCreate

	fiberError, parseOrValidationError := parseAndValidateAny(ctx, r.formValidator, &newKeeper)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	// parse keeper
	// if err := ctx.BodyParser(&newKeeper); err != nil {
	// 	logger.Log().Debug(ctx.UserContext(), err.Error())
	// 	return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	// }

	// make sure the create time is set
	// if keeper.CreatedAt.IsZero() {
	// 	keeper.CreatedAt = time.Now()
	// }

	// create keeper
	k := newKeeper.ToCoreNewKeeper()
	usrCtx := ctx.UserContext()
	if err := r.keeperService.Create(&usrCtx, k); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(keeper.Keepers{
		ID:          k.ID,
		UserID:      k.UserID,
		Description: k.Description,
		Price:       k.Price,
		Location:    k.Location,
		CreatedAt:   k.CreatedAt,
		UpdatedAt:   k.UpdatedAt,
	}))
}

func (r *Router) updateKeeperByID(ctx *fiber.Ctx) error {
	// get id
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 0)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	var updateKeeper keeper.KeepersUpdate

	// parse keeper
	fiberError, parseOrValidationError := parseAndValidateAny(ctx, r.formValidator, &updateKeeper)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	// update
	var usrCtx = ctx.UserContext()
	if err := r.keeperService.UpdateByID(&usrCtx, core.Keepers{
		ID:          int(id),
		Description: updateKeeper.Description,
		Price:       updateKeeper.Price,
		Location:    updateKeeper.Location,
	}); err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(updateKeeper))
}

func (r *Router) deleteKeeperByID(ctx *fiber.Ctx) error {
	// get id
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 0)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	// delete
	var usrCtx = ctx.UserContext()
	if err := r.keeperService.DeleteByID(&usrCtx, int(id)); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
