package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/keeper"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

func (r *Router) getKeepers(ctx *fiber.Ctx) error {
	var params keeper.GetAllKeepersParams

	fiberError, parseOrValidationError := parseQueryAndValidate(ctx, r.formValidator, &params)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), fiberError.Error())
		return fiberError
	}

	coreParams := params.ToCoreGetAllKeepersParams()

	coreKeepersDetails, err := r.keeperService.GetAll(ctx.UserContext(), coreParams)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	pagination := paginate(len(coreKeepersDetails), params.Limit, params.Offset)
	responseKeepers := keeper.ToKeepersResponse(pagination, coreKeepersDetails)

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(responseKeepers))
}

func (r *Router) getKeeperByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	k, err := r.keeperService.GetByID(ctx.UserContext(), id)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		if errors.Is(err, core.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(keeper.FromCoreKeeperDetails(k)))
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
	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &updateKeeper)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	updateKeeper.ID = id

	// update
	updatedKeeper, err := r.keeperService.UpdateByID(ctx.UserContext(), updateKeeper.ToCoreUpdatedKeeper())
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		if errors.Is(err, core.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(keeper.FromCoreKeeperDetails(updatedKeeper)))
}

func (r *Router) deleteKeeperByID(ctx *fiber.Ctx) error {
	// get id
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	// delete
	if err := r.keeperService.SoftDeleteByID(ctx.UserContext(), id); err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNoContent).JSON(model.ErrorResponse(err.Error()))
		}
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
