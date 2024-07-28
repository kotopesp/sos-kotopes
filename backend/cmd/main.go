package main

import (
	"log"

	"github.com/kotopesp/sos-kotopes/config"
	"github.com/kotopesp/sos-kotopes/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
