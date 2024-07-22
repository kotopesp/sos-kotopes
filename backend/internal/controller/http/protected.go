package http

import (
	"fmt"

	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
	"github.com/gofiber/fiber/v2"
)

func (r *Router) protected(ctx *fiber.Ctx) error {
	idItem := getPayloadItem(ctx, "id")
	idFloat, ok := idItem.(float64)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse("error while reading id from token"))
	}
	id := int(idFloat)
	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(fiber.Map{
		"UserID":  fmt.Sprintf("%d", id),
		"Message": "successfully accessed protected resource",
	}))
}
