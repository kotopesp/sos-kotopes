package core

import (
	"context"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/role"
	"time"
)

type (
	Role struct {
		Id          int       `gorm:"primary_key;autoIncrement" db:"id"`
		UserId      int       `db:"user_id"`
		Description string    `db:"description"`
		CreatedAt   time.Time `db:"created_at"`
		UpdatedAt   time.Time `db:"updated_at"`
	}
	Seeker struct {
		Id          int       `gorm:"primary_key;autoIncrement" db:"id"`
		UserId      int       `db:"user_id"`
		Description string    `db:"description"`
		CreatedAt   time.Time `db:"created_at"`
		UpdatedAt   time.Time `db:"updated_at"`
	}
	Keeper struct {
		Id          int    `gorm:"primary_key;autoIncrement" db:"id"`
		UserId      int    `db:"user_id"`
		Description string `db:"description"`
		CreatedAt   string `db:"created_at"`
		UpdatedAt   string `db:"updated_at"`
	}
	Vet struct {
		Id          int    `gorm:"primary_key;autoIncrement" db:"id"`
		UserId      int    `db:"user_id"`
		Description string `db:"description"`
		CreatedAt   string `db:"created_at"`
		UpdatedAt   string `db:"updated_at"`
	}
	RoleService interface {
		GetUserRoles(ctx context.Context, id int) (roles []role.Role, err error)
		GiveRoleToUser(ctx context.Context, id int, role role.GiveRole) (err error)
		DeleteUserRole(ctx context.Context, id int, role string) (err error)
		UpdateUserRole(ctx context.Context, id int, role string) (err error)
	}
	RoleStore interface {
		GetUserRoles(ctx context.Context, id int) (roles []role.Role, err error)
		GiveRoleToUser(ctx context.Context, id int, role role.GiveRole) (err error)
		DeleteUserRole(ctx context.Context, id int, role string) (err error)
		UpdateUserRole(ctx context.Context, id int, role string) (err error)
	}
)
