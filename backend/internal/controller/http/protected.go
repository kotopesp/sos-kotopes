package http

import (
	"fmt"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
	"github.com/gofiber/fiber/v2"
)

func (r *Router) protected(ctx *fiber.Ctx) error {
	// getting id from token
	idItem := getPayloadItem(ctx, "id")

	idFloat, ok := idItem.(float64)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse("error while reading id from token"))
	}

	id := int(idFloat)

	// getting username from token
	usernameItem := getPayloadItem(ctx, "username")

	usernameString, ok := usernameItem.(string)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse("error while reading username from token"))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(fiber.Map{
		"UserID":   fmt.Sprintf("%d", id),
		"Username": usernameString,
		"Message":  "successfully accessed protected resource",
	}))
}
