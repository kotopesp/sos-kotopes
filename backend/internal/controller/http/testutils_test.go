package http

import (
	"context"
	"testing"

	baseValidator "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/validator"
	websocketmanager "github.com/kotopesp/sos-kotopes/internal/controller/http/websockets"
	mocks "github.com/kotopesp/sos-kotopes/internal/core/mocks"
)

type (
	appDependencies struct {
		authService    *mocks.MockAuthService
		postService    *mocks.MockPostService
		commentService *mocks.MockCommentService
	}
)

var (
	secret = []byte("secret")
)

func newTestApp(t *testing.T) (*fiber.App, appDependencies) {
	app := fiber.New()
	ctx := context.Background()

	mockAuthService := mocks.NewMockAuthService(t)
	mockPostService := mocks.NewMockPostService(t)
	mockCommentService := mocks.NewMockCommentService(t)
	mockRoleService := mocks.NewMockRoleService(t)
	mockUserService := mocks.NewMockUserService(t)
	mockChatService := mocks.NewMockChatService(t)
	mockChatMemberService := mocks.NewMockChatMemberService(t)
	mockMessageService := mocks.NewMockMessageService(t)
	formValidatorService := validator.New(ctx, baseValidator.New())
	webSocketManager := websocketmanager.NewWebSocketManager()

	mockAuthService.On("GetJWTSecret").Return(secret)

	// mock your dependencies and put them here
	NewRouter(
		app,
		mockAuthService,
		mockCommentService,
		mockPostService,
		mockUserService,
		mockRoleService,
		mockChatService,
		mockMessageService,
		mockChatMemberService,
		formValidatorService,
		webSocketManager,
	)

	return app, appDependencies{
		authService:    mockAuthService,
		postService:    mockPostService,
		commentService: mockCommentService,
	}
}
