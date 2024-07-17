package core

import (
	"context"
)

type (
	Role struct {
		ID   int    `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}

	RoleStore interface {
		GetRoleByName(ctx context.Context, name string) (data Role, err error)
	}
)

// table name in db for gorm
func (Role) TableName() string {
	return "roles"
}
