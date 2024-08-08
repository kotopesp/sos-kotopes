package user

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"strings"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.UserStore {
	return &store{pg}
}

func (s *store) GetUserByUsername(ctx context.Context, username string) (data core.User, err error) {
	var user core.User
	err = s.DB.WithContext(ctx).First(&user, "username=?", username).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, core.ErrNoSuchUser
	}
	return user, err
}

func (s *store) GetUserByID(ctx context.Context, id int) (data core.User, err error) {
	var user core.User
	err = s.DB.WithContext(ctx).First(&user, "id=?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, core.ErrNoSuchUser
	}
	return user, err
}

func (s *store) AddUser(ctx context.Context, user core.User) (userID int, err error) {
	err = s.DB.WithContext(ctx).Create(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), "users_username_key") { // here I need to somehow catch the error of unique constraint violation
			return 0, core.ErrNotUniqueUsername
		}
	}
	return user.ID, err
}

func (s *store) GetUserByExternalID(ctx context.Context, externalID int) (data core.ExternalUser, err error) {
	var user core.ExternalUser
	err = s.DB.WithContext(ctx).First(&user, "external_id=?", externalID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, core.ErrNoSuchUser
	}
	return user, err
}

func (s *store) AddExternalUser(ctx context.Context, user core.User, externalUserID int, authProvider string) (userID int, err error) {
	tx := s.DB.WithContext(ctx).Begin()

	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	err = tx.Create(&user).Error
	if err != nil {
		return 0, err
	}

	err = tx.Create(&core.ExternalUser{
		ExternalID:   externalUserID,
		UserID:       user.ID,
		AuthProvider: authProvider,
	}).Error
	if err != nil {
		return 0, err
	}

	return user.ID, tx.Commit().Error
}
