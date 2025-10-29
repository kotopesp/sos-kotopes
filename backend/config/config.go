package config

import (
	"flag"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

type (
	Config struct {
		HTTP
		Log
		DB
		TLS
		Auth
		CORS
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
		JWTSecret            []byte
		VKClientID           string
		VKClientSecret       string
		VKCallback           string
		AccessTokenLifetime  int
		RefreshTokenLifetime int
	}

	CORS struct {
		AllowedOrigins string
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	port := flag.String("port", "localhost:8080", "port")
	logLevel := flag.String("log_level", string(logger.InfoLevel), "logger level")
	dbURL := flag.String("db_url", "host=localhost port=5432 user=postgres password=Dabzelosiqqq123 ", "url for connection to database")
	tlsCert := flag.String("tls_cert", "./tls/cert.pem", "path to tls certificate")
	tlsKey := flag.String("tls_key", "./tls/key.pem", "path to tls key")
	jwtSecret := flag.String("jwt_secret", "secret", "key that used to sign jwt")
	vkClientID := flag.String("vk_client_id", "", "vk id of our app")
	vkClientSecret := flag.String("vk_client_secret", "", "key that used to access vk api")
	vkCallback := flag.String("vk_callback", "https://59bf-91-223-89-38.ngrok-free.app/api/v1/auth/login/vk/callback", "callback for vk auth")
	accessTokenLifetime := flag.Int("access_token_lifetime", 10, "access token lifetime in minutes")
	refreshTokenLifetime := flag.Int("refresh_token_lifetime", 43800, "refresh token lifetime in minutes")
	frontendOrigins := flag.String("frontend_allowed_origins", "http://localhost:4200", "comma separated list of origins for the frontend")

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
			JWTSecret:            []byte(*jwtSecret),
			VKClientID:           *vkClientID,
			VKClientSecret:       *vkClientSecret,
			VKCallback:           *vkCallback,
			AccessTokenLifetime:  *accessTokenLifetime,
			RefreshTokenLifetime: *refreshTokenLifetime,
		},
		CORS: CORS{
			AllowedOrigins: *frontendOrigins,
		},
	}

	return cfg, nil
}
