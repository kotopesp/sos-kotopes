package keeperService

import (
	"context"

	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
)

type service struct {
	keeperStore core.KeeperStore
}

func New(keeperStore core.KeeperStore) core.KeeperService {
	return &service{keeperStore: keeperStore}
}

func (s *service) GetAll(ctx context.Context, params core.GetAllKeepersParams) ([]core.Keepers, int, error) {
	return s.keeperStore.GetAll(ctx, params)
}

func (s *service) GetByID(ctx context.Context, id int) (core.Keepers, error) {
	return s.keeperStore.GetByID(ctx, id)
}

func (s *service) Create(ctx context.Context, keeper core.Keepers) error {
	return s.keeperStore.Create(ctx, keeper)
}

func (s *service) DeleteById(ctx context.Context, id int) error {
	return s.keeperStore.DeleteById(ctx, id)
}

func (s *service) UpdateById(ctx context.Context, id int) error {
	return s.keeperStore.UpdateById(ctx, id)
}
