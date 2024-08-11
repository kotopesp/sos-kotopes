package userfavourite

import (
	"context"
	"errors"
	"fmt"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
)

type Store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.UserFavouriteStore {
	return &Store{pg}
}

// check for adding himself
// diff errors
// defer correct use
func (s *Store) AddUserToFavourite(ctx context.Context, favouriteUserID, userID int) (user core.User, err error) {
	tx := s.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return core.User{}, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var count int64
	err = tx.Table("favourite_persons").
		Model(&core.FavouriteUser{}).
		Where("person_id = ? AND user_id = ?", favouriteUserID, userID).
		Count(&count).Error
	if err != nil {
		return core.User{}, err
	}
	if count > 0 {
		return core.User{}, errors.New("user is already in favorites")
	}

	newFavorite := core.FavouriteUser{
		PersonID: favouriteUserID,
		UserID:   userID,
	}
	if err := tx.Table("favourite_persons").Create(&newFavorite).Error; err != nil {
		return core.User{}, err
	}

	var favouriteUser core.User
	err = tx.Table("users").Where("id = ?", favouriteUserID).First(&favouriteUser).Error

	return favouriteUser, tx.Commit().Error
}

func (s *Store) GetFavouriteUsers(ctx context.Context, userID int) (favouriteUsers []core.User, err error) {
	err = s.DB.WithContext(ctx).
		Model(&core.User{}).
		Joins("INNER JOIN favourite_persons ON users.id = favourite_persons.person_id").
		Where("favourite_persons.user_id = ?", userID).
		Find(&favouriteUsers).
		Error

	if err != nil {
		return nil, fmt.Errorf("failed to get favourite users: %w", err)
	}

	return favouriteUsers, nil
}

func (s *Store) DeleteUserFromFavourite(ctx context.Context, favouriteUserID, userID int) (err error) {
	panic("implement me")
}
