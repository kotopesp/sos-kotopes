package core

import (
	"context"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/role"
	"time"
)

type (
	Role struct {
		ID          int       `gorm:"id"`
		UserID      int       `gorm:"user_id"`
		Description string    `gorm:"description"`
		CreatedAt   time.Time `gorm:"created_at"`
		UpdatedAt   time.Time `gorm:"updated_at"`
	}
	Seeker struct {
		ID          int       `gorm:"id"`
		UserID      int       `gorm:"user_id"`
		Description string    `gorm:"description"`
		CreatedAt   time.Time `gorm:"created_at"`
		UpdatedAt   time.Time `gorm:"updated_at"`
	}
	Keeper struct {
		ID          int    `gorm:"id"`
		UserID      int    `gorm:"user_id"`
		Description string `gorm:"description"`
		CreatedAt   string `gorm:"created_at"`
		UpdatedAt   string `gorm:"updated_at"`
	}
	Vet struct {
		ID          int    `gorm:"id"`
		UserID      int    `gorm:"user_id"`
		Description string `gorm:"description"`
		CreatedAt   string `gorm:"created_at"`
		UpdatedAt   string `gorm:"updated_at"`
	}

	// todo
	RoleService interface {
		GetUserRoles(ctx context.Context, id int) (roles []Role, err error)
		GiveRoleToUser(ctx context.Context, id int, role role.GiveRole) (err error)
		DeleteUserRole(ctx context.Context, id int, role string) (err error)
		UpdateUserRole(ctx context.Context, id int, role role.UpdateRole) (err error)
	}
	RoleStore interface {
		GetUserRoles(ctx context.Context, id int) (roles []Role, err error)
		GiveRoleToUser(ctx context.Context, id int, role role.GiveRole) (err error)
		DeleteUserRole(ctx context.Context, id int, role string) (err error)
		UpdateUserRole(ctx context.Context, id int, role role.UpdateRole) (err error)
	}
)
