package core

import (
	"context"
)

type (
	AuthService interface {
		Login(ctx context.Context, user User) (accessToken string, refreshToken string, err error)
		Signup(ctx context.Context, userWithRoles UserWithRoles) error
		Refresh(ctx context.Context, id int) (accessToken string, err error)
	}
)
