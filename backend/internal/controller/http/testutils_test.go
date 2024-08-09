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
		authService *mocks.AuthService
		postService *mocks.PostService
	}
)

var (
	secret = []byte("secret")
)

func newTestApp(t *testing.T) (*fiber.App, appDependencies) {
	app := fiber.New()
	ctx := context.Background()

	mockAuthService := mocks.NewAuthService(t)
	mockPostService := mocks.NewPostService(t)
	mockRoleService := mocks.NewRoleService(t)
	mockUserService := mocks.NewUserService(t)
	formValidatorService := validator.New(ctx, baseValidator.New())

	mockAuthService.On("GetJWTSecret").Return(secret)

	// mock your dependencies and put them here
	NewRouter(
		app,
		mockAuthService,
		mockUserService,
		mockRoleService,
		formValidatorService,
		mockPostService,
	)

	return app, appDependencies{
		authService: mockAuthService,
		postService: mockPostService,
	}
}
