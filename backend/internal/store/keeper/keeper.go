package keeperStore

import (
	"context"
	"fmt"

	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/postgres"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.KeeperStore {
	return &store{pg}
}

// Create implements core.KeeperStore.
func (s *store) Create(ctx *context.Context, keeper core.Keeper) error {
	if err := s.DB.WithContext(*ctx).Create(&keeper).Error; err != nil {
		return err
	}
	return nil
}

// DeleteById implements core.KeeperStore.
func (s *store) DeleteById(ctx *context.Context, id int) error {
	panic("unimplemented")
}

// UpdateById implements core.KeeperStore.
func (s *store) UpdateById(ctx *context.Context, id int) error {
	panic("unimplemented")
}

func (s *store) GetAll(ctx *context.Context, params core.GetAllKeepersParams) ([]core.Keeper, error) {
	var keepers []core.Keeper
	query := s.DB.WithContext(*ctx).Model(&core.Keeper{})

	if params.SortBy != nil && params.SortOrder != nil {
		query = query.Order(*params.SortBy + " " + *params.SortOrder)
	}
	if params.Location != nil {
		query = query.Where("location = ?", *params.Location)
	}
	if params.MinRating != nil {
		query = query.Where("rating >= ?", *params.MinRating)
	}
	if params.MaxRating != nil {
		query = query.Where("rating <= ?", *params.MaxRating)
	}
	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}
	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	if err := query.Find(&keepers).Error; err != nil {
		return nil, fmt.Errorf("keepers.GetAll.find() error")
	}

	return keepers, nil
}

func (s *store) GetByID(ctx *context.Context, id int) (core.Keeper, error) {
	var keeper core.Keeper = core.Keeper{ID: id}

	if err := s.DB.WithContext(*ctx).First(&keeper).Error; err != nil {
		return core.Keeper{}, err
	}

	return keeper, nil
}
