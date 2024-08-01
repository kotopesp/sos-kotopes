package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/keeper"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"gorm.io/gorm"
)

func (r *Router) getKeepers(ctx *fiber.Ctx) error {
	var params keeper.GetAllKeepersParams

	fiberError, parseOrValidationError := parseAndValidateQueryAny(ctx, r.formValidator, &params)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	coreParams := params.FromKeeperRequest()
	coreKeepers, err := r.keeperService.GetAll(ctx.UserContext(), coreParams)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	totalKeepers := len(coreKeepers)
	currentCoreKeepers := coreKeepers[*coreParams.Offset:min(*coreParams.Offset+*coreParams.Limit, totalKeepers)]
	responseKeepers := fiber.Map{
		"meta": generatePaginationMeta(totalKeepers, params.Offset, params.Limit),
		"data": Map(currentCoreKeepers, func(k core.Keepers) keeper.Keepers {
			return keeper.Keepers{
				ID:          k.ID,
				UserID:      k.UserID,
				Description: k.Description,
				Price:       k.Price,
				Location:    k.Location,
				CreatedAt:   k.CreatedAt,
				UpdatedAt:   k.UpdatedAt,
			}
		}),
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(responseKeepers))
}

func (r *Router) getKeeperByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	k, err := r.keeperService.GetByID(ctx.UserContext(), int(id))
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

	// create keeper
	k := newKeeper.ToCoreNewKeeper()
	if err := r.keeperService.Create(ctx.UserContext(), k); err != nil {
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
	id, err := ctx.ParamsInt("id")
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
	if err := r.keeperService.UpdateByID(ctx.UserContext(), core.Keepers{
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
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	// delete
	if err := r.keeperService.DeleteByID(ctx.UserContext(), int(id)); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
