package role_store

import (
	"context"
	"errors"
	"fmt"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/role"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/postgres"
	"gorm.io/gorm"
	"time"
)

type Store struct {
	*postgres.Postgres
}

func NewRoleStore(pg *postgres.Postgres) core.RoleStore {
	return &Store{pg}
}

func (r *Store) GetUserRoles(ctx context.Context, id int) (roles []role.Role, err error) {

	var seeker role.Role
	if err := r.DB.WithContext(ctx).Table("seekers").
		Where("user_id = ?", id).
		Select("'seeker' as name, id, user_id, description, created_at, updated_at").
		First(&seeker).Error; err == nil {
		roles = append(roles, seeker)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error querying seeker role: %w", err)
	}

	var keeper role.Role
	if err := r.DB.WithContext(ctx).Table("keepers").
		Where("user_id = ?", id).
		Select("'keeper' as name, id, user_id, description, created_at, updated_at").
		First(&keeper).Error; err == nil {
		roles = append(roles, keeper)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error querying keeper role: %w", err)
	}

	var vet role.Role
	if err := r.DB.WithContext(ctx).Table("vets").
		Where("user_id = ?", id).
		Select("'vet' as name, id, user_id, description, created_at, updated_at").
		First(&vet).Error; err == nil {
		roles = append(roles, vet)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error querying vet role: %w", err)
	}

	return roles, nil
}

func (r *Store) GiveRoleToUser(ctx context.Context, id int, role role.GiveRole) (err error) {

	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	fmt.Println(id)
	var user core.User
	if err = tx.Table("users").First(&user, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return errors.New("user not found")
	}

	now := time.Now()

	switch role.Name {
	case "seeker":
		seeker := core.Seeker{
			UserId:      id,
			Description: role.Data,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err = tx.Table("seekers").Create(&seeker).Error; err != nil {
			tx.Rollback()
			return err
		}
	case "keeper":
		keeper := core.Role{
			UserId:      id,
			Description: role.Data,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err = tx.Table("keepers").Create(&keeper).Error; err != nil {
			tx.Rollback()
			return err
		}
	case "vet":
		vet := core.Role{
			UserId:      id,
			Description: role.Data,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err = tx.Table("vets").Create(&vet).Error; err != nil {
			tx.Rollback()
			return err
		}
	default:
		tx.Rollback()
		return errors.New("invalid role")
	}

	return tx.Commit().Error
}
func (r *Store) DeleteUserRole(ctx context.Context, id int, role string) (err error) {

	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	var user core.User
	if err = r.DB.Table("users").First(&user, "id = ?", id).Error; err != nil {
		return errors.New("user not found")
	}

	switch role {
	case "seeker":
		var seeker core.Seeker
		err = tx.Table("seekers").Where("user_id = ?", id).Delete(seeker).Error
	case "keeper":
		var keeper core.Keeper
		err = tx.Table("keeper").Where("user_id = ?", id).Delete(keeper).Error
	case "vet":
		var vet core.Vet
		err = tx.Table("vets").Where("user_id = ?", id).Delete(vet).Error
	default:
		tx.Rollback()
		return errors.New("invalid role name")
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error

}
func (r *Store) UpdateUserRole(ctx context.Context, id int, role string) (err error) {
	return nil
}
