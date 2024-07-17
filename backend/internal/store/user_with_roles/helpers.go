package userwithroles

import (
	"context"
	"strings"

	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gorm.io/gorm"
)

func addUser(ctx context.Context, tx *gorm.DB, data core.UserWithRoles) error {
	user := data.User
	if user != nil {
		err := tx.WithContext(ctx).Create(user).Error
		if err != nil {
			if strings.Contains(err.Error(), "users_username_key") {
				return ErrNotUniqueUsername
			}
			return err
		}
	}
	return nil
}

func addSeeker(ctx context.Context, tx *gorm.DB, data core.UserWithRoles) error {
	seeker := data.Seeker
	if seeker != nil {
		seeker.UserID = data.User.ID
		err := tx.WithContext(ctx).Create(seeker).Error
		if err != nil {
			return err
		}
		err = tx.WithContext(ctx).Create(&core.RoleUser{
			RoleID: 1,
			UserID: seeker.UserID,
		}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func addKeeper(ctx context.Context, tx *gorm.DB, data core.UserWithRoles) error {
	keeper := data.Keeper
	if keeper != nil {
		keeper.UserID = data.User.ID
		err := tx.WithContext(ctx).Create(keeper).Error
		if err != nil {
			return err
		}
		err = tx.WithContext(ctx).Create(&core.RoleUser{
			RoleID: 2,
			UserID: keeper.UserID,
		}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func addVet(ctx context.Context, tx *gorm.DB, data core.UserWithRoles) error {
	vet := data.Vet
	if vet != nil {
		vet.UserID = data.User.ID
		err := tx.WithContext(ctx).Create(vet).Error
		if err != nil {
			return err
		}
		err = tx.WithContext(ctx).Create(&core.RoleUser{
			RoleID: 3,
			UserID: vet.UserID,
		}).Error
		if err != nil {
			return err
		}
	}
	return nil
}
