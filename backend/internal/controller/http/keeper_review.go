package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	keeperreview "github.com/kotopesp/sos-kotopes/internal/controller/http/model/keeper_review"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"gorm.io/gorm"
)

func (r *Router) getKeeperReviews(ctx *fiber.Ctx) error {
	var params keeperreview.GetAllKeeperReviewsParams

	fiberError, parseOrValidationError := parseAndValidateQueryAny(ctx, r.formValidator, &params)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	usrCtx := ctx.UserContext()
	reviews, err := r.keeperReviewsService.GetAll(&usrCtx, params.FromKeeperReviewRequest())
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(Map(reviews, func(review core.KeeperReviews) keeperreview.KeeperReviews {
		return keeperreview.KeeperReviews{
			ID:        review.ID,
			AuthorID:  review.AuthorID,
			Content:   review.Content,
			Grade:     review.Grade,
			KeeperID:  review.KeeperID,
			IsDeleted: review.IsDeleted,
			DeletedAt: review.DeletedAt,
			CreatedAt: review.CreatedAt,
			UpdatedAt: review.UpdatedAt,
		}
	})))
}

func (r *Router) createKeeperReview(ctx *fiber.Ctx) error {
	var newReview keeperreview.KeeperReviewsCreate

	// parse review
	fiberError, parseOrValidationError := parseAndValidateAny(ctx, r.formValidator, &newReview)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	// create review
	usrCtx := ctx.UserContext()
	coreReview := newReview.ToCoreNewKeeperReview()
	if err := r.keeperReviewsService.Create(&usrCtx, coreReview); err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(keeperreview.KeeperReviews{
		ID:        coreReview.ID,
		AuthorID:  coreReview.AuthorID,
		Content:   coreReview.Content,
		Grade:     coreReview.Grade,
		KeeperID:  coreReview.KeeperID,
		IsDeleted: coreReview.IsDeleted,
		DeletedAt: coreReview.DeletedAt,
		CreatedAt: coreReview.CreatedAt,
		UpdatedAt: coreReview.UpdatedAt,
	}))
}

func (r *Router) updateKeeperReview(ctx *fiber.Ctx) error {
	// get id
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 0)
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
	var usrCtx = ctx.UserContext()
	if err := r.keeperReviewsService.UpdateByID(&usrCtx, core.KeeperReviews{
		ID:      int(id),
		Content: updateReview.Content,
		Grade:   updateReview.Grade,
	}); err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(updateReview))
}

func (r *Router) deleteKeeperReview(ctx *fiber.Ctx) error {
	// get id
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 0)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	// delete
	var usrCtx = ctx.UserContext()
	if err := r.keeperReviewsService.SoftDeleteByID(&usrCtx, int(id)); err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
