package user_store

import (
	"context"
	"errors"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/user"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core/user_core"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/postgres"
)

type Store struct {
	*postgres.Postgres
}

func NewUserStore(pg *postgres.Postgres) user_core.UserStore {
	return &Store{pg}
}

func (r *Store) UpdateUser(ctx context.Context, id int, update user.UpdateUser) error {

	tx := r.DB.Begin()

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	updates := make(map[string]interface{})
	if update.Username != nil {
		updates["username"] = *update.Username
	}
	if update.Password != nil {
		updates["password"] = *update.Password
	}

	if len(updates) == 0 {
		return errors.New("no fields to update")
	}

	result := tx.Table("users").Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("user not found")
	}

	return tx.Commit().Error
}

func (r *Store) GetUser(ctx context.Context, id int) (user_core.User, error) {
	var user user_core.User
	err := r.DB.Table("users").Where("id = ?", id).First(&user)

	if err != nil {
		return user, err.Error
	}

	return user, nil
}
