package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
	"gorm.io/gorm"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.UserStore {
	return &store{pg}
}

func (s *store) GetUser(ctx context.Context, id int) (user core.User, err error) {
	err = s.DB.WithContext(ctx).
		Table("users").
		Where("id = ?", id).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return core.User{}, core.ErrNoSuchUser
		}
		logger.Log().Error(ctx, err.Error())
		return core.User{}, err
	}
	return user, nil
}

func (s *store) UpdateUser(ctx context.Context, id int, update core.UpdateUser) (user core.User, err error) {
	tx := s.DB.WithContext(ctx).Begin()

	if tx.Error != nil {
		err = tx.Error
		logger.Log().Error(ctx, err.Error())
		tx.Rollback()
		return core.User{}, err
	}

	defer func() {
		if r := recover(); r != nil || err != nil {
			logger.Log().Error(ctx, err.Error())
			tx.Rollback()
		}
	}()

	updates := make(map[string]interface{})
	if update.Username != nil {
		updates["username"] = *update.Username
	}
	if update.Firstname != nil {
		updates["firstname"] = *update.Firstname
	}
	if update.Lastname != nil {
		updates["lastname"] = *update.Lastname
	}
	if update.Description != nil {
		updates["description"] = *update.Description
	}
	if update.Photo != nil {
		updates["photo"] = *update.Photo
	}
	// maybe delete
	if update.PasswordHash != nil {
		updates["password_hash"] = *update.PasswordHash
	}

	if len(updates) == 0 {
		logger.Log().Error(ctx, core.ErrEmptyUpdateRequest.Error())
		return core.User{}, core.ErrEmptyUpdateRequest
	} else {
		updates["updated_at"] = time.Now()
	}

	err = tx.WithContext(ctx).Table("users").Where("id = ?", id).Updates(updates).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return core.User{}, core.ErrNoSuchUser
		}
		logger.Log().Error(ctx, err.Error())
		return core.User{}, err
	}

	err = tx.WithContext(ctx).Table("users").Where("id = ?", id).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return core.User{}, core.ErrNoSuchUser
		}
		logger.Log().Error(ctx, err.Error())
		return core.User{}, err
	}

	return user, tx.Commit().Error
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
