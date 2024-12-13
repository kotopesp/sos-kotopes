package http

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/vet"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

// @Summary		Get all vets
// @Tags			vets
// @Description	Get all vets
// @ID				get-vets
// @Accept			json
// @Produce		json
// @Success		200	{object}	model.Response{data=core.Vets}
// @Failure		500	{object}	model.Response
// @Router			/vets [get]
func (r *Router) getVets(c *fiber.Ctx) error {
	params := core.GetAllVetParams{}
	// Here we can parse query params from the request to pass to the service layer
	if err := c.QueryParser(&params); err != nil {
		logger.Log().Error(c.UserContext(), err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	// Calling service method to get vets
	vets, err := r.vetService.GetAll(c.Context(), params)
	if err != nil {
		logger.Log().Error(c.UserContext(), err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	// Returning the response with the list of vets
	return c.Status(fiber.StatusOK).JSON(model.OKResponse(vets))
}

// @Summary				Get vet by user ID
// @Tags				vets
// @Description			Get vet by user ID.
// @ID					get-vet-by-user-id
// @Accept				json
// @Produce				json
// @Param				userID			path		int	true	"User ID"
// @Success				200	{object}	model.Response{data=core.VetsDetails}
// @Failure				404	{object}	model.Response
// @Failure				500	{object}	model.Response
// @Router				/vet/{userID} [get]
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

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(vet.FromCoreVetWithUser(v)))
}

// @Summary			Update vet by user id
// @Tags			vets
// @Description		Update vet by user id
// @ID				update-vet-by-user-id
// @Accept			json
// @Produce			json
// @Param			userID			path				int				true	"User ID"
// @Param			body			core.UpdateVets		true		"Vet data"
// @Success			200	{object}	model.Response{data=core.Vets}
// @Failure			400	{object}	model.Response
// @Failure			404	{object}	model.Response
// @Failure			500	{object}	model.Response
// @Router			/vets/{userID} [put]
func (r *Router) updateVetByUserID(ctx *fiber.Ctx) error {
	userID, err := ctx.ParamsInt("userID")
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	var updateVet vet.VetsUpdate

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

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(vet.FromCoreVetWithUser(updatedVet)))
}

// @Summary		Delete vet by user id
// @Tags			vets
// @Description	Delete vet by user id
// @ID				delete-vet-by-user-id
// @Accept			json
// @Produce		json
// @Param			id		path		int		true	"Vet ID"
// @Success		204
// @Failure		404	{object}	model.Response
// @Failure		500	{object}	model.Response
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
