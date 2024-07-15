package user

import (
	"context"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/user"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/postgres"
)

type UserStore struct {
	*postgres.Postgres
}

func NewUserStore(pg *postgres.Postgres) core.UserStore {
	return &UserStore{pg}
}

func (r *UserStore) ChangeName(ctx context.Context, id int, name string) (err error) {
	userModel := user.User{}
	err := r.DB.Save(id)

	return nil
}
