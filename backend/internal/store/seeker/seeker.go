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

func (s *store) CreateSeeker(ctx context.Context, seeker core.Seeker, equipment core.Equipment) (core.Seeker, error) {
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.Seeker{}, err
	}

	err := tx.WithContext(ctx).Table("equipments").Create(&equipment).Error // CreateEquipment теперь работает в транзакции
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.Seeker{}, err
	}

	seeker.CreatedAt = time.Now()
	seeker.UpdatedAt = time.Now()
	seeker.EquipmentID = equipment.ID

	err = tx.WithContext(ctx).Table("seekers").Create(&seeker).Error // CreateSeeker теперь работает в транзакции
	if err != nil {
		tx.Rollback()
		logger.Log().Error(ctx, err.Error())
		return core.Seeker{}, err
	}

	return seeker, nil
}

func (s *store) GetSeeker(ctx context.Context, id int) (core.Seeker, error) {
	var seeker core.Seeker

	if err := s.DB.WithContext(ctx).Table("seekers").Where("user_id = ?", id).First(&seeker).Error; err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.Seeker{}, core.ErrSeekerNotFound
		}

		logger.Log().Error(ctx, err.Error())
		return core.Seeker{}, err
	}

	return seeker, nil
}

func (s *store) UpdateSeeker(ctx context.Context, updateSeeker core.UpdateSeeker) (core.Seeker, error) {

	var seeker core.Seeker

	tx := s.DB.WithContext(ctx).Begin()

	err := tx.Error

	if err != nil {
		logger.Log().Error(ctx, err.Error())
		tx.Rollback()
		return core.Seeker{}, err
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
		return core.Seeker{}, core.ErrEmptyUpdateRequest
	}

	err = tx.WithContext(ctx).Table("seekers").Where("id = ?", updateSeeker.ID).Updates(updates).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return core.Seeker{}, core.ErrNoSuchUser
		}
		logger.Log().Error(ctx, err.Error())
		return core.Seeker{}, err
	}

	err = tx.WithContext(ctx).Table("seekers").Where("id = ?", updateSeeker.ID).First(&seeker).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return core.Seeker{}, core.ErrNoSuchUser
		}
		logger.Log().Error(ctx, err.Error())
		return core.Seeker{}, err
	}

	return seeker, nil
}
