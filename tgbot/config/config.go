package config

import "flag"

type (
	Config struct {
		Auth
		Bot
	}
	Auth struct {
		JWTSecret             []byte
		TelegramCallback      string
		TelegramTokenLifetime int
	}
	Bot struct {
		TelegramToken string
	}
)

func NewConfig() (*Config, error) {
	jwtSecret := flag.String("jwt_secret", "secret", "key that used to sign jwt")
	telegramCallback := flag.String("telegram_callback", "https://bad9-102-38-253-158.ngrok-free.app/api/v1/auth/login/telegram/callback", "callback for telegram auth")
	telegramToken := flag.String("telegram_token", "", "token to access bot")
	telegramTokenLifetime := flag.Int("telegram_token_lifetime", 2, "telegram token lifetime in minutes")

	flag.Parse()

	cfg := &Config{
		Auth: Auth{
			JWTSecret:             []byte(*jwtSecret),
			TelegramCallback:      *telegramCallback,
			TelegramTokenLifetime: *telegramTokenLifetime,
		},
		Bot: Bot{
			TelegramToken: *telegramToken,
		},
	}

	return cfg, nil
}
