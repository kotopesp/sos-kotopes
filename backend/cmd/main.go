package main

import (
	"log"

	"github.com/kotopesp/sos-kotopes/config"
	"github.com/kotopesp/sos-kotopes/internal/app"
)

//@title SOS Kotopes API
//@version 1.0
//@description This is the API for the SOS Kotopes project. It provides endpoints for managing the database of animals in need of help.

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.apiKey	ApiKeyAuthBasic
//	@in							header
//	@name						Authorization

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
