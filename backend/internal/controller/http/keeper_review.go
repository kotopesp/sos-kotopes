package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	keeperreview "github.com/kotopesp/sos-kotopes/internal/controller/http/model/keeper_review"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

func (r *Router) getKeeperReviews(ctx *fiber.Ctx) error {
	var params keeperreview.GetAllKeeperReviewsParams

	fiberError, parseOrValidationError := parseAndValidateQueryAny(ctx, r.formValidator, &params)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	reviews, err := r.keeperReviewsService.GetAll(ctx.UserContext(), params.FromKeeperReviewRequest())
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(Map(reviews, func(review core.KeeperReviews) keeperreview.KeeperReviewsResponse {
		return keeperreview.FromCoreKeeperReview(review)
	})))
}

func (r *Router) createKeeperReview(ctx *fiber.Ctx) error {
	var newReview keeperreview.KeeperReviewsCreate

	// parse review
	fiberError, parseOrValidationError := parseAndValidateAny(ctx, r.formValidator, &newReview)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	authorId, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}
	newReview.AuthorID = authorId

	// create review
	coreReview := newReview.ToCoreNewKeeperReview()
	if err := r.keeperReviewsService.Create(ctx.UserContext(), coreReview); err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(keeperreview.FromCoreKeeperReview(coreReview)))
}

func (r *Router) updateKeeperReview(ctx *fiber.Ctx) error {
	// get id
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	var updateReview keeperreview.KeeperReviewsUpdate

	// parse review
	fiberError, parseOrValidationError := parseAndValidateAny(ctx, r.formValidator, &updateReview)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	// update
	if err := r.keeperReviewsService.UpdateByID(ctx.UserContext(), core.KeeperReviews{
		ID:      int(id),
		Content: updateReview.Content,
		Grade:   updateReview.Grade,
	}); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		if errors.Is(err, core.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (r *Router) deleteKeeperReview(ctx *fiber.Ctx) error {
	// get id
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	// delete
	if err := r.keeperReviewsService.SoftDeleteByID(ctx.UserContext(), int(id)); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		if errors.Is(err, core.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNoContent).JSON(model.ErrorResponse(err.Error()))
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}