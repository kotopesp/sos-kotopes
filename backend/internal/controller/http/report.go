package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/report"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

// @Summary		Create a report
// @Description	Create a report for a specific
// @Tags			reports
// @Accept			json
// @Produce		json
// @Param			body	body	report.CreateRequestBodyReport	true	"Report data"	"Report data"
//
// @Success		201		"Report created successfully"
// @Failure		400		{object}	model.Response							"Invalid request body"
// @Failure		401		{object}	model.Response							"Unauthorized: Invalid or missing token"
// @Failure		404		{object}	model.Response							"Content not found"
// @Failure		409		{object}	model.Response							"Conflict: Report already exists"
// @Failure		422		{object}	model.Response{data=validator.Response}	"Validation error"
// @Failure		500		{object}	model.Response							"Internal server error"
// @Security		ApiKeyAuthBasic
// @Router			/reports [post]
func (r *Router) createReport(ctx *fiber.Ctx) error {
	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}

	var createReport report.CreateRequestBodyReport

	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &createReport)
	if fiberError != nil {
		logger.Log().Error(ctx.UserContext(), fiberError.Error())
		return fiberError
	}

	if parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), parseOrValidationError.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ErrorResponse("Invalid request body"))
	}

	coreReport := createReport.ToCoreReport(userID)

	err = r.reportService.CreateReport(ctx.UserContext(), coreReport)
	if err != nil {
		if errors.Is(err, core.ErrTargetNotFound) {
			logger.Log().Error(ctx.UserContext(), core.ErrTargetNotFound.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(core.ErrTargetNotFound.Error()))
		}

		if errors.Is(err, core.ErrDuplicateReport) {
			logger.Log().Error(ctx.UserContext(), core.ErrDuplicateReport.Error())
			return ctx.Status(fiber.StatusConflict).JSON(model.ErrorResponse(core.ErrDuplicateReport.Error()))
		}

		if errors.Is(err, core.ErrInvalidReportableType) {
			logger.Log().Error(ctx.UserContext(), core.ErrInvalidReportableType.Error())
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(core.ErrInvalidReportableType.Error()))
		}

		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse("Internal server error"))
	}

	return ctx.SendStatus(fiber.StatusCreated)
}
