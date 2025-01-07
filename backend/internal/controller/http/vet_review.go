package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	vetreview "github.com/kotopesp/sos-kotopes/internal/controller/http/model/vet_review"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

// @Summary		get vet reviews
// @Description	Fetch reviews of a vet based on optional pagination
// @Tags			vet review
// @Accept			json
// @Produce		json
// @Param			limit	query		int	false	"Limit"		default(10)
// @Param			offset	query		int	false	"Offset"	default(0)
// @Param			id		path		int	true	"Vet ID"
// @Success		200		{object}	model.Response{data=[]vet_review.VetReviewsResponse}
// @Failure		400		{object}	model.Response
// @Failure		500		{object}	model.Response
// @Router			/vets/{id}/vet_reviews [get]
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
//
//	@Summary		create vet review
//	@Description	Create review on a vet
//
//	@Tags			vet review
//
//	@Accept			json
//	@Produce		json
//
//	@Param			body	body		vet_review.VetReviewsCreate		true	"Create vet review"
//	@Success		201		{object}	vet_review.VetReviewsResponse	"Success response"
//	@Failure		400		{object}	model.Response
//	@Failure		401		{object}	model.Response
//	@Failure		500		{object}	model.Response
//
//	@Security		ApiKeyAuthBasic
//
//	@Router			/vets/{id}/vet_reviews [post]
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
//
//	@Summary		Update vet review
//	@Description	Updates the vet review details
//	@Tags			vet review
//	@Accept			json
//	@Produce		json
//
//	@Param			id		path		int													true	"Vet ID"
//	@Param			body	body		vet_review.VetReviewsUpdate							true	"Create vet review"
//	@Success		200		{object}	model.Response{data=vet_review.VetReviewsResponse}	"Success response"
//	@Failure		400		{object}	model.Response										"Invalid ID or request body"
//	@Failure		404		{object}	model.Response										"Review not found"
//
//	@Failure		401		{object}	model.Response
//	@Failure		500		{object}	model.Response
//
//	@Security		ApiKeyAuthBasic
//
//	@Router			/vet_reviews/{id} [patch]
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

// @Summary		Delete vet review
// @Description	Deletes the vet review
// @Tags			vet review
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Review ID"
// @Success		204	{object}	model.Response
// @Failure		400	{object}	model.Response	"Invalid ID or request body"
// @Failure		404	{object}	model.Response	"Review not found"
// @Failure		401	{object}	model.Response
// @Failure		500	{object}	model.Response
// @Security		ApiKeyAuthBasic
// @Router			/vet_reviews/{id} [delete]
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
