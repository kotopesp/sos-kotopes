package http

import "github.com/gofiber/fiber/v2"

func (r *Router) ping(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).SendString("Pong")
}
