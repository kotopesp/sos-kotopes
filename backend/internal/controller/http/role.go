package http

import (
	"fmt"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/role"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (r *Router) GetUserRoles(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	if idStr == "" {
		logger.Log().Debug(ctx.UserContext(), "Error: id is required")
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "id is required",
		})
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}
	userRoles, err := r.roleService.GetUserRoles(ctx.UserContext(), id)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	if len(userRoles) == 0 {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "User has no roles",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(userRoles)
}

func (r *Router) GiveRoleToUser(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	if idStr == "" {
		logger.Log().Debug(ctx.UserContext(), "Error: id is required")
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "id is required",
		})
	}
	id, err := strconv.Atoi(idStr)

	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	var body role.GiveRole
	if err = ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	err = r.roleService.GiveRoleToUser(ctx.UserContext(), id, body)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(id)
}

func (r *Router) DeleteUserRole(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	if idStr == "" {
		logger.Log().Debug(ctx.UserContext(), "Error: id is required")
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "id is required",
		})
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	var roleName role.GiveRole
	if err = ctx.BodyParser(&roleName); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	name := roleName.Name
	err = r.roleService.DeleteUserRole(ctx.UserContext(), id, name)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(id)
}

func (r *Router) UpdateUserRoles(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	if idStr == "" {
		logger.Log().Debug(ctx.UserContext(), "Error: id is required")
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "id is required",
		})
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}
	fmt.Println(id)
	return nil
}
