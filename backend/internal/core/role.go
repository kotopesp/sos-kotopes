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
	Seeker struct {
		ID          int       `gorm:"id"`
		UserID      int       `gorm:"user_id"`
		Description string    `gorm:"description"`
		CreatedAt   time.Time `gorm:"created_at"`
		UpdatedAt   time.Time `gorm:"updated_at"`
	}
	Keeper struct {
		ID          int       `gorm:"id"`
		UserID      int       `gorm:"user_id"`
		Description string    `gorm:"description"`
		CreatedAt   time.Time `gorm:"created_at"`
		UpdatedAt   time.Time `gorm:"updated_at"`
	}
	Vet struct {
		ID          int       `gorm:"id"`
		UserID      int       `gorm:"user_id"`
		Description string    `gorm:"description"`
		CreatedAt   time.Time `gorm:"created_at"`
		UpdatedAt   time.Time `gorm:"updated_at"`
	}
	RoleDetails struct {
		Name        string    `gorm:"-"`
		ID          int       `gorm:"id"`
		UserID      int       `gorm:"user_id"`
		Description string    `gorm:"description"`
		CreatedAt   time.Time `gorm:"created_at"`
		UpdatedAt   time.Time `gorm:"updated_at"`
	}
	GiveRole struct {
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
		GiveRoleToUser(ctx context.Context, id int, role GiveRole) (err error)
		DeleteUserRole(ctx context.Context, id int, role string) (err error)
		UpdateUserRole(ctx context.Context, id int, role UpdateRole) (err error)
	}
	RoleStore interface {
		GetUserRoles(ctx context.Context, id int) (roles map[string]Role, err error)
		GiveRoleToUser(ctx context.Context, id int, role GiveRole) (err error)
		DeleteUserRole(ctx context.Context, id int, role string) (err error)
		UpdateUserRole(ctx context.Context, id int, role UpdateRole) (err error)
	}
)

// errors
var (
	ErrUserDoNotHaveRole = errors.New("user don't have this role")
)

func (r *Role) ToRoleDetails(name string) RoleDetails {
	if r == nil {
		return RoleDetails{}
	}
	return RoleDetails{
		Name:        name,
		ID:          r.ID,
		UserID:      r.UserID,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}
