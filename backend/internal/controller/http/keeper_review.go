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

	fiberError, parseOrValidationError := parseQueryAndValidate(ctx, r.formValidator, &params)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	reviews, err := r.keeperService.GetAllReviews(ctx.UserContext(), params.FromKeeperReviewRequest())
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(Map(reviews, keeperreview.FromCoreKeeperReviewDetails)))
}

func (r *Router) createKeeperReview(ctx *fiber.Ctx) error {
	var newReview keeperreview.KeeperReviewsCreate

	// parse review
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

	// create review
	coreReview := newReview.ToCoreNewKeeperReview()
	if err := r.keeperService.CreateReview(ctx.UserContext(), coreReview); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
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

	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}

	var updateReview keeperreview.KeeperReviewsUpdate

	// parse review
	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &updateReview)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	updateReview.ID = id
	updateReview.AuthorID = userID

	// update
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

func (r *Router) deleteKeeperReview(ctx *fiber.Ctx) error {
	// get id
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

	// delete
	if err := r.keeperService.SoftDeleteReviewByID(ctx.UserContext(), id, userID); err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNoContent).JSON(model.ErrorResponse(err.Error()))
		}
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
