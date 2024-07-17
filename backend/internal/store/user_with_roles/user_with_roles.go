package userwithroles

import (
	"context"

	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/postgres"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.UserWithRolesStore {
	return &store{pg}
}

func (s *store) AddUserWithRoles(ctx context.Context, data core.UserWithRoles) error {
	var err error
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	err = addUser(ctx, tx, data)
	if err != nil || data.User == nil {
		return err
	}

	err = addSeeker(ctx, tx, data)
	if err != nil {
		return err
	}

	err = addKeeper(ctx, tx, data)
	if err != nil {
		return err
	}

	err = addVet(ctx, tx, data)
	if err != nil {
		return err
	}

	return tx.Commit().Error
}
