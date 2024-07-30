package http

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/role"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

func (r *Router) GetUserRoles(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	userRoles, err := r.roleService.GetUserRoles(ctx.UserContext(), id)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	if len(userRoles) == 0 {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "User has no roles",
		})
	}

	posts := make([]role.Role, 0, len(userRoles))
	for i := range userRoles {
		posts = append(posts, role.ToRole(&userRoles[i]))
	}

	return ctx.Status(fiber.StatusOK).JSON(userRoles)
}

func (r *Router) GiveRoleToUser(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	var body core.GiveRole // fix
	if err = ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	err = r.roleService.GiveRoleToUser(ctx.UserContext(), id, body)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrNoSuchUser):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		default:
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(id)
}

func (r *Router) DeleteUserRole(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
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
		switch {
		case errors.Is(err, core.ErrNoSuchUser):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		default:
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(id)
}

func (r *Router) UpdateUserRoles(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}
	var updateRole core.UpdateRole //fix
	if err = ctx.BodyParser(&updateRole); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	err = r.roleService.UpdateUserRole(ctx.UserContext(), id, updateRole)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrNoSuchUser):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		default:
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(id)
}
