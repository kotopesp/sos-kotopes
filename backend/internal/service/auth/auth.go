package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kotopesp/sos-kotopes/internal/core"

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
	return time.Now().Add(time.Hour * time.Duration(s.authServiceConfig.RefreshTokenLifetime))
}

func setUserData(session *core.RefreshSession, token string, id int, expires time.Time) {
	session.RefreshToken = token
	session.UserID = id
	session.ExpiresAt = expires
}

// LoginBasic Login through username and password
func (s *service) LoginBasic(ctx context.Context, user core.User, refreshSession core.RefreshSession) (accessToken, refreshToken *string, err error) {
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

	setUserData(&refreshSession, s.generateRefreshToken(), dbUser.ID, s.getRefreshTokenExpiresAt())

	fingerprintHash, err := bcrypt.GenerateFromPassword([]byte(refreshSession.FingerprintHash), bcryptCost)
	if err != nil {
		return nil, nil, err
	}
	refreshSession.FingerprintHash = string(fingerprintHash)

	err = s.refreshSessionStore.CreateRefreshSession(ctx, refreshSession)
	if err != nil {
		return nil, nil, err
	}

	return at, &refreshSession.RefreshToken, nil
}

// SignupBasic Signup through username and password (can be additional fields)
func (s *service) SignupBasic(ctx context.Context, user core.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcryptCost)
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
	}, core.RefreshSession{})
}

// signupVK Creating external user
func (s *service) signupVK(ctx context.Context, user core.User, externalUserID int, authProvider string) (userID int, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcryptCost)
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
func (s *service) Refresh(ctx context.Context, rs core.RefreshSession) (accessToken, refreshToken *string, err error) {
	// check if refresh token is valid
	session, err := s.refreshSessionStore.GetRefreshSessionByToken(ctx, rs.RefreshToken)
	if err != nil {
		return nil, nil, core.ErrUnauthorized
	}

	user, err := s.userStore.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, nil, err
	}

	accessToken, err = s.generateAccessToken(rs.UserID, user.Username)
	if err != nil {
		return nil, nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(session.FingerprintHash), []byte(rs.FingerprintHash))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, nil, core.ErrUnauthorized
		}
		return nil, nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, nil, core.ErrUnauthorized
	}

	setUserData(&rs, rs.RefreshToken, rs.UserID, s.getRefreshTokenExpiresAt())

	err = s.refreshSessionStore.UpdateRefreshSession(ctx, session.ID, rs)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, &rs.RefreshToken, nil
}
