package role_store

import (
	"context"
	"errors"
	"fmt"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/postgres"
	"gorm.io/gorm"
)

type Store struct {
	*postgres.Postgres
}

func NewUserStore(pg *postgres.Postgres) core.RoleStore {
	return &Store{pg}
}

func (r *Store) GetUserRoles(ctx context.Context, id int) (roles []core.Role, err error) {

	var seeker core.Role
	if err := r.DB.WithContext(ctx).Table("seekers").
		Where("user_id = ?", id).
		Select("'seeker' as name, id, user_id, description, created_at, updated_at").
		First(&seeker).Error; err == nil {
		roles = append(roles, seeker)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error querying seeker role: %w", err)
	}

	var keeper core.Role
	if err := r.DB.WithContext(ctx).Table("keepers").
		Where("user_id = ?", id).
		Select("'keeper' as name, id, user_id, description, created_at, updated_at").
		First(&keeper).Error; err == nil {
		roles = append(roles, keeper)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error querying keeper role: %w", err)
	}

	var vet core.Role
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

func (r *Store) GiveRoleToUser(ctx context.Context, id int, role string) (err error) {
	return nil
}
func (r *Store) DeleteUserRole(ctx context.Context, id int, role string) (err error) {
	return nil
}
func (r *Store) UpdateUserRole(ctx context.Context, id int, role string) (err error) {
	return nil
}
