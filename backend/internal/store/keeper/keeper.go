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

func (s *store) Create(ctx context.Context, keeper core.Keepers) error {
	if err := s.DB.WithContext(ctx).Create(&keeper).Error; err != nil {
		return err
	}
	return nil
}

func (s *store) SoftDeleteByID(ctx context.Context, id int) error {
	err := s.DB.WithContext(ctx).Model(&core.Keepers{}).Where("id = ?", id).Updates(core.Keepers{IsDeleted: true, DeletedAt: time.Now()}).Error

	return err
}

func (s *store) UpdateByID(ctx context.Context, keeper core.UpdateKeepers) (core.Keepers, error) {
	keeper.UpdatedAt = time.Now()

	var updatedKeeper core.Keepers

	result := s.DB.WithContext(ctx).Model(&core.Keepers{}).Where("id = ? AND is_deleted = ?", keeper.ID, false).Updates(keeper).First(&updatedKeeper, keeper.ID)
	if result.Error != nil {

		if errors.Is(result.Error, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.Keepers{}, core.ErrRecordNotFound
		}

		logger.Log().Error(ctx, result.Error.Error())
		return core.Keepers{}, result.Error
	}

	return updatedKeeper, nil
}

func (s *store) GetAll(ctx context.Context, params core.GetAllKeepersParams) ([]core.Keepers, error) {
	var keepers []core.Keepers
	query := s.DB.WithContext(ctx).Model(&core.Keepers{})

	if params.MinPrice != nil {
		query = query.Where("price >= ?", *params.MinPrice)
	}
	if params.MaxPrice != nil {
		query = query.Where("price <= ?", *params.MaxPrice)
	}

	query = query.Select("keepers.*, AVG(keeper_reviews.grade) as avg_grade").
		Joins("left join keeper_reviews on keeper_reviews.keeper_id = keepers.id").
		Group("keepers.id")

	if params.MinRating != nil {
		query = query.Having("AVG(keeper_reviews.grade) >= ?", *params.MinRating)
	}
	if params.MaxRating != nil {
		query = query.Having("AVG(keeper_reviews.grade) <= ?", *params.MaxRating)
	}

	query = query.Order(*params.SortBy + " " + *params.SortOrder)

	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}
	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	if err := query.Find(&keepers).Error; err != nil {
		logger.Log().Debug(ctx, err.Error())
		return nil, err
	}

	return keepers, nil
}

func (s *store) GetByID(ctx context.Context, id int) (core.Keepers, error) {
	var keeper = core.Keepers{ID: id}

	if err := s.DB.WithContext(ctx).First(&keeper).Error; err != nil {
		return core.Keepers{}, err
	}

	return keeper, nil
}
