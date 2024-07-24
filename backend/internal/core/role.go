package core

import (
	"context"
)

type (
	Role struct {
		Id          int    `gorm:"primary_key;autoIncrement" db:"id"`
		UserId      int    `db:"user_id"`
		description string `db:"description"`
		createdAt   string `db:"created_at"`
		updatedAt   string `db:"updated_at"`
	}
	RoleService interface {
		GetUserRoles(ctx context.Context, id int) (roles []Role, err error)
		GiveRoleToUser(ctx context.Context, id int, role string) (err error)
		DeleteUserRole(ctx context.Context, id int, role string) (err error)
		UpdateUserRole(ctx context.Context, id int, role string) (err error)
	}
	RoleStore interface {
		GetUserRoles(ctx context.Context, id int) (roles []Role, err error)
		GiveRoleToUser(ctx context.Context, id int, role string) (err error)
		DeleteUserRole(ctx context.Context, id int, role string) (err error)
		UpdateUserRole(ctx context.Context, id int, role string) (err error)
	}
)
