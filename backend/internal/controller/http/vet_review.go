package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	vetreview "github.com/kotopesp/sos-kotopes/internal/controller/http/model/vet_review"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

// Get Vet Reviews
func (r *Router) getVetReviews(ctx *fiber.Ctx) error {
	var params vetreview.GetAllVetReviewsParams

	fiberError, parseOrValidationError := parseQueryAndValidate(ctx, r.formValidator, &params)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	reviews, err := r.vetService.GetAllReviews(ctx.UserContext(), params.FromVetReviewRequest())
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	mappedReviews := make([]vetreview.VetReviewsResponseWithUser, len(reviews))
	for i, review := range reviews {
		mappedReviews[i] = vetreview.FromCoreVetReviewDetails(review)
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(mappedReviews))
}

// Create Vet Review
func (r *Router) createVetReview(ctx *fiber.Ctx) error {
	var newReview vetreview.VetReviewsCreate

	// Parse review
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

	// Create review
	coreReview := newReview.ToCoreNewVetReview()
	if err := r.vetService.CreateReview(ctx.UserContext(), coreReview); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(vetreview.FromCoreVetReview(coreReview)))
}

// Update Vet Review
func (r *Router) updateVetReview(ctx *fiber.Ctx) error {
	// Get ID
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

	var updateReview vetreview.VetReviewsUpdate

	// Parse review
	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &updateReview)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	updateReview.ID = id
	updateReview.AuthorID = userID

	// Update
	updatedVetReview, err := r.vetService.UpdateReviewByID(ctx.UserContext(), updateReview.ToCoreUpdateVetReview())
	if err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(vetreview.FromCoreVetReviewDetails(updatedVetReview)))
}

// Delete Vet Review
func (r *Router) deleteVetReview(ctx *fiber.Ctx) error {
	// Get ID
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

	// Delete
	if err := r.vetService.SoftDeleteReviewByID(ctx.UserContext(), id, userID); err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNoContent).JSON(model.ErrorResponse(err.Error()))
		}
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
