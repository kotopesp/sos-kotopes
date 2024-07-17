package http

import (
	"errors"
	"fmt"

	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/user"
	userwithroles "gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/user_with_roles"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/validator"
	auth "gitflic.ru/spbu-se/sos-kotopes/internal/service/auth_jwt"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/logger"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func authErrorHandler(ctx *fiber.Ctx, err error) error {
	return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
}

func (r *Router) protectedMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256,
			Key:    auth.Secret,
		},
		ErrorHandler: authErrorHandler,
	})
}

func (r *Router) refreshTokenMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256,
			Key:    auth.Secret,
		},
		ErrorHandler: authErrorHandler,
		TokenLookup:  "cookie:refresh_token",
	})
}

func (r *Router) login(ctx *fiber.Ctx) error {
	apiUser := user.User{}
	if err := ctx.BodyParser(&apiUser); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	errs := validator.Validate(apiUser)
	if len(errs) > 0 {
		logger.Log().Error(ctx.UserContext(), fmt.Sprintf("%v", errs))
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ErrorResponse(errs))
	}

	coreUser := apiUser.ToCoreUser()
	accessToken, refreshToken, err := r.authService.Login(ctx.UserContext(), *coreUser)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		// `Secure` enable for https
		// Secure:   true,
		SameSite: "Strict",
	})
	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(accessToken))
}

func (r *Router) signup(ctx *fiber.Ctx) error {
	apiUserWithRoles := userwithroles.UserWithRoles{}
	if err := ctx.BodyParser(&apiUserWithRoles); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}
	errs := validator.Validate(apiUserWithRoles)
	if len(errs) > 0 {
		logger.Log().Error(ctx.UserContext(), fmt.Sprintf("%v", errs))
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ErrorResponse(errs))
	}
	coreUserWithRoles := apiUserWithRoles.ToCoreUserWithRoles()

	err := r.authService.Signup(ctx.UserContext(), *coreUserWithRoles)
	if err == nil {
		return ctx.Status(fiber.StatusOK).JSON(model.OKResponse("user created successfully"))
	}

	logger.Log().Error(ctx.UserContext(), err.Error())
	switch {
	case errors.Is(err, auth.ErrNotUniqueUsername):
		errNotUniqueUsername := model.ErrNotUniqueUsername
		errNotUniqueUsername.Value = coreUserWithRoles.User.Username
		errs = append(errs, errNotUniqueUsername)
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ErrorResponse(errs))
	case errors.Is(err, auth.ErrUserIsNil):
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
}

func (r *Router) refresh(ctx *fiber.Ctx) error {
	sub, _ := getPayloadItem(ctx, "sub")
	id := int(sub.(float64))

	accessToken, err := r.authService.Refresh(ctx.UserContext(), id)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(accessToken))
}

// getting items from token payload
func getPayloadItem(ctx *fiber.Ctx, key string) (any, bool) {
	user, ok := ctx.Locals("user").(*jwt.Token)
	if !ok {
		return nil, false
	}
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return nil, false
	}
	return claims[key], true
}
