package name

import (
	"context"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

type service struct {
	entityStore core.EntityStore
}

func New(store core.EntityStore) core.EntityService {
	return &service{
		entityStore: store,
	}
}

func (s *service) GetAll(_ context.Context, _ core.GetAllParams) (data []core.Entity, total int, err error) {
	panic("implement me")
}

func (s *service) GetByID(_ context.Context, _ int) (data core.Entity, err error) {
	panic("implement me")
}
