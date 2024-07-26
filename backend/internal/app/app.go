package app

import (
	"context"
	v1 "gitflic.ru/spbu-se/sos-kotopes/internal/controller/http"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/internal/service/auth"
	"gitflic.ru/spbu-se/sos-kotopes/internal/service/name"
	"gitflic.ru/spbu-se/sos-kotopes/internal/service/role_service"
	"gitflic.ru/spbu-se/sos-kotopes/internal/service/user_service"
	"gitflic.ru/spbu-se/sos-kotopes/internal/store/entity"
	"gitflic.ru/spbu-se/sos-kotopes/internal/store/role_store"
	"gitflic.ru/spbu-se/sos-kotopes/internal/store/user_store"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"os"
	"os/signal"
	"syscall"

	"gitflic.ru/spbu-se/sos-kotopes/config"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/logger"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	ctx := context.Background()

	// Init logger
	logger.New(cfg.Log.Level)

	// Postgres connection
	pg, err := postgres.New(ctx, cfg.DB.URL)
	if err != nil {
		logger.Log().Fatal(ctx, "error with connection to database: %s", err.Error())
	}
	defer pg.Close(ctx)

	// Stores
	entityStore := entity.New(pg)
	userStore := user_store.NewUserStore(pg)
	roleStore := role_store.NewRoleStore(pg)
	// Services
	roleService := role_service.NewRoleService(roleStore)
	entityService := name.New(entityStore)
	userService := user_service.NewUserService(userStore)
	authService := auth.New(
		userStore,
		core.AuthServiceConfig{
			JWTSecret: cfg.JWTSecret,
			//VKClientID:     cfg.VKClientID,
			//VKClientSecret: cfg.VKClientSecret,
			//VKCallback:     cfg.VKCallback,
		},
	)

	// HTTP Server
	app := fiber.New(fiber.Config{
		CaseSensitive:            true,
		StrictRouting:            false,
		EnableSplittingOnParsers: true,
	})
	app.Use(recover.New())
	app.Use(cors.New())

	v1.NewRouter(app, entityService, authService, userService, roleService)

	logger.Log().Info(ctx, "server was started on %s", cfg.HTTP.Port)
	err = app.Listen(cfg.HTTP.Port)
	if err != nil {
		logger.Log().Fatal(ctx, "server was stopped: %s", err.Error())
	}

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Log().Info(ctx, "signal %s received", s.String())
	case <-ctx.Done():
		return
	}

	// Shutdown
	err = app.Shutdown()
	if err != nil {
		logger.Log().Fatal(ctx, "error with shutdown server: %s", err.Error())
	}
}
