package core

import (
	"context"
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

	RoleDetails struct {
		ID          int
		Name        string
		UserID      int
		Username    string
		Description string
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}
	GivenRole struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	UpdateRole struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
	}

	RoleService interface {
		GetUserRoles(ctx context.Context, id int) (roles []RoleDetails, err error)
		GiveRoleToUser(ctx context.Context, id int, role GivenRole) (addedRole RoleDetails, err error)
		DeleteUserRole(ctx context.Context, id int, role string) (err error)
		UpdateUserRole(ctx context.Context, id int, role UpdateRole) (updatedRole RoleDetails, err error)
	}
	RoleStore interface {
		GetUserRoles(ctx context.Context, id int) (roles map[string]Role, err error)
		GiveRoleToUser(ctx context.Context, id int, role GivenRole) (addedRole Role, err error)
		DeleteUserRole(ctx context.Context, id int, role string) (err error)
		UpdateUserRole(ctx context.Context, id int, role UpdateRole) (updatedRole Role, err error)
	}
)

const Seeker = "seeker"
const Keeper = "keeper"
const Vet = "vet"
