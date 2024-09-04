package http

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"mime/multipart"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/validator"
	"github.com/kotopesp/sos-kotopes/internal/core"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/user"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

// authErrorHandler Error handler if user is not authorized
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

// loginBasic Login through username and password
//
//	@Summary		Login through username and password
//	@Tags			auth
//	@Description	Login through username and password
//	@ID				login-basic
//	@Accept			json
//	@Produce		json
//	@Param			request	body		user.Login	true	"User"
//	@Success		200		{object}	model.Response
//	@Failure		400		{object}	model.Response
//	@Failure		401		{object}	model.Response
//	@Failure		500		{object}	model.Response
//	@Router			/auth/login [post]
func (r *Router) loginBasic(ctx *fiber.Ctx) error {
	var apiUser user.Login
	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &apiUser)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	coreUser := apiUser.ToCoreUser()

	accessToken, refreshToken, err := r.authService.LoginBasic(ctx.Context(), coreUser)
	if err != nil {
		if errors.Is(err, core.ErrInvalidCredentials) {
			logger.Log().Info(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
		}

		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	setRefreshTokenCookie(ctx, *refreshToken)
	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(accessToken))
}

// getPhotoBytes Getting photo from request body
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

// signup Signup through username and password (user can have additional field like photo description)
//
//	@Summary		Signup through username and password
//	@Tags			auth
//	@Description	Signup through username and password
//	@ID				signup
//	@Accept			json
//	@Produce		json
//	@Param			username	formData	string	true	"Username"
//	@Param			password	formData	string	true	"Password"
//	@Param			firstname	formData	string	false	"Firstname"
//	@Param			lastname	formData	string	false	"Lastname"
//	@Param			description	formData	string	false	"Description"
//	@Param			photo		formData	file	false	"Photo"
//	@Success		201			{object}	any
//	@Failure		400			{object}	model.Response
//	@Failure		422			{object}	model.Response
//	@Failure		500			{object}	model.Response
//	@Router			/auth/signup [post]
func (r *Router) signup(ctx *fiber.Ctx) error {
	var apiUser user.User
	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &apiUser)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	photo, err := ctx.FormFile("photo")
	if err != nil {
		apiUser.Photo = nil
	} else {
		photoBytes, err := getPhotoBytes(photo)
		if err != nil {
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
		}
		apiUser.Photo = photoBytes
	}

	coreUser := apiUser.ToCoreUser()

	err = r.authService.SignupBasic(ctx.UserContext(), coreUser)
	if err == nil {
		return ctx.SendStatus(fiber.StatusCreated)
	}

	if errors.Is(err, core.ErrNotUniqueUsername) {
		logger.Log().Info(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ErrorResponse(fiber.Map{
			"validation_errors": []validator.ResponseError{model.ErrNotUniqueUsername(coreUser.Username)},
		}))
	}

	logger.Log().Error(ctx.UserContext(), err.Error())

	return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
}

// getting items from token payload
func getPayloadItem(ctx *fiber.Ctx, key string) any {
	token, ok := ctx.Locals("user").(*jwt.Token)
	if !ok {
		return nil
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}
	return claims[key]
}

// refresh Refreshing token if one expires
//
//	@Summary		Refresh token
//	@Tags			auth
//	@Description	Refresh token
//	@ID				refresh-token
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	model.Response
//	@Failure		401	{object}	model.Response
//	@Failure		500	{object}	model.Response
//	@Security		ApiKeyAuthCookie
//	@Router			/auth/token/refresh [post]
func (r *Router) refresh(ctx *fiber.Ctx) error {
	id, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(model.ErrInvalidTokenID.Error()))
	}

	accessToken, err := r.authService.Refresh(ctx.UserContext(), id)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(accessToken))
}

// generateState State generator to protect from CSRF
func generateState(length int) (*string, error) {
	b := make([]byte, length)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	state := base64.URLEncoding.EncodeToString(b)

	return &state, nil
}

// loginVK Login through VK
func (r *Router) loginVK(ctx *fiber.Ctx) error {
	cfg := r.authService.ConfigVK()

	state, err := generateState(16)
	if err != nil {
		logger.Log().Info(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "state",
		Value:    *state,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	url := cfg.AuthCodeURL(*state)

	return ctx.Redirect(url)
}

// setRefreshTokenCookie Setting refresh token in cookie
func setRefreshTokenCookie(ctx *fiber.Ctx, refreshToken string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})
}

// callback is invoked when the user login through VK
func (r *Router) callback(ctx *fiber.Ctx) error {
	token, err := r.authService.ConfigVK().Exchange(ctx.Context(), ctx.FormValue("code"))
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	if ctx.FormValue("state") != ctx.Cookies("state") {
		logger.Log().Error(ctx.UserContext(), model.ErrOAuthStateMismatch.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrOAuthStateMismatch.Error())
	}

	accessToken, refreshToken, err := r.authService.AuthorizeVK(ctx.Context(), token.AccessToken)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	setRefreshTokenCookie(ctx, *refreshToken)

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(accessToken))
}
