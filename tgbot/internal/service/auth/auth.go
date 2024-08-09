package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kotopesp/tgbot/internal/core"
)

type service struct {
	authServiceConfig core.AuthServiceConfig
}

func New(authServiceConfig core.AuthServiceConfig) core.AuthService {
	return &service{
		authServiceConfig: authServiceConfig,
	}
}

func (s *service) GetAuthURL(id int) (*string, error) {
	token, err := s.generateToken(id)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(
		s.authServiceConfig.TelegramCallback+"?token=%s",
		*token,
	)

	return &url, nil
}

func (s *service) generateToken(id int) (*string, error) {
	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Duration(s.authServiceConfig.TelegramTokenLifetime) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(s.authServiceConfig.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
