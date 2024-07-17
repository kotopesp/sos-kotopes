package auth

import (
	"context"
	"errors"

	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	userwithrolesStore "gitflic.ru/spbu-se/sos-kotopes/internal/store/user_with_roles"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userStore          core.UserStore
	userWithRolesStore core.UserWithRolesStore
}

func New(
	userStore core.UserStore,
	userWithRolesStore core.UserWithRolesStore,
) core.AuthService {
	return &service{
		userStore:          userStore,
		userWithRolesStore: userWithRolesStore,
	}
}

// return values (`access token`, `refresh token`, error)
func (s *service) Login(ctx context.Context, user core.User) (string, string, error) {
	dbUser, err := s.userStore.GetUserByUsername(ctx, user.Username)
	if err != nil {
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return "", "", err
	}

	at, err := generateAccessToken(dbUser.ID, dbUser.Username)
	if err != nil {
		return "", "", err
	}

	rt, err := generateRefreshToken(dbUser.ID)
	if err != nil {
		return "", "", err
	}

	return at, rt, nil
}

func (s *service) Signup(ctx context.Context, userWithRoles core.UserWithRoles) error {
	user := userWithRoles.User
	if user == nil {
		return ErrUserIsNil
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userWithRoles.User.Password), 12)
	if err != nil {
		return err
	}
	userWithRoles.User.Password = string(hashedPassword)
	if err := s.userWithRolesStore.AddUserWithRoles(ctx, userWithRoles); err != nil {
		if errors.Is(err, userwithrolesStore.ErrNotUniqueUsername) {
			return ErrNotUniqueUsername
		}
		return err
	}
	return nil
}

func (s *service) Refresh(ctx context.Context, id int) (accessToken string, err error) {
	user, err := s.userStore.GetUserByID(ctx, id)
	if err != nil {
		return "", err
	}

	accessToken, err = generateAccessToken(id, user.Username)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
