package http

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/seeker"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

// @Summary		Get seeker
// @Tags			seeker
// @Description	Get seeker by id
// @ID				get-seeker
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"User ID"
// @Success		200	{object}	model.Response{data=seeker.ResponseSeekers}
// @Failure		400	{object}	model.Response
// @Failure		500	{object}	model.Response
// @Security		ApiKeyAuthBasic
// @Router			/seekers/{id}  [get]
func (r *Router) getSeeker(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	if err != nil {
		logger.Log().Debug(ctx.Context(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	getSeeker, err := r.seekerService.GetSeeker(ctx.UserContext(), id)

	if err != nil {
		switch {
		case errors.Is(err, core.ErrNoSuchUser):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		default:
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
		}
	}

	responseSeeker := seeker.ToResponseSeeker(&getSeeker)
	return ctx.Status(fiber.StatusOK).JSON(responseSeeker)
}

// @Summary		Create a seeker
// @Tags			seeker
// @Description	Create a seeker
// @ID				create-seeker
// @Accept			json
// @Produce		json
// @Param			request	body		seeker.ResponseSeekers	true	"Seeker"
// @Success		200		{object}	model.Response{data=seeker.ResponseSeekers}
// @Failure		400		{object}	model.Response
// @Failure		500		{object}	model.Response
// @Security		ApiKeyAuthBasic
// @Router			/seekers [post]
func (r *Router) createSeeker(ctx *fiber.Ctx) error {
	var responseSeeker seeker.ResponseSeekers

	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &responseSeeker)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	coreSeeker, err := r.seekerService.CreateSeeker(ctx.UserContext(), responseSeeker.ToCoreSeeker())
	if err != nil {
		switch {
		case errors.Is(err, core.ErrNoSuchUser):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		default:
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
		}
	}

	responseSeeker = seeker.ToResponseSeeker(&coreSeeker)
	return ctx.Status(fiber.StatusOK).JSON(responseSeeker)
}

func (r *Router) updateSeeker(ctx *fiber.Ctx) error { return nil }