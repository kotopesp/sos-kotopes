package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kotopesp/sos-kotopes/internal/core"
	refreshsession "github.com/kotopesp/sos-kotopes/internal/store/refresh_session"
	"golang.org/x/crypto/bcrypt"
)

const (
	vkAuthProvider = "vk"
	bcryptCost     = 12
)

var authProvidersPasswordPlugs = map[string]string{
	vkAuthProvider: "vk_password",
}

type service struct {
	userStore           core.UserStore
	refreshSessionStore core.RefreshSessionStore
	authServiceConfig   core.AuthServiceConfig
}

func New(
	userStore core.UserStore,
	refreshSessionStore core.RefreshSessionStore,
	authServiceConfig core.AuthServiceConfig,
) core.AuthService {
	return &service{
		userStore:           userStore,
		refreshSessionStore: refreshSessionStore,
		authServiceConfig:   authServiceConfig,
	}
}

// GetJWTSecret need to be accessed from middleware
func (s *service) GetJWTSecret() []byte {
	return s.authServiceConfig.JWTSecret
}

func (s *service) getRefreshTokenExpiresAt() time.Time {
	return time.Now().Add(time.Minute * time.Duration(s.authServiceConfig.RefreshTokenLifetime))
}

func setRefreshSessionData(session *core.RefreshSession, token string, id int, expires time.Time) {
	session.RefreshToken = token
	session.UserID = id
	session.ExpiresAt = expires
}

// LoginBasic Login through username and password
func (s *service) LoginBasic(ctx context.Context, user core.User) (accessToken, refreshToken *string, err error) {
	coreUser, err := s.userStore.GetUserByUsername(ctx, user.Username)
	if err != nil {
		if errors.Is(err, core.ErrNoSuchUser) {
			return nil, nil, core.ErrInvalidCredentials
		}
		return nil, nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(coreUser.PasswordHash), []byte(user.PasswordHash))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, nil, core.ErrInvalidCredentials
		}
		return nil, nil, err
	}

	at, err := s.generateAccessToken(coreUser.ID, coreUser.Username)
	if err != nil {
		return nil, nil, err
	}

	var refreshSession core.RefreshSession
	setRefreshSessionData(&refreshSession, s.generateRefreshToken(), coreUser.ID, s.getRefreshTokenExpiresAt())

	rt := refreshSession.RefreshToken
	refreshSession.RefreshToken = s.generateHash(refreshSession.RefreshToken)

	err = s.refreshSessionStore.CountSessionsAndDelete(ctx, coreUser.ID)
	if err != nil {
		return nil, nil, err
	}

	err = s.refreshSessionStore.UpdateRefreshSession(
		ctx,
		refreshsession.ByNothing(),
		refreshSession,
	)
	if err != nil {
		return nil, nil, err
	}

	return at, &rt, nil
}

// SignupBasic Signup through username and password (can be additional fields)
func (s *service) SignupBasic(ctx context.Context, user core.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcryptCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)

	if _, err := s.userStore.CreateUser(ctx, user); err != nil {
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

	accessToken, refreshToken, err = s.loginVK(ctx, vkUserID)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, err
}

// loginVK Signup if user not exists, then login
func (s *service) loginVK(ctx context.Context, externalUserID int) (accessToken, refreshToken *string, err error) {
	externalUser, err := s.userStore.GetUserByExternalID(ctx, externalUserID)

	var userID int

	if err != nil {
		if errors.Is(err, core.ErrNoSuchUser) {
			userID, err = s.signupVK(ctx, core.User{
				Username:     uuid.New().String(),
				PasswordHash: authProvidersPasswordPlugs[vkAuthProvider],
			}, externalUserID, vkAuthProvider)
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

	if err != nil {
		return nil, nil, err
	}

	return s.LoginBasic(ctx, core.User{
		Username:     user.Username,
		PasswordHash: authProvidersPasswordPlugs[vkAuthProvider],
	})
}

// signupVK Creating external user
func (s *service) signupVK(ctx context.Context, user core.User, externalUserID int, authProvider string) (userID int, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcryptCost)
	if err != nil {
		return 0, err
	}

	user.PasswordHash = string(hashedPassword)
	userID, err = s.userStore.CreateExternalUser(ctx, user, externalUserID, authProvider)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// Refresh Getting new accessToken, when another one expires; need to have refreshToken in cookie
func (s *service) Refresh(
	ctx context.Context,
	refreshSession core.RefreshSession,
) (accessToken, refreshToken *string, err error) {
	hashedToken := s.generateHash(refreshSession.RefreshToken)
	dbSession, err := s.refreshSessionStore.GetRefreshSessionByToken(ctx, hashedToken)
	if err != nil {
		return nil, nil, core.ErrUnauthorized
	}

	user, err := s.userStore.GetUserByID(ctx, dbSession.UserID)
	if err != nil {
		return nil, nil, err
	}

	accessToken, err = s.generateAccessToken(dbSession.UserID, user.Username)
	if err != nil {
		return nil, nil, err
	}

	if dbSession.ExpiresAt.Before(time.Now()) {
		return nil, nil, core.ErrUnauthorized
	}

	setRefreshSessionData(&refreshSession, s.generateRefreshToken(), dbSession.UserID, s.getRefreshTokenExpiresAt())

	rt := refreshSession.RefreshToken
	refreshSession.RefreshToken = s.generateHash(refreshSession.RefreshToken)

	err = s.refreshSessionStore.UpdateRefreshSession(
		ctx,
		refreshsession.ByID(dbSession.ID),
		refreshSession,
	)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, &rt, nil
}
