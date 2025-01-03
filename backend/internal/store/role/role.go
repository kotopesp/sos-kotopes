package role

import (
	"context"
	"errors"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
	"gorm.io/gorm"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.RoleStore {
	return &store{pg}
}

func (s *store) GiveRoleToUser(ctx context.Context, id int, givenRole core.GivenRole) (addedRole core.Role, err error) {
	tx := s.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		logger.Log().Error(ctx, err.Error())
		tx.Rollback()
		return core.Role{}, tx.Error
	}

	defer func() {
		if r := recover(); r != nil || err != nil {
			logger.Log().Error(ctx, err.Error())
			tx.Rollback()
		}
	}()

	var user core.User
	if err = tx.Table("users").First(&user, "id = ?", id).Error; err != nil {
		return core.Role{}, core.ErrNoSuchUser
	}

	var tableName string
	switch givenRole.Name {
	case core.Seeker:
		tableName = "seekers"
	case core.KeeperRole:
		tableName = "keepers"
	case core.Vet:
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
	}
	addedRole = role

	return addedRole, tx.Commit().Error
}

func (s *store) GetUserRoles(ctx context.Context, id int) (roles map[string]core.Role, err error) {
	roles = make(map[string]core.Role)
	tx := s.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		logger.Log().Error(ctx, err.Error())
		tx.Rollback()
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil || err != nil {
			logger.Log().Error(ctx, err.Error())
			tx.Rollback()
		}
	}()

	tableNames := []string{"seekers", "keepers", "vets"}
	roleNames := []string{core.Seeker, core.KeeperRole, core.Vet}

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

func (s *store) UpdateUserRole(ctx context.Context, id int, updateRole core.UpdateRole) (updatedRole core.Role, err error) {
	tx := s.DB.WithContext(ctx).Begin()

	if tx.Error != nil {
		tx.Rollback()
		logger.Log().Error(ctx, err.Error())
		return core.Role{}, err
	}

	defer func() {
		if r := recover(); r != nil || err != nil {
			logger.Log().Error(ctx, err.Error())
			tx.Rollback()
		}
	}()

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
	case core.Seeker:
		tableName = "seekers"
	case core.KeeperRole:

		tableName = "keepers"
	case core.Vet:
		tableName = "vets"
	default:
		return core.Role{}, core.ErrInvalidRole
	}

	var count int64
	err = tx.Table(tableName).Where("user_id = ?", id).Count(&count).Error
	if err != nil {
		return core.Role{}, err
	}
	if count == 0 {
		return core.Role{}, core.ErrUserRoleNotFound
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

func (s *store) DeleteUserRole(ctx context.Context, id int, roleName string) (err error) {
	tx := s.DB.WithContext(ctx).Begin()

	if tx.Error != nil {
		logger.Log().Error(ctx, err.Error())
		tx.Rollback()
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil || err != nil {
			logger.Log().Error(ctx, err.Error())
			tx.Rollback()
		}
	}()

	var role core.Role
	switch roleName {
	case core.Seeker:
		err = tx.Table("seekers").Where("user_id = ?", id).Delete(role).Error
	case core.KeeperRole:
		err = tx.Table("keepers").Where("user_id = ?", id).Delete(role).Error
	case core.Vet:
		err = tx.Table("vets").Where("user_id = ?", id).Delete(role).Error
	default:
		return core.ErrInvalidRole
	}

	if err != nil {
		return err
	}

	return tx.Commit().Error

}
