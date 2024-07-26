package entity

import (
	"context"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.EntityStore {
	return &store{pg}
}

func (s *store) GetAll(_ context.Context, _ core.GetAllParams) (data []core.Entity, err error) {
	panic("implement me")
}

func (s *store) GetByID(_ context.Context, _ int) (data core.Entity, err error) {
	panic("implement me")
}
