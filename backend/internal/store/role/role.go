package role

import (
	"context"

	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/postgres"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.RoleStore {
	return &store{pg}
}

func (s *store) GetRoleByName(ctx context.Context, name string) (data core.Role, err error) {
	role := core.Role{Name: name}
	result := s.DB.WithContext(ctx).First(&role)
	data = role
	err = result.Error
	return
}
