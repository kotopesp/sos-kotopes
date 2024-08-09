package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

const (
	vkAuthProvider       = "vk"
	telegramAuthProvider = "telegram"
)

var authProvidersPasswordPlugs = map[string]string{
	vkAuthProvider:       "vk_password",
	telegramAuthProvider: "telegram_password",
}

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

func (s *service) GetTelegramAuthBotURL() string {
	return s.authServiceConfig.TelegramAuthBotURL
}

// GetJWTSecret need to be accessed from middleware
func (s *service) GetJWTSecret() []byte {
	return s.authServiceConfig.JWTSecret
}

// LoginBasic Login through username and password
func (s *service) LoginBasic(ctx context.Context, user core.User) (accessToken, refreshToken *string, err error) {
	dbUser, err := s.userStore.GetUserByUsername(ctx, user.Username)
	if err != nil {
		if errors.Is(err, core.ErrNoSuchUser) {
			return nil, nil, core.ErrInvalidCredentials
		}
		return nil, nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(user.PasswordHash))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, nil, core.ErrInvalidCredentials
		}
		return nil, nil, err
	}

	at, err := s.generateAccessToken(dbUser.ID, dbUser.Username)
	if err != nil {
		return nil, nil, err
	}

	rt, err := s.generateRefreshToken(dbUser.ID)
	if err != nil {
		return nil, nil, err
	}

	return at, rt, nil
}

// SignupBasic Signup through username and password (can be additional fields)
func (s *service) SignupBasic(ctx context.Context, user core.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), 12)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)

	if _, err := s.userStore.AddUser(ctx, user); err != nil {
		if errors.Is(err, core.ErrNotUniqueUsername) {
			return core.ErrNotUniqueUsername
		}
		return err
	}

	return nil
}

// AuthorizeVK Authorization through VK (user automatically signs up if not exists): getting info from vk, signup, login
func (s *service) AuthorizeVK(ctx context.Context, token string) (accessToken, refreshToken *string, err error) {
	vkUserID, err := s.getVKUserID(token)
	if err != nil {
		return nil, nil, err
	}

	accessToken, refreshToken, err = s.loginExternal(ctx, vkUserID, vkAuthProvider)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, err
}

func (s *service) AuthorizeTelegram(ctx context.Context, telegramUserID int) (accessToken, refreshToken *string, err error) {
	accessToken, refreshToken, err = s.loginExternal(ctx, telegramUserID, telegramAuthProvider)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, err
}

// loginVK Signup if user not exists, then login
func (s *service) loginExternal(ctx context.Context,
	externalUserID int,
	authProvider string) (accessToken, refreshToken *string, err error) {
	externalUser, err := s.userStore.GetUserByExternalID(ctx, externalUserID)

	var userID int

	if err != nil {
		if errors.Is(err, core.ErrNoSuchUser) {
			userID, err = s.signupExternal(ctx, core.User{
				Username:     uuid.New().String(),
				PasswordHash: authProvidersPasswordPlugs[authProvider],
			}, externalUserID, authProvider)
			if err != nil {
				return nil, nil, err
			}
		} else {
			return nil, nil, err
		}
	} else {
		userID = externalUser.UserID
	}

	user, err := s.userStore.GetUserByID(ctx, userID)

	logger.Log().Debug(ctx, fmt.Sprint(err))

	if err != nil {
		return nil, nil, err
	}

	return s.LoginBasic(ctx, core.User{
		Username:     user.Username,
		PasswordHash: authProvidersPasswordPlugs[authProvider],
	})
}

// signupVK Creating external user
func (s *service) signupExternal(ctx context.Context, user core.User, externalUserID int, authProvider string) (userID int, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), 12)
	if err != nil {
		return 0, err
	}

	user.PasswordHash = string(hashedPassword)
	userID, err = s.userStore.AddExternalUser(ctx, user, externalUserID, authProvider)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// Refresh Getting new accessToken, when another one expires; need to have refreshToken in cookie
func (s *service) Refresh(ctx context.Context, id int) (accessToken *string, err error) {
	user, err := s.userStore.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	accessToken, err = s.generateAccessToken(id, user.Username)
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}
