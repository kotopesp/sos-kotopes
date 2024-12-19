package seeker

import (
	"context"
	"errors"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
	"gorm.io/gorm"
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

	err := tx.WithContext(ctx).Table("equipments").Create(&equipment).Error
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.Seeker{}, err
	}

	seeker.CreatedAt = time.Now()
	seeker.UpdatedAt = time.Now()
	seeker.EquipmentID = equipment.ID

	err = tx.WithContext(ctx).Table("seekers").Create(&seeker).Error
	if err != nil {
		tx.Rollback()
		logger.Log().Error(ctx, err.Error())
		return core.Seeker{}, err
	}

	if err := tx.Commit().Error; err != nil {
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

func (s *store) UpdateSeeker(ctx context.Context, id int, updateSeeker map[string]interface{}) (core.Seeker, error) {
	var seeker core.Seeker

	updateSeeker["updated_at"] = time.Now()

	err := s.DB.WithContext(ctx).Table("seekers").Where("id = ?", id).Updates(updateSeeker).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return core.Seeker{}, core.ErrNoSuchUser
		}
		logger.Log().Error(ctx, err.Error())
		return core.Seeker{}, err
	}

	err = s.DB.WithContext(ctx).Table("seekers").Where("id = ?", id).First(&seeker).Error
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

func (s *store) DeleteSeeker(ctx context.Context, userID int) error {
	updates := make(map[string]interface{})
	updates["is_deleted"] = true

	err := s.DB.WithContext(ctx).Table("seekers").Where("user_id = ?", userID).Updates(updates).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return core.ErrNoSuchUser
		}
		logger.Log().Error(ctx, err.Error())
		return err
	}

	return nil
}

func (s *store) GetAllSeekers(ctx context.Context, params core.GetAllSeekersParams) ([]core.Seeker, error) {
	var seekers []core.Seeker
	query := s.DB.WithContext(ctx).Model(&core.Seeker{})
	query.Where("is_deleted = false")

	if params.AnimalType != nil {
		query = query.Where("animal_type = ?", *params.AnimalType)
	}

	if params.Location != nil {
		query = query.Where("location = ?", *params.Location)
	}

	if params.Price != nil {
		if *params.Price < 0 {
			query = query.Where("price < 0")
		}

		if *params.Price == 0 {
			query = query.Where("price == 0")
		}

		if *params.Price > 0 {
			query = query.Where("price > 0")
		}
	}

	if params.HaveCar != nil {
		query = query.Where("have_car = ?", *params.HaveCar)
	}

	query = query.Order(*params.SortBy + " " + *params.SortOrder)

	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}

	if err := query.Preload("User").Find(&seekers).Error; err != nil {
		logger.Log().Debug(ctx, err.Error())
		return nil, err
	}

	return seekers, nil
}
