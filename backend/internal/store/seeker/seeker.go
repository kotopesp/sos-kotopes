package seeker

import (
	"context"
	"errors"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
	"gorm.io/gorm"
	"reflect"
	"time"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.SeekersStore {
	return &store{pg}
}

func (s *store) CreateSeeker(ctx context.Context, seekers core.Seekers) (core.Seekers, error) {
	// Set the creation and update timestamps
	seekers.CreatedAt = time.Now()
	seekers.UpdatedAt = time.Now()

	if err := s.DB.WithContext(ctx).Table("seekers").Create(&seekers).Error; err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.Seekers{}, err
	}

	return seekers, nil
}

func (s *store) GetSeeker(ctx context.Context, id int) (core.Seekers, error) {
	var seeker core.Seekers

	if err := s.DB.WithContext(ctx).Table("seekers").Where("id = ?", id).First(&seeker).Error; err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.Seekers{}, core.ErrAnimalNotFound
		}

		logger.Log().Error(ctx, err.Error())
		return core.Seekers{}, err
	}

	return seeker, nil
}

func (s *store) UpdateSeeker(ctx context.Context, updateSeeker core.UpdateSeekers) (core.Seekers, error) {

	var seeker core.Seekers

	tx := s.DB.WithContext(ctx).Begin()

	err := tx.Error

	if err != nil {
		logger.Log().Error(ctx, err.Error())
		tx.Rollback()
		return core.Seekers{}, err
	}

	defer func() {
		if r := recover(); r != nil || err != nil {
			logger.Log().Error(ctx, err.Error())
			tx.Rollback()
		}
	}()

	updates := make(map[string]interface{})

	v := reflect.ValueOf(updateSeeker)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := v.Type().Field(i).Name

		if !field.IsNil() {
			updates[fieldName] = field.Interface()
		}
	}

	updates["updated_at"] = time.Now()

	// Проверяем, есть ли изменения
	if len(updates) == 0 {
		logger.Log().Error(ctx, core.ErrEmptyUpdateRequest.Error())
		return core.Seekers{}, core.ErrEmptyUpdateRequest
	}

	err = tx.WithContext(ctx).Table("seekers").Where("id = ?", updateSeeker.ID).Updates(updates).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return core.Seekers{}, core.ErrNoSuchUser
		}
		logger.Log().Error(ctx, err.Error())
		return core.Seekers{}, err
	}

	err = tx.WithContext(ctx).Table("seekers").Where("id = ?", updateSeeker.ID).First(&seeker).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return core.Seekers{}, core.ErrNoSuchUser
		}
		logger.Log().Error(ctx, err.Error())
		return core.Seekers{}, err
	}

	return seeker, nil
}
