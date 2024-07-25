package http

import (
	"strconv"
	"time"

	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (r *Router) getKeeperReviews(ctx *fiber.Ctx) error {
	params := core.GetAllKeeperReviewsParams{}

	if limit := ctx.Query("limit"); limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err == nil {
			*params.Limit = limitInt
		}
	}

	if offset := ctx.Query("offset"); offset != "" {
		offsetInt, err := strconv.Atoi(offset)
		if err == nil {
			*params.Offset = offsetInt
		}
	}

	usrCtx := ctx.UserContext()
	reviews, err := r.KeeperReviewsService.GetAll(&usrCtx, params)

	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(reviews))
}

func (r *Router) createKeeperReview(ctx *fiber.Ctx) error {
	var review core.KeeperReviews

	// parse review
	if err := ctx.BodyParser(&review); err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	// make sure the create time is set
	if review.CreatedAt.IsZero() {
		review.CreatedAt = time.Now()
	}

	// check grade boundaries
	if review.Grade < 1 || review.Grade > 5 {
		errMsg := "Grade must be between 1 and 5"
		logger.Log().Debug(ctx.UserContext(), errMsg)
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(errMsg))
	}

	// create review
	usrCtx := ctx.UserContext()
	if err := r.KeeperReviewsService.Create(&usrCtx, review); err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(review))
}

func (r *Router) updateKeeperReview(ctx *fiber.Ctx) error {
	// get id
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	var review core.KeeperReviews

	// parse review
	if err := ctx.BodyParser(&review); err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	review.ID = id

	// update
	var usrCtx = ctx.UserContext()
	if err := r.KeeperReviewsService.UpdateById(&usrCtx, review); err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(review))
}

func (r *Router) deleteKeeperReview(ctx *fiber.Ctx) error {
	// get id
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	// delete
	var usrCtx = ctx.UserContext()
	if err := r.KeeperReviewsService.DeleteById(&usrCtx, id); err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
