package http

import (
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
	"github.com/gofiber/fiber/v2"
)

func (r *Router) protected(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse("successfully accessed protected resource"))
}
