package role

import (
	"context"
	"errors"
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

func (s *store) GetUserRoles(ctx context.Context, id int) (roles map[string]core.Role, err error) {
	roles = make(map[string]core.Role)
	tx := s.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	tableNames := []string{"seekers", "keepers", "vets"}
	roleNames := []string{Seeker, Keeper, Vet}

	for i, name := range tableNames {
		var role core.Role
		key := roleNames[i]
		if err = tx.Table(name).
			Where("user_id = ?", id).
			First(&role).Error; err == nil {
			roles[key] = role
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	return roles, tx.Commit().Error
}

func (s *store) GiveRoleToUser(ctx context.Context, id int, givenRole core.GivenRole) (addedRole core.Role, err error) {

	tx := s.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return core.Role{}, tx.Error
	}

	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	var user core.User
	if err = tx.Table("users").First(&user, "id = ?", id).Error; err != nil {
		return core.Role{}, core.ErrNoSuchUser
	}

	var tableName string
	switch givenRole.Name {
	case Seeker:
		tableName = "seekers"
	case Keeper:
		tableName = "keepers"
	case Vet:
		tableName = "vets"
	default:
		return core.Role{}, core.ErrInvalidRole
	}
	now := time.Now()
	role := core.Role{
		UserID:      id,
		Description: givenRole.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := tx.Table(tableName).Create(&role).Error; err != nil {
		return core.Role{}, err
	} else {
		addedRole = role
	}

	return addedRole, tx.Commit().Error
}
func (s *store) DeleteUserRole(ctx context.Context, id int, roleName string) (err error) {

	tx := s.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	var role core.Role
	switch roleName {
	case Seeker:
		err = tx.Table("seekers").Where("user_id = ?", id).Delete(role).Error
	case Keeper:
		err = tx.Table("keepers").Where("user_id = ?", id).Delete(role).Error
	case Vet:
		err = tx.Table("vets").Where("user_id = ?", id).Delete(role).Error
	default:
		tx.Rollback()
		return core.ErrInvalidRole
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error

}

func (s *store) UpdateUserRole(ctx context.Context, id int, updateRole core.UpdateRole) (updatedRole core.Role, err error) {
	tx := s.DB.WithContext(ctx).Begin()

	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		tx.Rollback()
		return core.Role{}, err
	}

	updates := make(map[string]interface{})
	if updateRole.Description != nil {
		updates["description"] = *updateRole.Description
	}

	if len(updates) == 0 {
		return core.Role{}, core.ErrNoFieldsToUpdate
	} else {
		updates["updated_at"] = time.Now()
	}

	roleName := updateRole.Name
	var tableName string
	switch roleName {
	case Seeker:
		tableName = "seekers"
	case Keeper:

		tableName = "keepers"
	case Vet:
		tableName = "vets"
	default:
		return core.Role{}, core.ErrInvalidRole
	}
	err = tx.Table(tableName).Where("user_id = ?", id).Updates(updates).Error
	if err != nil {
		return core.Role{}, err
	}

	err = tx.Table(tableName).Where("user_id = ?", id).First(&updatedRole).Error
	if err != nil {
		return core.Role{}, err
	}
	return updatedRole, tx.Commit().Error
}
