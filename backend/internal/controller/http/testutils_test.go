package http

import (
	"context"
	baseValidator "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/validator"
	"github.com/kotopesp/sos-kotopes/internal/core/mocks"
	"testing"
)

type (
	appDependencies struct {
		entityService *mocks.EntityService
		authService   *mocks.AuthService
	}
)

var (
	secret = []byte("secret")
)

func newTestApp(t *testing.T) (*fiber.App, appDependencies) {
	app := fiber.New()
	ctx := context.Background()

	mockEntityService := mocks.NewEntityService(t)
	mockAuthService := mocks.NewAuthService(t)
	formValidatorService := validator.New(ctx, baseValidator.New())

	mockAuthService.On("GetJWTSecret").Return(secret)

	// mock your dependencies and put them here
	NewRouter(
		app,
		mockEntityService,
		formValidatorService,
		mockAuthService,
	)

	return app, appDependencies{
		entityService: mockEntityService,
		authService:   mockAuthService,
	}
}
