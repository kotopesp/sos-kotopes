package core

type (
	AuthServiceConfig struct {
		JWTSecret             []byte
		TelegramCallback      string
		TelegramTokenLifetime int
	}

	AuthService interface {
		GetAuthURL(id int) (url *string, err error)
	}
)
