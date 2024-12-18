package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/keeper"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

// Retrieves a list of keepers with optional filtering and pagination
//
//	@Summary		get keepers
//	@Description	Fetch a list of keepers based on optional filters and pagination
//	@Tags			keeper
//	@Accept			json
//	@Produce		json
//	@Param			location				query		string	false	"Location"
//	@Param			min_rating				query		float64	false	"Minimum rating"
//	@Param			max_rating				query		float64	false	"Maximum rating"
//	@Param			min_price				query		float64	false	"Minimum price"
//	@Param			max_price				query		float64	false	"Maximum price"
//	@Param			has_cage				query		bool	false	"Has cage"
//	@Param			boarding_duration		query		string	false	"Boarding duration"
//	@Param			boarding_compensation	query		string	false	"Boarding compensation"
//	@Param			animal_acceptance		query		string	false	"Animal acceptance"
//	@Param			animal_category			query		string	false	"Animal category"
//	@Param			limit					query		int		false	"Limit"		default(10)
//	@Param			offset					query		int		false	"Offset"	default(0)
//	@Success		200						{object}	model.Response{data=keeper.ResponseKeepers}
//	@Failure		400						{object}	model.Response
//	@Failure		500						{object}	model.Response
//	@Router			/keepers [get]
func (r *Router) getKeepers(ctx *fiber.Ctx) error {
	var params keeper.GetAllKeepersParams

	fiberError, parseOrValidationError := parseQueryAndValidate(ctx, r.formValidator, &params)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), fiberError.Error())
		return fiberError
	}

	coreParams := params.ToCoreGetAllKeepersParams()

	coreKeepers, err := r.keeperService.GetAllKeepers(ctx.UserContext(), coreParams)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	pagination := paginate(len(coreKeepers), *params.Limit, *params.Offset)
	responseKeepers := keeper.ToModelResponseKeepers(pagination, coreKeepers)

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(responseKeepers))
}

// Retrieves a keeper by id
//
//	@Summary		Get keeper
//	@Description	Retrieves the details of a keeper by its ID
//	@Tags			keeper
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Keeper ID"
//	@Success		200	{object}	model.Response{data=keeper.ResponseKeeper}
//	@Failure		400	{object}	model.Response	"Invalid ID"
//	@Failure		404	{object}	model.Response	"Keeper not found"
//	@Failure		500	{object}	model.Response
//	@Router			/keepers/{id} [get]
func (r *Router) getKeeper(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	k, err := r.keeperService.GetKeepeByID(ctx.UserContext(), id)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		if errors.Is(err, core.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(keeper.ToModelResponseKeeper(k)))
}

// Create a keeper
//
//	@Summary		create keeper
//	@Description	Create a keeper
//	@Tags			keeper
//	@Accept			json
//	@Produce		json
//	@Param			body	body		keeper.CreateKeeper	true	"Create keeper"
//	@Success		201		{object}	model.Response{data=keeper.ResponseKeeper}
//	@Failure		400		{object}	model.Response
//	@Failure		401		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Router			/keepers [post]
func (r *Router) createKeeper(ctx *fiber.Ctx) error {
	var newKeeper keeper.CreateKeeper

	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &newKeeper)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}

	newKeeper.UserID = userID

	createdKeeper, err := r.keeperService.CreateKeeper(ctx.UserContext(), newKeeper.ToCoreKeeper())
	if err != nil {
		switch {
		case errors.Is(err, core.ErrKeeperUserAlreadyKeeper):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
		default:
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(keeper.ToModelResponseKeeper(createdKeeper)))
}

// Updates a keeper by id
//
//	@Summary		Update keeper
//	@Description	Updates the keeper details such as description, price, location, etc.
//	@Tags			keeper
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Keeper ID"
//	@Param			body	body		keeper.UpdateKeeper	true	"Updated keeper"
//	@Success		200		{object}	model.Response{data=keeper.ResponseKeeper}
//	@Failure		400		{object}	model.Response	"Invalid ID or request body"
//	@Failure		404		{object}	model.Response	"Keeper not found"
//	@Failure		401		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Router			/keepers/{id} [patch]
func (r *Router) updateKeeper(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}

	var updateKeeper keeper.UpdateKeeper

	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &updateKeeper)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	updatedKeeper, err := r.keeperService.UpdateKeeper(ctx.UserContext(), id, userID, updateKeeper.ToCoreUpdateKeeper())
	if err != nil {
		switch {
		case errors.Is(err, core.ErrRecordNotFound):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		case errors.Is(err, core.ErrKeeperUserIDMissmatch):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
		default:
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(keeper.ToModelResponseKeeper(updatedKeeper)))
}

// Delete keeper by id
//
//	@Summary		Delete keeper
//	@Description	Deletes a keeper by its ID.
//	@Tags			keeper
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int				true	"Keeper ID"
//	@Success		204	{object}	model.Response	"Successfully deleted"
//	@Failure		400	{object}	model.Response	"Invalid ID"
//	@Failure		401	{object}	model.Response
//	@Failure		404	{object}	model.Response	"Keeper not found"
//	@Failure		500	{object}	model.Response
//	@Router			/keepers/{id} [delete]
func (r *Router) deleteKeeper(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}

	if err := r.keeperService.DeleteKeeper(ctx.UserContext(), id, userID); err != nil {
		switch {
		case errors.Is(err, core.ErrRecordNotFound):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		case errors.Is(err, core.ErrKeeperUserIDMissmatch):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
		default:
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
		}
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
