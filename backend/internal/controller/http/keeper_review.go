package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	keeperreview "github.com/kotopesp/sos-kotopes/internal/controller/http/model/keeper_review"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

// Retrieves reviews of a keeper with optional pagination
//
//	@Summary		get keeper reviews
//	@Description	Fetch reviews of a keeper based on optional pagination
//	@Tags			keeper review
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	false	"Limit"		default(10)
//	@Param			offset	query		int	false	"Offset"	default(0)
//	@Param			id		path		int	true	"Keeper ID"
//	@Success		200		{object}	model.Response{data=[]keeperreview.ResponseKeeperReview}
//	@Failure		400		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Router			/keepers/{id}/keeper_reviews [get]
func (r *Router) getKeeperReviews(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	var params keeperreview.GetAllKeeperReviewsParams

	fiberError, parseOrValidationError := parseQueryAndValidate(ctx, r.formValidator, &params)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	coreReviews, err := r.keeperService.GetAllReviews(ctx.UserContext(), id, params.ToCoreGetAllKeeperReviewsParams())
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	responseReviews := Map(coreReviews, keeperreview.ToModelResponseKeeperReview)

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(responseReviews))
}

// Create review on a keeper
//
//	@Summary		create keeper review
//	@Description	Create review on a keeper
//	@Tags			keeper review
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int								true	"Keeper ID"
//	@Param			body	body		keeperreview.CreateKeeperReview	true	"Create keeper review"
//	@Success		201		{object}	model.Response{data=keeperreview.ResponseKeeperReview}
//	@Failure		400		{object}	model.Response
//	@Failure		401		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Router			/keepers/{id}/keeper_reviews [post]
func (r *Router) createKeeperReview(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	var newReview keeperreview.CreateKeeperReview

	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &newReview)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}

	coreNewReview := newReview.ToCoreKeeperReview()
	coreNewReview.AuthorID = userID
	coreNewReview.KeeperID = id

	createdReview, err := r.keeperService.CreateReview(ctx.UserContext(), coreNewReview)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrKeeperReviewToItself):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
		default:
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(keeperreview.ToModelResponseKeeperReview(createdReview)))
}

// Updates a review on a keeper
//
//	@Summary		Update keeper review
//	@Description	Updates the keeper review details
//	@Tags			keeper review
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int								true	"Keeper ID"
//	@Param			body	body		keeperreview.UpdateKeeperReview	true	"Updated keeper review"
//	@Success		200		{object}	model.Response{data=keeperreview.ResponseKeeperReview}
//	@Failure		400		{object}	model.Response	"Invalid ID or request body"
//	@Failure		404		{object}	model.Response	"Review not found"
//	@Failure		401		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Router			/keeper_reviews/{id} [patch]
func (r *Router) updateKeeperReview(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}

	var updateReview keeperreview.UpdateKeeperReview

	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &updateReview)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}
	print(45)
	updatedReview, err := r.keeperService.UpdateReview(ctx.UserContext(), id, userID, updateReview.ToCoreUpdateKeeperReview())
	print(23)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrRecordNotFound):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		case errors.Is(err, core.ErrKeeperReviewUserIDMissmatch):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
		default:
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(keeperreview.ToModelResponseKeeperReview(updatedReview)))
}

// Delete a review on a keeper
//
//	@Summary		Delete keeper review
//	@Description	Deletes the keeper review
//	@Tags			keeper review
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Keeper ID"
//	@Success		204	{object}	model.Response
//	@Failure		400	{object}	model.Response	"Invalid ID or request body"
//	@Failure		404	{object}	model.Response	"Review not found"
//	@Failure		401	{object}	model.Response
//	@Failure		500	{object}	model.Response
//	@Router			/keeper_reviews/{id} [delete]
func (r *Router) deleteKeeperReview(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}

	if err := r.keeperService.DeleteReview(ctx.UserContext(), id, userID); err != nil {
		switch {
		case errors.Is(err, core.ErrRecordNotFound):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		case errors.Is(err, core.ErrKeeperReviewToItself):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
		default:
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
		}
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
