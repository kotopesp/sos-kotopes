package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/report"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

// @Summary      Create a report
// @Description  Create a report for a specific post
// @Tags         reports
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        body  		     report.CreateRequestBodyReport  true  "Report data"
// @Success      201   {object}  model.Response                 "Report created successfully"
// @Failure      400   {object}  model.Response                 "Invalid request body or validation error"
// @Failure      401   {object}  model.Response                 "Unauthorized: Invalid or missing token"
// @Failure      404   {object}  model.Response                 "Post not found"
// @Failure      409   {object}  model.Response                 "Conflict: Report already exists"
// @Failure      422   {object}  model.Response{data=validator.Response}  "Unprocessable entity: Validation error"
// @Failure      500   {object}  model.Response                 "Internal server error"
// @Router       /reports [post]
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

		if errors.Is(err, core.ErrDuplicateReport) {
			logger.Log().Error(ctx.UserContext(), core.ErrDuplicateReport.Error())
			return ctx.Status(fiber.StatusConflict).JSON(model.ErrorResponse(core.ErrDuplicateReport.Error()))
		}

		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.SendStatus(fiber.StatusCreated)
}
