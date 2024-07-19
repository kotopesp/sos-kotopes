package http

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (r *Router) getKeepers(c *fiber.Ctx) error {
	keepers, err := r.keepers_service.GetKeepers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// TODO имеет смысл реализовать возвращение не keepers
	// TODO а responseOK
	return c.Status(fiber.StatusOK).JSON(keepers)
}

func (r *Router) getKeeperByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid keeper ID"})
	}

	keeper, err := r.keepersService.GetKeeperByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Keeper not found"})
	}

	return c.JSON(keeper)
}

func (r *Router) createKeeper(c *fiber.Ctx) error {
	var keeper core.Keeper
	if err := c.BodyParser(&keeper); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request body"})
	}

	if err := r.keepersService.CreateKeeper(&keeper); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(keeper)
}

func (r *Router) updateKeeper(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid keeper ID"})
	}

	var keeper core.Keeper
	if err := c.BodyParser(&keeper); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request body"})
	}

	if err := r.keepersService.UpdateKeeper(id, &keeper); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(keeper)
}

func (r *Router) deleteKeeper(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid keeper ID"})
	}

	if err := r.keepersService.DeleteKeeper(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
