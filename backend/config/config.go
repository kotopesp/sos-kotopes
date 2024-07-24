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
		TLS
		Auth
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

	TLS struct {
		TLSCert string
		TLSKey  string
	}

	Auth struct {
		JWTSecret      []byte
		VKClientID     string
		VKClientSecret string
		VKCallback     string
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	port := flag.String("port", "localhost:8080", "port")
	logLevel := flag.String("log_level", string(logger.InfoLevel), "logger level")
	dbURL := flag.String("db_url", "", "url for connection to database")
	tlsCert := flag.String("tls_cert", "./tls/cert.pem", "path to tls certificate")
	tlsKey := flag.String("tls_key", "./tls/key.pem", "path to tls key")
	jwtSecret := flag.String("jwt_secret", "secret", "key that used to sign jwt")
	vkClientID := flag.String("vk_client_id", "", "vk id of our app")
	vkClientSecret := flag.String("vk_client_secret", "", "key that used to access vk api")
	vkCallback := flag.String("vk_callback", "https://08ec-102-38-225-73.ngrok-free.app/api/v1/auth/login/vk/callback", "callback for vk auth")

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
		TLS: TLS{
			TLSCert: *tlsCert,
			TLSKey:  *tlsKey,
		},
		Auth: Auth{
			JWTSecret:      []byte(*jwtSecret),
			VKClientID:     *vkClientID,
			VKClientSecret: *vkClientSecret,
			VKCallback:     *vkCallback,
		},
	}

	return cfg, nil
}
