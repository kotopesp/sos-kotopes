package app

import (
	"log"

	"github.com/kotopesp/tgbot/config"
	v1 "github.com/kotopesp/tgbot/internal/controller/bot"
	"github.com/kotopesp/tgbot/internal/core"
	"github.com/kotopesp/tgbot/internal/service/auth"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func Run(cfg *config.Config) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	bot, err := telego.NewBot(cfg.TelegramToken, telego.WithDefaultDebugLogger())
	if err != nil {
		log.Fatal(err)
	}

	updates, err := bot.UpdatesViaLongPolling(nil)
	if err != nil {
		log.Fatal(err)
	}

	defer bot.StopLongPolling()

	app, err := th.NewBotHandler(bot, updates)
	if err != nil {
		log.Fatal(err)
	}

	// Services
	authService := auth.New(core.AuthServiceConfig{
		JWTSecret:             cfg.JWTSecret,
		TelegramCallback:      cfg.TelegramCallback,
		TelegramTokenLifetime: cfg.TelegramTokenLifetime,
	})

	v1.NewRouter(
		app,
		authService,
	)

	defer app.Stop()
	app.Start()
}
