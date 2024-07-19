package core

import (
	"context"

	"golang.org/x/oauth2"
)

type (
	AuthServiceConfig struct {
		JWTSecret      []byte
		VKClientID     string
		VKClientSecret string
		VKCallback     string
	}

	AuthService interface {
		GetJWTSecret() []byte
		Login(ctx context.Context, user User) (accessToken string, refreshToken string, err error)
		LoginVK(ctx context.Context, externalUserID int) (string, string, error)
		Signup(ctx context.Context, user User) error
		Refresh(ctx context.Context, id int) (accessToken string, err error)
		ConfigVK() *oauth2.Config
		GetVKUserID(token string) (int, error)
		GetVKLoginPageURL() string
	}
)
