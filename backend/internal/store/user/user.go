package user

import (
	"context"
	"errors"
	"strings"

	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/postgres"
	"gorm.io/gorm"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.UserStore {
	return &store{pg}
}

func (s *store) GetUserByUsername(ctx context.Context, username string) (data core.User, err error) {
	user := core.User{}
	result := s.DB.WithContext(ctx).First(&user, "username=?", username)
	data = user
	err = result.Error
	return
}

func (s *store) GetUserByID(ctx context.Context, id int) (data core.User, err error) {
	user := core.User{}
	result := s.DB.WithContext(ctx).First(&user, id)
	data = user
	err = result.Error
	return
}

func (s *store) AddUser(ctx context.Context, user core.User) (int, error) {
	err := s.DB.WithContext(ctx).Create(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), "users_username_key") { // here I need to somehow catch the error of unique constraint violation
			return 0, ErrNotUniqueUsername
		}
	}
	return user.ID, err
}

func (s *store) GetUserByExternalID(ctx context.Context, extID int) (core.User, error) {
	user := core.User{}
	err := s.DB.WithContext(ctx).First(&user, "ext_id=?", extID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, ErrNoSuchUser
		}
		return user, err
	}
	return user, nil
}
