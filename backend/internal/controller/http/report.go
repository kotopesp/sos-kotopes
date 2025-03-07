package http

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/report"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

func (r *Router) createReport(ctx *fiber.Ctx) error {
	var createReport report.CreateRequestBodyReport

	fiberError, parseOrValidationError := parseQueryAndValidate(ctx, r.formValidator, &createReport)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), fiberError.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(core.ErrPostNotFound.Error()))
	}
	coreReport := createReport.ToCoreReport()

	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}

	coreReport.UserID = userID

	err = r.reportService.CreateReport(ctx.UserContext(), coreReport)
	if err != nil {
		if errors.Is(err, core.ErrPostNotFound) {
			logger.Log().Error(ctx.UserContext(), core.ErrPostNotFound.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(core.ErrPostNotFound.Error()))
		}

		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.SendStatus(fiber.StatusCreated)
}
