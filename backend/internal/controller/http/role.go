package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/role"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

//	@Summary		Give role to user
//	@Tags			role
//	@Description	Give role to user
//	@ID				give-role-to-user
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuthBasic
//	@Param			request	body		role.GivenRole	true	"Role"
//	@Success		201		{object}	model.Response{data=role.Role}
//	@Failure		400		{object}	model.Response
//	@Failure		404		{object}	model.Response
//
//	@Failure		422		{object}	model.Response{data=[]validator.ResponseError}
//
//	@Failure		500		{object}	model.Response
//	@Router			/users/roles [post]
func (r *Router) giveRoleToUser(ctx *fiber.Ctx) error {
	id, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
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
	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(addedRole))
}

//	@Summary		Get user roles
//	@Tags			role
//	@Description	Get user roles
//	@ID				get-user-roles
//	@Accept			json
//	@Produce		json
//	@Param			request	path		int	true	"User ID"
//	@Success		200		{object}	model.Response{data=[]role.Role}
//	@Failure		400		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Router			/users/{id}/roles [get]
func (r *Router) getUserRoles(ctx *fiber.Ctx) error {
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
		return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(core.RoleDetails{}))
	}

	roles := make([]role.Role, 0, len(userRoles))
	for i := range userRoles {
		roles = append(roles, role.ToRole(&userRoles[i]))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(roles))
}

//	@Summary		Update user roles
//	@Tags			role
//	@Description	Update user roles
//	@ID				update-user-roles
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuthBasic
//	@Param			request	body		role.UpdateRole	true	"Role"
//	@Success		200		{object}	model.Response{data=role.Role}
//	@Failure		400		{object}	model.Response
//	@Failure		404		{object}	model.Response
//
//	@Failure		422		{object}	model.Response{data=[]validator.ResponseError}
//
//	@Failure		500		{object}	model.Response
//	@Router			/users/roles [patch]
func (r *Router) updateUserRoles(ctx *fiber.Ctx) error {
	id, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
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
		case errors.Is(err, core.ErrUserRoleNotFound):
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

//	@Summary		Delete user role
//	@Tags			role
//	@Description	Delete user role
//	@ID				delete-user-role
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuthBasic
//	@Param			request	body	role.DeleteRole	true	"Role"
//	@Success		204
//	@Failure		400	{object}	model.Response
//
//	@Failure		422	{object}	model.Response{data=[]validator.ResponseError}
//
//	@Failure		500	{object}	model.Response
//	@Router			/users/roles [delete]
func (r *Router) deleteUserRole(ctx *fiber.Ctx) error {
	id, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	var deleteRole role.DeleteRole
	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &deleteRole)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	err = r.roleService.DeleteUserRole(ctx.UserContext(), id, deleteRole.Name)
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

	return ctx.Status(fiber.StatusNoContent).JSON(id)
}
