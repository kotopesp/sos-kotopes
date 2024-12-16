package http

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	vetModel "github.com/kotopesp/sos-kotopes/internal/controller/http/model/vet"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

//	@Summary		Get all vets
//	@Tags			vets
//	@Description	Get all vets
//	@ID				get-vets
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	model.Response{data=core.Vets}
//	@Failure		500	{object}	model.Response
//	@Router			/vets [get]
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

//	@Summary		Get vet by user ID
//	@Tags			vets
//	@Description	Get vet by user ID.
//	@ID				get-vet-by-user-id
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		200		{object}	model.Response{data=core.VetsDetails}
//	@Failure		404		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Router			/vet/{userID} [get]
func (r *Router) getVetByUserID(ctx *fiber.Ctx) error {
	userID, err1 := ctx.ParamsInt("userID")
	if err1 != nil {
		logger.Log().Debug(ctx.UserContext(), err1.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err1.Error()))
	}

	v, err := r.vetService.GetByUserID(ctx.Context(), userID)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		if errors.Is(err, core.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(vetModel.FromCoreVetWithUser(v)))
}

//	@Summary		Create a vet profile
//	@Tags			vets
//	@Description	Create a new vet profile
//	@ID				create-vet
//	@Accept			json
//	@Produce		json
//	@Param			is_organization			formData	bool	true	"Is Organization"
//	@Param			patronymic				formData	string	false	"Patronymic"
//	@Param			education				formData	string	false	"Education"
//	@Param			org_name				formData	string	false	"Organization Name"
//	@Param			location				formData	string	true	"Location"
//	@Param			org_email				formData	string	false	"Organization Email"
//	@Param			inn_number				formData	string	false	"INN Number"
//	@Param			is_remote_consulting	formData	bool	false	"Is Remote Consulting"
//	@Param			is_inpatient			formData	bool	false	"Is Inpatient"
//	@Param			description				formData	string	false	"Description"
//	@Success		201						{object}	model.Response{data=vet.VetsResponse}
//	@Failure		400						{object}	model.Response
//	@Failure		401						{object}	model.Response
//	@Failure		422						{object}	model.Response{data=validator.Response}
//	@Failure		500						{object}	model.Response
//	@Security		ApiKeyAuthBasic
//	@Router			/vets [post]
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

	vetRequest.UserID = userID

	coreVet := vetRequest.ToCoreNewVet()

	err = r.vetService.Create(ctx.UserContext(), coreVet)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse("Vet profile created successfully"))
}

//	@Summary		Update vet by user id
//	@Tags			vets
//	@Description	Update vet by user id
//	@ID				update-vet-by-user-id
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int				true	"User ID"
//	@Param			body	body		core.UpdateVets	true	"Vet data"
//	@Success		200		{object}	model.Response{data=core.Vets}
//	@Failure		400		{object}	model.Response
//	@Failure		404		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Router			/vets/{userID} [put]
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

//	@Summary		Delete vet by user id
//	@Tags			vets
//	@Description	Delete vet by user id
//	@ID				delete-vet-by-user-id
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Vet ID"
//	@Success		204
//	@Failure		404	{object}	model.Response
//	@Failure		500	{object}	model.Response
//	@Router			/vets/{id} [delete]
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