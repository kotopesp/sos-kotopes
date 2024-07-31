package core

import (
	"context"
	"errors"
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
		Name        string    `gorm:"-"`
		ID          int       `gorm:"id"`
		Username    string    `gorm:"user_id"`
		Description string    `gorm:"description"`
		CreatedAt   time.Time `gorm:"created_at"`
		UpdatedAt   time.Time `gorm:"updated_at"`
	}
	GivenRole struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	UpdateRole struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
	}

	// todo
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

// errors
var (
	ErrInvalidRole      = errors.New("invalid role name")
	ErrNoFieldsToUpdate = errors.New("no fields to update")
)

func (r *Role) ToRoleDetails(roleName, username string) RoleDetails {
	if r == nil {
		return RoleDetails{}
	}
	return RoleDetails{
		Name:        roleName,
		ID:          r.ID,
		Username:    username,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}
