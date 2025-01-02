package keepestore

import (
	"context"
	"errors"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.KeeperStore {
	return &store{pg}
}

func (s *store) CreateKeeper(ctx context.Context, keeper core.Keeper) (data core.Keeper, err error) {
	keeperExist := false
	if err := s.DB.WithContext(ctx).Table(keeper.TableName()).Select("1").Where("user_id = ? AND is_deleted = ?", keeper.UserID, false).Limit(1).Find(&keeperExist).Error; err != nil {
		return core.Keeper{}, err
	}
	if keeperExist {
		return core.Keeper{}, core.ErrKeeperUserAlreadyKeeper
	}
	if err := s.DB.WithContext(ctx).Create(&keeper).Error; err != nil {
		return core.Keeper{}, err
	}

	if err := s.DB.WithContext(ctx).Preload("User").First(&keeper).Error; err != nil {
		return core.Keeper{}, err
	}

	return keeper, nil
}

func (s *store) DeleteKeeper(ctx context.Context, id int) error {
	now := time.Now()
	err := s.DB.WithContext(ctx).Model(&core.Keeper{}).Where("id = ?", id).Updates(core.Keeper{IsDeleted: true, DeletedAt: &now}).Error

	return err
}

func (s *store) UpdateKeeper(ctx context.Context, id int, keeper core.Keeper) (core.Keeper, error) {
	var updatedKeeper core.Keeper
	updatedKeeper.UpdatedAt = time.Now()

	if err := s.DB.WithContext(ctx).Model(&core.Keeper{}).Where("id = ? AND is_deleted = ?", id, false).Updates(keeper).Error; err != nil {
		logger.Log().Debug(ctx, err.Error())
		return core.Keeper{}, err
	}

	if err := s.DB.WithContext(ctx).Preload("User").First(&updatedKeeper, "id = ? AND is_deleted = ?", id, false).Error; err != nil {
		switch {
		case errors.Is(err, core.ErrRecordNotFound):
			logger.Log().Debug(ctx, core.ErrRecordNotFound.Error())
			return core.Keeper{}, core.ErrRecordNotFound
		default:
			logger.Log().Debug(ctx, err.Error())
			return core.Keeper{}, err
		}
	}

	return updatedKeeper, nil
}

func (s *store) GetAllKeepers(ctx context.Context, params core.GetAllKeepersParams) (data []core.Keeper, err error) {
	var keepersExist core.Keeper
	if err := s.DB.WithContext(ctx).Model(&core.Keeper{}).Where("keepers.is_deleted = ?", false).First(&keepersExist).Error; err != nil {
		logger.Log().Debug(ctx, err.Error())
		print(err.Error())
		return []core.Keeper{}, nil
	}

	var keepers []core.Keeper
	query := s.DB.WithContext(ctx).Model(&core.Keeper{}).Where("keepers.is_deleted = ?", false)

	if params.MinPrice != nil {
		query = query.Where("price >= ?", *params.MinPrice)
	}
	if params.MaxPrice != nil {
		query = query.Where("price <= ?", *params.MaxPrice)
	}

	query = query.Select("keepers.*, AVG(keeper_reviews.grade) as avg_grade").
		Joins("LEFT JOIN keeper_reviews ON keeper_reviews.keeper_id = keepers.id").
		Group("keepers.id")

	if params.MinRating != nil {
		query = query.Having("AVG(keeper_reviews.grade) >= ?", *params.MinRating)
	}
	if params.MaxRating != nil {
		query = query.Having("AVG(keeper_reviews.grade) <= ?", *params.MaxRating)
	}

	if params.HasCage != nil {
		query = query.Where("has_cage = ?", *params.HasCage)
	}
	if params.BoardingDuration != nil {
		query = query.Where("boarding_duration = ?", *params.BoardingDuration)
	}
	if params.BoardingCompensation != nil {
		query = query.Where("boarding_compensation = ?", *params.BoardingCompensation)
	}
	if params.AnimalAcceptance != nil {
		query = query.Where("animal_acceptance = ?", *params.AnimalAcceptance)
	}
	if params.AnimalCategory != nil {
		query = query.Where("animal_category = ?", *params.AnimalCategory)
	}
	if params.LocationID != nil {
		query = query.Where("location_id = ?", *params.LocationID)
	}

	query = query.Order(*params.SortBy + " " + *params.SortOrder)

	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}
	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	if err := query.Preload("User").Find(&keepers).Error; err != nil {
		logger.Log().Debug(ctx, err.Error())
		return nil, err
	}

	return keepers, nil
}

func (s *store) GetKeeperByID(ctx context.Context, id int) (core.Keeper, error) {
	var keeper = core.Keeper{ID: id}

	if err := s.DB.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).Preload("User").First(&keeper).Error; err != nil {
		return core.Keeper{}, err
	}

	return keeper, nil
}
