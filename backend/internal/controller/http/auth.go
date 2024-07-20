package http

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/user"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/validator"
	"gitflic.ru/spbu-se/sos-kotopes/internal/service/auth"
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
			Key:    r.authService.GetJWTSecret(),
		},
		ErrorHandler: authErrorHandler,
	})
}

func (r *Router) refreshTokenMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256,
			Key:    r.authService.GetJWTSecret(),
		},
		ErrorHandler: authErrorHandler,
		TokenLookup:  "cookie:refresh_token",
	})
}

func (r *Router) loginBasic(ctx *fiber.Ctx) error {
	apiUser := user.User{}
	if err := ctx.BodyParser(&apiUser); err != nil {
		logger.Log().Error(ctx.Context(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	errs := validator.Validate(apiUser)
	if len(errs) > 0 {
		logger.Log().Error(ctx.Context(), fmt.Sprintf("%v", errs))
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ErrorResponse(fiber.Map{
			"validation_errors": errs,
		}))
	}

	coreUser := apiUser.ToCoreUser()
	accessToken, refreshToken, err := r.authService.Login(ctx.Context(), coreUser)
	if err != nil {
		logger.Log().Error(ctx.Context(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	setRefreshTokenCookie(ctx, refreshToken)
	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(fiber.Map{
		"access_token": accessToken,
	}))
}

func getPhotoBytes(photo *multipart.FileHeader) (*[]byte, error) {
	file, err := photo.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	photoBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return &photoBytes, nil
}

func (r *Router) signup(ctx *fiber.Ctx) error {
	apiUser := user.User{}
	if err := ctx.BodyParser(&apiUser); err != nil {
		logger.Log().Error(ctx.Context(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}
	photo, err := ctx.FormFile("photo")
	if err != nil {
		apiUser.Photo = nil
	} else {
		photoBytes, err := getPhotoBytes(photo)
		if err != nil {
			logger.Log().Error(ctx.Context(), err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
		}
		apiUser.Photo = photoBytes
	}
	errs := validator.Validate(apiUser)
	if len(errs) > 0 {
		logger.Log().Error(ctx.Context(), fmt.Sprintf("%v", errs))
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ErrorResponse(fiber.Map{
			"validation_errors": errs,
		}))
	}
	coreUser := apiUser.ToCoreUser()

	err = r.authService.Signup(ctx.Context(), coreUser)
	if err == nil {
		return ctx.Status(fiber.StatusOK).JSON(model.OKResponse("user created successfully"))
	}

	logger.Log().Error(ctx.Context(), err.Error())
	if errors.Is(err, auth.ErrNotUniqueUsername) {
		errs = append(errs, model.ErrNotUniqueUsername(coreUser.Username))
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ErrorResponse(fiber.Map{
			"validation_errors": errs,
		}))
	}
	return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
}

func (r *Router) refresh(ctx *fiber.Ctx) error {
	sub, _ := getPayloadItem(ctx, "sub")
	id := int(sub.(float64))

	accessToken, err := r.authService.Refresh(ctx.UserContext(), id)
	if err != nil {
		logger.Log().Error(ctx.Context(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(fiber.Map{
		"access_token": accessToken,
	}))
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

func generateState(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (r *Router) loginVK(ctx *fiber.Ctx) error {
	cfg := r.authService.ConfigVK()

	state, err := generateState(16)
	if err != nil {
		logger.Log().Info(ctx.Context(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "state",
		Value:    state,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	url := cfg.AuthCodeURL(state)
	return ctx.Redirect(url)
}

func setRefreshTokenCookie(ctx *fiber.Ctx, refreshToken string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})
}

func (r *Router) callback(ctx *fiber.Ctx) error {
	token, err := r.authService.ConfigVK().Exchange(ctx.Context(), ctx.FormValue("code"))
	if err != nil {
		logger.Log().Error(ctx.Context(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	if ctx.FormValue("state") != ctx.Cookies("state") {
		err = errors.New("states do not match")
		logger.Log().Error(ctx.Context(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}
	vkUserID, err := r.authService.GetVKUserID(token.AccessToken)
	if err != nil {
		logger.Log().Error(ctx.Context(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	accessToken, refreshToken, err := r.authService.LoginVK(ctx.Context(), vkUserID)
	if err != nil {
		logger.Log().Error(ctx.Context(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	setRefreshTokenCookie(ctx, refreshToken)
	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(fiber.Map{
		"access_token": accessToken,
	}))
}
