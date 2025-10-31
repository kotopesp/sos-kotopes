package seeker

import (
	"context"
	"errors"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.SeekersStore {
	return &store{pg}
}

func (s *store) SeekersDB(ctx context.Context) *gorm.DB {
	return s.DB.WithContext(ctx).Table("seekers")
}

func (s *store) CreateSeeker(ctx context.Context, seeker core.Seeker) (createSeeker core.Seeker, err error) {
	tx := s.SeekersDB(ctx).Begin()

	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	seeker.CreatedAt = time.Now()
	seeker.UpdatedAt = time.Now()

	if err = tx.Create(&seeker).Error; err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.Seeker{}, err
	}

	if err = tx.
		Where("user_id = ?", seeker.UserID).
		Preload("User").
		First(&seeker).Error; err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return core.Seeker{}, core.ErrNoSuchUser
		}

		logger.Log().Error(ctx, err.Error())
		return core.Seeker{}, err
	}

	return seeker, tx.Commit().Error
}

func (s *store) GetSeeker(ctx context.Context, id int) (core.Seeker, error) {
	var seeker core.Seeker

	err := s.SeekersDB(ctx).
		Where("user_id = ?", id).
		Preload("User").
		First(&seeker).Error
	if err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return core.Seeker{}, core.ErrSeekerNotFound
		}

		logger.Log().Error(ctx, err.Error())
		return core.Seeker{}, err
	}

	return seeker, nil
}

func (s *store) UpdateSeeker(ctx context.Context, id int, updateSeeker map[string]interface{}) (seeker core.Seeker, err error) {
	tx := s.SeekersDB(ctx).Begin()

	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	if err = tx.
		Where("user_id = ?", id).
		Preload("User").
		First(&seeker).Error; err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return core.Seeker{}, core.ErrNoSuchUser
		}

		logger.Log().Error(ctx, err.Error())
		return core.Seeker{}, err
	}

	if err = tx.
		Model(&seeker).
		Where("user_id = ?", id).
		Clauses(clause.Returning{}).
		Updates(updateSeeker).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return core.Seeker{}, core.ErrSeekerNotFound
		}
		logger.Log().Error(ctx, err.Error())
		return core.Seeker{}, err
	}

	return seeker, tx.Commit().Error
}

func (s *store) DeleteSeeker(ctx context.Context, userID int) error {
	updates := make(map[string]interface{})
	updates["is_deleted"] = true

	if err := s.SeekersDB(ctx).
		Where("user_id = ?", userID).
		Updates(updates).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return core.ErrSeekerNotFound
		}
		logger.Log().Error(ctx, err.Error())
		return err
	}

	return nil
}

func createQuery(query *gorm.DB, params core.GetAllSeekersParams) *gorm.DB {
	query.Where("is_deleted = false")

	if params.AnimalType != nil {
		query = query.Where("animal_type = ?", *params.AnimalType)
	}

	if params.LocationID != nil {
		query = query.Where("location_id = ?", *params.LocationID)
	}

	if params.MinPrice != nil {
		query = query.Where("price > ?", *params.MinPrice)
	}

	if params.MaxPrice != nil {
		query = query.Where("price < ?", *params.MaxPrice)
	}

	if params.MinEquipmentRental != nil {
		query = query.Where("equipment_rental > ?", params.MinEquipmentRental)
	}

	if params.MaxEquipmentRental != nil {
		query = query.Where("equipment_rental < ?", params.MaxEquipmentRental)
	}

	if params.HaveMetalCage != nil {
		query = query.Where("have_metal_cage = ?", *params.HaveMetalCage)
	}

	if params.HavePlasticCage != nil {
		query = query.Where("have_plastic_cage = ?", *params.HavePlasticCage)
	}

	if params.HaveNet != nil {
		query = query.Where("have_net = ?", *params.HaveNet)
	}

	if params.HaveLadder != nil {
		query = query.Where("have_ladder = ?", *params.HaveLadder)
	}

	if params.HaveOther != nil {
		query = query.Where("have_other != ''")
	}

	if params.HaveCar != nil {
		query = query.Where("have_car = ?", *params.HaveCar)
	}

	query = query.Order(*params.SortBy + " " + *params.SortOrder)
	query = query.Offset(*params.Offset)
	query = query.Limit(*params.Limit)

	return query
}

func (s *store) GetAllSeekers(ctx context.Context, params core.GetAllSeekersParams) ([]core.Seeker, error) {
	var seekers []core.Seeker
	query := s.SeekersDB(ctx).Model(&core.Seeker{})
	query = createQuery(query, params)
	err := query.Preload("User",
		func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", core.UserActive)
		}).Find(&seekers).Error
	if err != nil {
		logger.Log().Debug(ctx, err.Error())
		return nil, err
	}

	return seekers, nil
}
