package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	v1 "gitflic.ru/spbu-se/sos-kotopes/internal/controller/http"
	chatservice "gitflic.ru/spbu-se/sos-kotopes/internal/service/chat"
	"gitflic.ru/spbu-se/sos-kotopes/internal/service/name"
	chatstore "gitflic.ru/spbu-se/sos-kotopes/internal/store/chat"
	"gitflic.ru/spbu-se/sos-kotopes/internal/store/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

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
	chatStore := chatstore.New(pg)
	// Services
	entityService := name.New(entityStore)
	chatService := chatservice.New(chatStore)

	// HTTP Server
	app := fiber.New(fiber.Config{
		CaseSensitive:            true,
		StrictRouting:            false,
		EnableSplittingOnParsers: true,
	})
	app.Use(recover.New())
	app.Use(cors.New())

	v1.NewRouter(app, entityService, nil, chatService)

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
