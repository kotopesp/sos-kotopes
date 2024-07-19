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
func (s *store) Create(ctx context.Context, keeper core.Keepers) error {
	if err := s.DB.WithContext(ctx).Create(&keeper).Error; err != nil {
		return err
	}
	return nil
}

// DeleteById implements core.KeeperStore.
func (s *store) DeleteById(ctx context.Context, id int) error {
	panic("unimplemented")
}

// UpdateById implements core.KeeperStore.
func (s *store) UpdateById(ctx context.Context, id int) error {
	panic("unimplemented")
}

func (s *store) GetAll(ctx context.Context, params core.GetAllKeepersParams) ([]core.Keepers, int, error) {
	var keepers []core.Keepers
	var count int64
	panic("impl")
	return keepers, int(count), nil
}

func (s *store) GetByID(ctx context.Context, id int) (core.Keepers, error) {
	var keeper core.Keepers = core.Keepers{ID: id}

	if err := s.DB.WithContext(ctx).First(&keeper).Error; err != nil {
		return core.Keepers{}, err
	}

	return keeper, nil
}
