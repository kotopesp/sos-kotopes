package keepersService

import (
	"context"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
)

type keeperService struct {
	keepersStore core.KeepersStore
}

func New(store core.KeepersStore) core.keeperService {
	return &service{
		keepersStore: store,
	}
}

func (s *keeperService) GetAll(ctx context.Context, params core.GetAllParams) (data []core.Keeper, total int, err error) {
	panic("implement me")
}

func (s *keeperService) GetByID(ctx context.Context, ID int) (keeper core.Keeper, err error) {
	panic("implement me")
}

func (s *keeperService) Create(ctx context.Context, ID int) (keeper core.Keeper, err error) {
	panic("implement me")
}

func (s *keeperService) Update(ctx context.Context, ID int) (err error) {
	panic("implement me")
}

func (s *keeperService) Delete(ctx context.Context, ID int) (err error) {
	panic("implement me")
}
