package user

import (
	"context"

	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/postgres"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.UserStore {
	return &store{pg}
}

func (s *store) GetUserByUsername(ctx context.Context, username string) (data core.User, err error) {
	user := core.User{}
	result := s.DB.WithContext(ctx).First(&user, "username=?", username)
	data = user
	err = result.Error
	return
}

func (s *store) GetUserByID(ctx context.Context, id int) (data core.User, err error) {
	user := core.User{}
	result := s.DB.WithContext(ctx).First(&user, id)
	data = user
	err = result.Error
	return
}
