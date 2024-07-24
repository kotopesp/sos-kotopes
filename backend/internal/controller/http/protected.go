package http

import (
	"fmt"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
	"github.com/gofiber/fiber/v2"
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
