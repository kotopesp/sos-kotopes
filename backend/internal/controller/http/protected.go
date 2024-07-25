package http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
)

func (r *Router) protected(ctx *fiber.Ctx) error {
	// getting id from token
	id, err := getIDFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	// getting username from token
	username, err := getUsernameFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(fiber.Map{
		"UserID":   fmt.Sprintf("%d", id),
		"Username": username,
		"Message":  "successfully accessed protected resource",
	}))
}
