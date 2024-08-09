package core

import (
	"context"

	"golang.org/x/oauth2"
)

type (
	AuthServiceConfig struct {
		JWTSecret            []byte
		VKClientID           string
		VKClientSecret       string
		VKCallback           string
		AccessTokenLifetime  int
		RefreshTokenLifetime int
		TelegramAuthBotURL   string
	}

	AuthService interface {
		GetJWTSecret() []byte
		LoginBasic(ctx context.Context, user User) (accessToken, refreshToken *string, err error)
		SignupBasic(ctx context.Context, user User) error
		Refresh(ctx context.Context, id int) (accessToken *string, err error)
		ConfigVK() *oauth2.Config
		AuthorizeVK(ctx context.Context, token string) (accessToken, refreshToken *string, err error)
		AuthorizeTelegram(ctx context.Context, telegramUserID int) (accessToken, refreshToken *string, err error)
		GetTelegramAuthBotURL() string
	}
)

const (
	VKGetUsersURL = "https://api.vk.com/method/users.get"
	VKAPIVersion  = "5.199"
)
