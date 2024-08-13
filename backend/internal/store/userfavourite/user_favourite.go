package userfavourite

import (
	"context"
	"errors"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
	"gorm.io/gorm"
)

type Store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.UserFavouriteStore {
	return &Store{pg}
}

func (s *Store) AddUserToFavourite(ctx context.Context, favouriteUserID, userID int) (favouriteUser core.User, err error) {
	tx := s.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		logger.Log().Error(ctx, tx.Error.Error())
		return core.User{}, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err = tx.WithContext(ctx).Table("users").Where("id = ?", favouriteUserID).First(&favouriteUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return core.User{}, core.ErrNoSuchUser
		}
		logger.Log().Error(ctx, err.Error())
		return core.User{}, err
	}

	err = tx.WithContext(ctx).Table("users").Where("id = ?", userID).First(&core.User{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return core.User{}, core.ErrNoSuchUser
		}
		logger.Log().Error(ctx, err.Error())
		return core.User{}, err
	}

	if favouriteUserID == userID {
		logger.Log().Debug(ctx, core.ErrCantAddYourselfIntoFavourites.Error())
		return core.User{}, core.ErrCantAddYourselfIntoFavourites
	}

	var count int64
	err = tx.Table("favourite_persons").
		Model(&core.FavouriteUser{}).
		Where("person_id = ? AND user_id = ?", favouriteUserID, userID).
		Count(&count).Error
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.User{}, err
	}

	if count > 0 {
		logger.Log().Debug(ctx, core.ErrUserAlreadyInFavourites.Error())
		return core.User{}, core.ErrUserAlreadyInFavourites
	}

	newFavorite := core.FavouriteUser{
		PersonID: favouriteUserID,
		UserID:   userID,
	}
	if err = tx.Table("favourite_persons").Create(&newFavorite).Error; err != nil {
		return core.User{}, err
	}

	return favouriteUser, tx.Commit().Error
}

func (s *Store) GetFavouriteUsers(ctx context.Context, userID int) (favouriteUsers []core.User, err error) {
	// is it works correcly?
	err = s.DB.WithContext(ctx).
		Model(&core.User{}).
		Joins("INNER JOIN favourite_persons ON users.id = favourite_persons.person_id").
		Where("favourite_persons.user_id = ?", userID).
		Find(&favouriteUsers).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return nil, core.ErrNoSuchUser
		}
		logger.Log().Error(ctx, err.Error())
		return nil, err
	}

	return favouriteUsers, nil
}

func (s *Store) DeleteUserFromFavourite(ctx context.Context, favouriteUserID, userID int) (err error) {
	err = s.DB.WithContext(ctx).Table("favourite_persons").
		Where("person_id = ? AND user_id = ?", favouriteUserID, userID).
		Delete(core.FavouriteUser{}).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return core.ErrNoSuchUser
		}
		logger.Log().Error(ctx, err.Error())
		return err
	}

	return nil
}
