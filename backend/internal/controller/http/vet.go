package http

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	vetModel "github.com/kotopesp/sos-kotopes/internal/controller/http/model/vet"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

// @Summary		Get all vets
// @Tags			vets
// @Description	Get all vets with optional filtering, pagination, and sorting
// @Param			Location	query		string	false	"Location"
// @Param			MinRating	query		float64	false	"Minimum rating"
// @Param			MaxRating	query		float64	false	"Maximum rating"
// @Param			MinPrice	query		float64	false	"Minimum price"
// @Param			MaxPrice	query		float64	false	"Maximum price"
// @Param			SortBy		query		string	false	"Sort by field (name, rating, price)"
// @Param			SortOrder	query		string	false	"Sort order (asc, desc)"
// @Param			Limit		query		int		false	"Limit"		default(10)
// @Param			Offset		query		int		false	"Offset"	default(0)
// @Success		200			{object}	model.Response{data=core.Vets}
// @Failure		400			{object}	model.Response
// @Failure		500			{object}	model.Response
// @Router			/vets [get]
func (r *Router) getVets(c *fiber.Ctx) error {
	params := core.GetAllVetParams{}
	if err := c.QueryParser(&params); err != nil {
		logger.Log().Error(c.UserContext(), err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	vets, err := r.vetService.GetAll(c.Context(), params)
	if err != nil {
		logger.Log().Error(c.UserContext(), err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(model.OKResponse(vets))
}

// @Summary		Get vet by user ID
// @Tags			vets
// @Description	Get vet by user ID.
// @ID				get-vet-by-user-id
// @Accept			json
// @Produce		json
// @Param			userID	path		int	true	"User ID"
// @Success		200		{object}	model.Response{data=core.VetsDetails}
// @Failure		404		{object}	model.Response
// @Failure		500		{object}	model.Response
// @Router			/vets/{userID} [get]
func (r *Router) getVetByUserID(ctx *fiber.Ctx) error {
	println(1)
	userID, err1 := ctx.ParamsInt("userID")
	if err1 != nil {
		println(2)
		logger.Log().Debug(ctx.UserContext(), err1.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err1.Error()))
	}

	v, err := r.vetService.GetByUserID(ctx.Context(), userID)
	println(1)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		if errors.Is(err, core.ErrRecordNotFound) {
			println(1)
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(vetModel.FromCoreVetWithUser(v)))
}

// @Summary		Create a vet profile
// @Tags			vets
// @Description	Create a new vet profile
// @ID				create-vet
// @Accept			json
// @Produce		json
// @Param			request	body		vet.VetsCreate	true	"Vet"
// @Success		201		{object}	model.Response{data=vet.VetsResponse}
// @Failure		400		{object}	model.Response
// @Failure		401		{object}	model.Response
// @Failure		422		{object}	model.Response{data=validator.Response}
// @Failure		500		{object}	model.Response
// @Security		ApiKeyAuthBasic
// @Router			/vets [post]
func (r *Router) createVet(ctx *fiber.Ctx) error {
	var vetRequest vetModel.VetsCreate

	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &vetRequest)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(core.ErrFailedToGetUserIDFromToken))
	}

	_, err = r.vetService.GetByUserID(ctx.Context(), userID)
	if err == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Vet profile already exists for this user"))
	}

	if !errors.Is(err, core.ErrRecordNotFound) {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	vetRequest.UserID = userID
	coreVet := vetRequest.ToCoreNewVet()

	err = r.vetService.Create(ctx.UserContext(), coreVet)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse("Vet profile created successfully"))
}

// @Summary		Update vet by user id
// @Tags			vets
// @Description	Update vet by user id
// @ID				update-vet-by-user-id
// @Accept			json
// @Produce		json
// @Param			userID	path		int				true	"User ID"
// @Param			body	body		core.UpdateVets	true	"Vet data"
// @Success		200		{object}	model.Response{data=core.Vets}
// @Failure		400		{object}	model.Response
// @Failure		404		{object}	model.Response
// @Failure		500		{object}	model.Response
// @Security		ApiKeyAuthBasic
// @Router			/vets/{userID} [patch]
func (r *Router) updateVetByUserID(ctx *fiber.Ctx) error {
	userID, err := ctx.ParamsInt("userID")
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	var updateVet vetModel.VetsUpdate

	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &updateVet)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	updateVet.UserID = userID

	updatedVet, err := r.vetService.UpdateByUserID(ctx.UserContext(), updateVet.ToCoreUpdateVet())
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		if errors.Is(err, core.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(vetModel.FromCoreVetWithUser(updatedVet)))
}

// @Summary		Delete vet by user id
// @Tags			vets
// @Description	Delete vet by user id
// @ID				delete-vet-by-user-id
// @Accept			json
// @Produce		json
// @Param			id	path	int	true	"Vet ID"
// @Success		204
// @Failure		404	{object}	model.Response
// @Failure		500	{object}	model.Response
// @Security		ApiKeyAuthBasic
// @Router			/vets/{id} [delete]
func (r *Router) deleteVetByUserID(ctx *fiber.Ctx) error {
	userID, err := ctx.ParamsInt("userID")
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	if err := r.vetService.DeleteByUserID(ctx.UserContext(), userID); err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
