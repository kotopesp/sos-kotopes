package role

import (
	"context"
	"errors"
	"fmt"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/role"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
	"gorm.io/gorm"
	"time"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.RoleStore {
	return &store{pg}
}

const Seeker = "seeker"
const Keeper = "keeper"
const Vet = "vet"

func (s *store) GetUserRoles(ctx context.Context, id int) (roles []core.Role, err error) {

	var seeker core.Role
	if err := s.DB.WithContext(ctx).Table("seekers").
		Where("user_id = ?", id).
		Select("'seeker' as name, id, user_id, description, created_at, updated_at").
		First(&seeker).Error; err == nil {
		roles = append(roles, seeker)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error querying seeker role: %w", err)
	}

	var keeper core.Role
	if err := s.DB.WithContext(ctx).Table("keepers").
		Where("user_id = ?", id).
		Select("'keeper' as name, id, user_id, description, created_at, updated_at").
		First(&keeper).Error; err == nil {
		roles = append(roles, keeper)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error querying keeper role: %w", err)
	}

	var vet core.Role
	if err := s.DB.WithContext(ctx).Table("vets").
		Where("user_id = ?", id).
		Select("'vet' as name, id, user_id, description, created_at, updated_at").
		First(&vet).Error; err == nil {
		roles = append(roles, vet)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error querying vet role: %w", err)
	}

	return roles, nil
}

func (s *store) GiveRoleToUser(ctx context.Context, id int, givenRole role.GiveRole) (err error) {

	tx := s.DB.Begin()
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

	switch givenRole.Name {
	case Seeker:
		seeker := core.Seeker{
			UserID:      id,
			Description: givenRole.Description,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err = tx.Table("seekers").Create(&seeker).Error; err != nil {
			tx.Rollback()
			return err
		}
	case Keeper:
		keeper := core.Role{
			UserID:      id,
			Description: givenRole.Description,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err = tx.Table("keepers").Create(&keeper).Error; err != nil {
			tx.Rollback()
			return err
		}
	case Vet:
		vet := core.Role{
			UserID:      id,
			Description: givenRole.Description,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err = tx.Table("vets").Create(&vet).Error; err != nil {
			tx.Rollback()
			return err
		}
	default:
		tx.Rollback()
		return errors.New("invalid givenRole")
	}

	return tx.Commit().Error
}
func (s *store) DeleteUserRole(ctx context.Context, id int, roleName string) (err error) {

	tx := s.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	var user core.User
	if err = s.DB.Table("users").First(&user, "id = ?", id).Error; err != nil {
		return errors.New("user not found")
	}

	switch roleName {
	case Seeker:
		var seeker core.Seeker
		err = tx.Table("seekers").Where("user_id = ?", id).Delete(seeker).Error
	case Keeper:
		var keeper core.Keeper
		err = tx.Table("keepers").Where("user_id = ?", id).Delete(keeper).Error
	case Vet:
		var vet core.Vet
		err = tx.Table("vets").Where("user_id = ?", id).Delete(vet).Error
	default:
		tx.Rollback()
		return errors.New("invalid roleName name")
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error

}

// check for user existing
func (s *store) UpdateUserRole(ctx context.Context, id int, updateRole role.UpdateRole) (err error) {
	tx := s.DB.Begin()

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	updates := make(map[string]interface{})
	if updateRole.Description != nil {
		updates["description"] = *updateRole.Description
	}

	if len(updates) == 0 {
		return errors.New("no fields to update")
	} else {
		updates["updated_at"] = time.Now()
	}
	roleName := updateRole.Name
	switch roleName {
	case Seeker:
		result := tx.Table("seekers").Where("id = ?", id).Updates(updates)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	case Keeper:
		result := tx.Table("keepers").Where("id = ?", id).Updates(updates)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	case Vet:
		result := tx.Table("vets").Where("id = ?", id).Updates(updates)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	default:
		tx.Rollback()
		return errors.New("invalid updateRole name")
	}

	return tx.Commit().Error
}
