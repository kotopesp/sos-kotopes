package http

import (
	"errors"
	"fmt"
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

	return ctx.Status(fiber.StatusOK).JSON(posts)
}

func (r *Router) GiveRoleToUser(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	var givenRole role.GivenRole
	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &givenRole)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}
	coreGivenRole := givenRole.ToCoreGivenRole()

	addedRole, err := r.roleService.GiveRoleToUser(ctx.UserContext(), id, coreGivenRole)
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
	return ctx.Status(fiber.StatusCreated).JSON(addedRole)
}

func (r *Router) DeleteUserRole(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	var givenRole role.GivenRole
	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &givenRole)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}
	coreGivenRole := givenRole.ToCoreGivenRole()

	name := coreGivenRole.Name
	err = r.roleService.DeleteUserRole(ctx.UserContext(), id, name)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrNoSuchUser):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusNoContent).JSON(model.ErrorResponse(err.Error()))
		default:
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
		}
	}
	fmt.Println(err)
	return ctx.Status(fiber.StatusNoContent).JSON(id)
}

func (r *Router) UpdateUserRoles(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}
	var updateRole role.UpdateRole
	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &updateRole)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}
	coreUpdateRole := updateRole.ToCoreUpdateRole()

	updatedRole, err := r.roleService.UpdateUserRole(ctx.UserContext(), id, coreUpdateRole)
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
	modelRole := role.ToRole(&updatedRole)
	return ctx.Status(fiber.StatusOK).JSON(modelRole)
}
