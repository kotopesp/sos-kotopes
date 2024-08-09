package main

import (
	"log"

	"github.com/kotopesp/tgbot/config"
	"github.com/kotopesp/tgbot/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
