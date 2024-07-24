package auth

import (
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
)

const (
	vkPasswordPlug = "vk"
)

type service struct {
	userStore         core.UserStore
	authServiceConfig core.AuthServiceConfig
}

func New(
	userStore core.UserStore,
	authServiceConfig core.AuthServiceConfig,
) core.AuthService {
	return &service{
		userStore:         userStore,
		authServiceConfig: authServiceConfig,
	}
}

// need to be accessed from middleware
func (s *service) GetJWTSecret() []byte {
	return s.authServiceConfig.JWTSecret
}
