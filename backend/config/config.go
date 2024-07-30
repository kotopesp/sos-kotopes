package config

import (
	"flag"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/logger"
)

type (
	Config struct {
		HTTP
		Log
		DB
	}

	HTTP struct {
		Port string
	}

	Log struct {
		Level string
	}

	DB struct {
		URL string
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	port := flag.String("port", "localhost:8080", "port")
	logLevel := flag.String("log_level", string(logger.InfoLevel), "logger level")
	dbURL := flag.String("db_url", "", "url for connection to database")

	flag.Parse()

	cfg := &Config{
		HTTP: HTTP{
			Port: *port,
		},
		Log: Log{
			Level: *logLevel,
		},
		DB: DB{
			URL: *dbURL,
		},
	}

	return cfg, nil
}
