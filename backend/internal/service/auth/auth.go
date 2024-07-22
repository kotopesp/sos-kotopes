package auth

import (
	"context"
	"errors"

	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	userStorePkg "gitflic.ru/spbu-se/sos-kotopes/internal/store/user"
	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
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

func (s *service) Login(ctx context.Context, user core.User) (accessToken, refreshToken string, err error) {
	dbUser, err := s.userStore.GetUserByUsername(ctx, user.Username)
	if err != nil {
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(user.PasswordHash))
	if err != nil {
		return "", "", err
	}

	at, err := s.generateAccessToken(dbUser.ID, dbUser.Username)
	if err != nil {
		return "", "", err
	}

	rt, err := s.generateRefreshToken(dbUser.ID)
	if err != nil {
		return "", "", err
	}

	return at, rt, nil
}

func (s *service) Signup(ctx context.Context, user core.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), 12)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)
	if _, err := s.userStore.AddUser(ctx, user); err != nil {
		if errors.Is(err, userStorePkg.ErrNotUniqueUsername) {
			return ErrNotUniqueUsername
		}
		return err
	}
	return nil
}

func (s *service) LoginVK(ctx context.Context, externalUserID int) (accessToken, refreshToken string, err error) {
	user, err := s.userStore.GetUserByExternalID(ctx, externalUserID)
	user.ExternalID = &externalUserID
	user.PasswordHash = vkPasswordPlug
	if err != nil {
		if errors.Is(err, userStorePkg.ErrNoSuchUser) {
			user.Username = uuid.New().String()
			err = s.Signup(ctx, user)
			if err != nil {
				return "", "", err
			}
		} else {
			return "", "", err
		}
	}
	return s.Login(ctx, core.User{
		Username:     user.Username,
		PasswordHash: vkPasswordPlug,
	})
}

func (s *service) Refresh(ctx context.Context, id int) (accessToken string, err error) {
	user, err := s.userStore.GetUserByID(ctx, id)
	if err != nil {
		return "", err
	}

	accessToken, err = s.generateAccessToken(id, user.Username)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
