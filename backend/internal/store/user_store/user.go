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
		return errors.New("user_service not found")
	}

	return tx.Commit().Error
}

func (r *Store) GetUser(ctx context.Context, id int) (user core.User, err error) {
	err = r.DB.Table("users").Where("id = ?", id).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *Store) GetUserPosts(ctx context.Context, id int) (posts []core.Post, err error) {
	err = r.DB.
		Where("user_id = ?", id).
		Order("created_at DESC").
		Find(&posts).Error

	if err != nil {
		return nil, err
	}

	return posts, nil
}
