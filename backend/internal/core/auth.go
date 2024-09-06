package core

import (
	"context"
	"time"

	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type (
	AuthServiceConfig struct {
		JWTSecret            []byte
		VKClientID           string
		VKClientSecret       string
		VKCallback           string
		AccessTokenLifetime  int
		RefreshTokenLifetime int
	}

	AuthService interface {
		GetJWTSecret() []byte
		LoginBasic(ctx context.Context, user User) (accessToken, refreshToken *string, err error)
		SignupBasic(ctx context.Context, user User) error
		Refresh(ctx context.Context, refreshSession RefreshSession) (accessToken, refreshToken *string, err error)
		ConfigVK() *oauth2.Config
		AuthorizeVK(ctx context.Context, token string) (accessToken, refreshToken *string, err error)
	}

	RefreshSession struct {
		ID           int       `gorm:"column:id"`
		UserID       int       `gorm:"column:user_id"`
		RefreshToken string    `gorm:"column:refresh_token"`
		ExpiresAt    time.Time `gorm:"column:expires_at"`
	}

	RefreshSessionStore interface {
		UpdateRefreshSession(
			ctx context.Context,
			param UpdateRefreshSessionParam,
			refreshSession RefreshSession,
		) error
		CountSessionsAndDelete(ctx context.Context, userID int) error
		GetRefreshSessionByToken(ctx context.Context, token string) (data RefreshSession, err error)
	}

	UpdateRefreshSessionParam func(*gorm.DB) error
)

const (
	VKGetUsersURL = "https://api.vk.com/method/users.get"
	VKAPIVersion  = "5.199"
)

func (RefreshSession) TableName() string {
	return "refresh_sessions"
}
