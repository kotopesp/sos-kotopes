package user_store

import (
	"context"
	"errors"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/user"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/postgres"
)

type Store struct {
	*postgres.Postgres
}

func NewUserStore(pg *postgres.Postgres) core.UserStore {
	return &Store{pg}
}

func (r *Store) UpdateUser(ctx context.Context, id int, update user.UpdateUser) error {

	tx := r.DB.WithContext(ctx).Begin()

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	updates := make(map[string]interface{})
	if update.Username != nil {
		updates["username"] = *update.Username
	}
	if update.PasswordHash != nil {
		updates["password"] = *update.PasswordHash
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
