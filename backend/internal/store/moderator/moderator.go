package moderator

import (
	"context"
	"errors"
	"fmt"
	"github.com/lib/pq"

	"gorm.io/gorm"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.ModeratorStore {
	return &store{pg}
}

// GetModerator - получение структуры модератора по его id.
func (s *store) GetModerator(ctx context.Context, id int) (moderator core.Moderator, err error) {
	err = s.DB.WithContext(ctx).
		Table(moderator.TableName()).
		Where("id = ?", id).
		First(&moderator).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return moderator, core.ErrNoSuchModerator
		}
		logger.Log().Debug(ctx, err.Error())
		return moderator, err
	}

	logger.Log().Info(ctx, fmt.Sprintf("%+v", moderator))

	return moderator, nil
}

// GetModeratorByUsername - получение структуры модератора по юзернейму.
func (s *store) GetModeratorByUsername(ctx context.Context, username string) (moderator core.Moderator, err error) {
	err = s.DB.WithContext(ctx).
		Table(moderator.TableName()).
		Where("username = ?", username).
		First(&moderator).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return moderator, core.ErrNoSuchModerator
		}
		logger.Log().Debug(ctx, err.Error())
		return moderator, err
	}
	logger.Log().Info(ctx, fmt.Sprintf("%+v", moderator))

	return moderator, nil
}

// UpdateModerator - метод позволяет обновить информацию о модераторе, получая его id, а так же структуру
// core.UpdateModerator, создает транзакцию, проверяет что поменялось, и обновляет, в противном случае ловит ошибку или панику, и откатывает транзакцию
func (s *store) UpdateModerator(ctx context.Context, id int, update core.UpdateModerator) (updatedModerator core.Moderator, err error) {
	tx := s.DB.WithContext(ctx).Begin()

	if tx.Error != nil {
		logger.Log().Error(ctx, tx.Error.Error())
		tx.Rollback()
		return core.Moderator{}, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Log().Error(ctx, err.Error())
			tx.Rollback()
		}
	}()

	updates := make(map[string]interface{}, 8)

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
	if update.PasswordHash != nil {
		updates["password_hash"] = *update.PasswordHash
	}

	if len(updates) == 0 {
		logger.Log().Debug(ctx, core.ErrEmptyUpdateRequest.Error())
		return core.Moderator{}, core.ErrEmptyUpdateRequest
	} else {
		updates["updated_at"] = time.Now()
	}

	err = tx.WithContext(ctx).Table(updatedModerator.TableName()).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return core.Moderator{}, core.ErrNoSuchModerator
		}
		logger.Log().Debug(ctx, err.Error())
		return core.Moderator{}, err
	}

	err = tx.WithContext(ctx).Table(updatedModerator.TableName()).Where("id = ?", id).First(&updatedModerator).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return core.Moderator{}, core.ErrNoSuchModerator
		}
		logger.Log().Debug(ctx, err.Error())
		return core.Moderator{}, err
	}

	return updatedModerator, tx.Commit().Error
}

func (s *store) AddModerator(ctx context.Context, moderator core.Moderator) (moderatorID int, err error) {
	err = s.DB.Create(&moderator).Error
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, core.ErrNotUniqueUsername
		}
		return 0, err
	}
	return moderator.ID, nil
}
