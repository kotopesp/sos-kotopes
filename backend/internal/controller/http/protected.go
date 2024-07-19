package http

import (
	"fmt"

	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
	"github.com/gofiber/fiber/v2"
)

func (r *Router) protected(ctx *fiber.Ctx) error {
	idItem, _ := getPayloadItem(ctx, "id")
	id := int(idItem.(float64))
	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(fiber.Map{
		"UserID":  fmt.Sprintf("%d", id),
		"Message": "successfully accessed protected resource",
	}))
}
