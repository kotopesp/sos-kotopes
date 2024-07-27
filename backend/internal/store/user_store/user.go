package user_store

import (
	"context"
	"errors"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/user"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/postgres"
	"time"
)

type Store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.UserStore {
	return &Store{pg}
}

// проверка на существование пользователя?
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
	if update.FirstName != nil {
		updates["firstname"] = *update.FirstName
	}
	if update.LastName != nil {
		updates["lastname"] = *update.LastName
	}
	if update.Description != nil {
		updates["description"] = *update.Description
	}
	if update.Photo != nil {
		updates["photo"] = *update.Photo
	}
	if update.PasswordHash != nil {
		updates["password_hash"] = *update.PasswordHash
	}
	if update.IsDeleted != nil {
		updates["is_deleted"] = *update.IsDeleted
	}
	if update.DeletedAt != nil {
		updates["deleted_at"] = *update.DeletedAt
	}
	if update.CreatedAt != nil {
		updates["created_at"] = *update.CreatedAt
	}

	if len(updates) == 0 {
		return errors.New("no fields to update")
	} else {
		updates["updated_at"] = time.Now()
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
