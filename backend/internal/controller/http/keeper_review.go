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
//	@Summary		get keeper reviews
//	@Description	Fetch reviews of a keeper based on optional pagination
//	@Tags			keeper review
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	false	"Limit"		default(10)
//	@Param			offset	query		int	false	"Offset"	default(0)
//	@Param			id		path		int	true	"Keeper ID"
//	@Success		200		{object}	model.Response{data=[]keeperreview.KeeperReviewsResponseWithUser}
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

	reviews, err := r.keeperService.GetAllReviews(ctx.UserContext(), params.FromKeeperReviewRequest(), id)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(Map(reviews, keeperreview.FromCoreKeeperReviewDetails)))
}

// Create review on a keeper
//	@Summary		create keeper review
//	@Description	Create review on a keeper
//	@Tags			keeper review
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Keeper ID"
//	@Success		200	{object}	model.Response{data=keeperreview.KeeperReviewsResponse}
//	@Failure		400	{object}	model.Response
//	@Failure		401	{object}	model.Response
//	@Failure		500	{object}	model.Response
//	@Router			/keepers/{id}/keeper_reviews [post]
func (r *Router) createKeeperReview(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	var newReview keeperreview.KeeperReviewsCreate

	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &newReview)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	authorID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}
	newReview.AuthorID = authorID
	newReview.KeeperID = id

	coreReview := newReview.ToCoreNewKeeperReview()
	if err := r.keeperService.CreateReview(ctx.UserContext(), coreReview); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(keeperreview.FromCoreKeeperReview(coreReview)))
}

// Updates a review on a keeper
//	@Summary		Update keeper review
//	@Description	Updates the keeper review details
//	@Tags			keeper review
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int									true	"Keeper ID"
//	@Param			body	body		keeperreview.KeeperReviewsUpdate	true	"Updated keeper review details"
//	@Success		200		{object}	model.Response{data=keeperreview.KeeperReviewsResponseWithUser}
//	@Failure		400		{object}	model.Response	"Invalid ID or request body"
//	@Failure		404		{object}	model.Response	"Review not found"
//	@Failure		401		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Router			/keeper_reviews/{id} [put]
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

	var updateReview keeperreview.KeeperReviewsUpdate

	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &updateReview)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	updateReview.ID = id
	updateReview.AuthorID = userID

	updatedKeeperReview, err := r.keeperService.UpdateReviewByID(ctx.UserContext(), updateReview.ToCoreUpdateKeeperReview())
	if err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(keeperreview.FromCoreKeeperReviewDetails(updatedKeeperReview)))
}

// Delete a review on a keeper
//	@Summary		Delete keeper review
//	@Description	Deletes the keeper review
//	@Tags			keeper review
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int									true	"Keeper ID"
//	@Param			body	body		keeperreview.KeeperReviewsUpdate	true	"Updated keeper review details"
//	@Success		200		{object}	model.Response
//	@Failure		204		{object}	model.Response	"Review not found"
//	@Failure		400		{object}	model.Response	"Invalid ID or request body"
//	@Failure		401		{object}	model.Response
//	@Failure		500		{object}	model.Response
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

	if err := r.keeperService.SoftDeleteReviewByID(ctx.UserContext(), id, userID); err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNoContent).JSON(model.ErrorResponse(err.Error()))
		}
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
