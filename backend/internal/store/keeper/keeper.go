package keeperStore

import (
	"context"

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
	panic("impl")
	return keepers, nil
}

func (s *store) GetByID(ctx *context.Context, id int) (core.Keeper, error) {
	var keeper core.Keeper = core.Keeper{ID: id}

	if err := s.DB.WithContext(*ctx).First(&keeper).Error; err != nil {
		return core.Keeper{}, err
	}

	return keeper, nil
}
