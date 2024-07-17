package core

import (
	"context"
)

type (
	UserWithRoles struct {
		User   *User
		Keeper *Keeper
		Seeker *Seeker
		Vet    *Vet
	}

	UserWithRolesStore interface {
		AddUserWithRoles(ctx context.Context, data UserWithRoles) error
	}
)
