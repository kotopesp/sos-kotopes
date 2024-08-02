package auth

import (
	"context"
	"errors"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/internal/core/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestLoginBasic(t *testing.T) {
	mockUserStore := mocks.NewUserStore(t)
	ctx := context.Background()

	authService := New(mockUserStore, core.AuthServiceConfig{
		JWTSecret: []byte("secret"),
	})

	tests := []struct {
		name         string
		ctx          context.Context
		argUser      core.User
		mockRetUser  core.User
		mockRetError error
		wantErr      error
	}{
		{
			name: "success",
			argUser: core.User{
				Username:     "Rondrean",
				PasswordHash: "0sLGcJAm96L6b01AeGbJ",
			},
			mockRetUser: core.User{
				ID:           1,
				Username:     "Rondrean",
				PasswordHash: "$2a$12$u3U2peGqPmD4yk0bJ0h5VOU1woza0F9uauPfAgHcU5gI/NYflKvtm",
			},
		},
		{
			name: "invalid username",
			argUser: core.User{
				Username:     "Bdulka",
				PasswordHash: "0sLGcJAm96L6b01AeGbJ",
			},
			mockRetError: core.ErrNoSuchUser,
			wantErr:      core.ErrInvalidCredentials,
		},
		{
			name: "invalid password",
			argUser: core.User{
				Username:     "Arista",
				PasswordHash: "1gCCMoQWr3b4w1MurMWK",
			},
			mockRetUser: core.User{
				ID:           2,
				Username:     "Arista",
				PasswordHash: "$2a$12$u3U2peGqPmD4yk0bJ0h5VOU1woza0F9uauPfAgHcU5gI/NYflKvtm",
			},
			wantErr: core.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore.On("GetUserByUsername", ctx, tt.argUser.Username).Return(tt.mockRetUser, tt.mockRetError)

			_, _, err := authService.LoginBasic(ctx, tt.argUser)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestSignupBasic(t *testing.T) {
	mockUserStore := mocks.NewUserStore(t)
	ctx := context.Background()

	authService := New(mockUserStore, core.AuthServiceConfig{
		JWTSecret: []byte("secret"),
	})

	tests := []struct {
		name         string
		ctx          context.Context
		argUser      core.User
		mockRetError error
		wantErr      error
	}{
		{
			name: "success",
			argUser: core.User{
				Username:     "Rondrean",
				PasswordHash: "0sLGcJAm96L6b01AeGbJ",
			},
		},
		{
			name: "not unique username",
			argUser: core.User{
				Username:     "Bdulka",
				PasswordHash: "0sLGcJAm96L6b01AeGbJ",
			},
			mockRetError: core.ErrNotUniqueUsername,
			wantErr:      core.ErrNotUniqueUsername,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore.On("AddUser", ctx, mock.MatchedBy(func(user core.User) bool {
				return user.Username == tt.argUser.Username
			})).Return(0, tt.mockRetError)

			err := authService.SignupBasic(ctx, tt.argUser)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestRefresh(t *testing.T) {
	mockUserStore := mocks.NewUserStore(t)
	ctx := context.Background()

	authService := New(mockUserStore, core.AuthServiceConfig{
		JWTSecret: []byte("secret"),
	})

	errInvalidID := errors.New("invalid id")

	tests := []struct {
		name         string
		id           int
		mockRetUser  core.User
		mockRetError error
		wantErr      error
	}{
		{
			name: "success",
			id:   1,
			mockRetUser: core.User{
				ID:           1,
				Username:     "Rondrean",
				PasswordHash: "$2a$12$u3U2peGqPmD4yk0bJ0h5VOU1woza0F9uauPfAgHcU5gI/NYflKvtm",
			},
			mockRetError: nil,
			wantErr:      nil,
		},
		{
			name:         "invalid id",
			id:           -1,
			mockRetError: errInvalidID,
			wantErr:      errInvalidID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore.On("GetUserByID", ctx, tt.id).Return(tt.mockRetUser, tt.mockRetError)

			_, err := authService.Refresh(ctx, tt.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
