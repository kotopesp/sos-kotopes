package auth

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kotopesp/sos-kotopes/internal/core"
	mocks "github.com/kotopesp/sos-kotopes/internal/core/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	secret = []byte("secret")

	errInvalidToken = errors.New("invalid token")
)

func validateToken(tokenString string) (id int, username *string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})
	if err != nil {
		return 0, nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		id, ok := claims["id"].(float64)
		if !ok {
			return 0, nil, errInvalidToken
		}

		username, _ := claims["username"].(string)

		return int(id), &username, nil
	}

	return 0, nil, errInvalidToken
}

func TestLoginBasic(t *testing.T) {
	t.Parallel()
	mockUserStore := mocks.NewMockUserStore(t)
	ctx := context.Background()

	authService := New(mockUserStore, core.AuthServiceConfig{
		JWTSecret:            secret,
		AccessTokenLifetime:  2,
		RefreshTokenLifetime: 43800,
	})

	tests := []struct {
		name         string
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

			accessToken, refreshToken, err := authService.LoginBasic(ctx, tt.argUser)
			assert.ErrorIs(t, err, tt.wantErr)
			if err == nil {
				id, username, err := validateToken(*accessToken)
				assert.ErrorIs(t, err, nil)
				assert.Equal(t, tt.mockRetUser.Username, *username)
				assert.Equal(t, tt.mockRetUser.ID, id)

				id, _, err = validateToken(*refreshToken)
				assert.ErrorIs(t, err, nil)
				assert.Equal(t, tt.mockRetUser.ID, id)
			}
		})
	}
}

func TestSignupBasic(t *testing.T) {
	t.Parallel()
	mockUserStore := mocks.NewMockUserStore(t)
	ctx := context.Background()

	authService := New(mockUserStore, core.AuthServiceConfig{
		JWTSecret:            secret,
		AccessTokenLifetime:  2,
		RefreshTokenLifetime: 43800,
	})

	tests := []struct {
		name         string
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
			})).Return(0, tt.mockRetError).Once()

			err := authService.SignupBasic(ctx, tt.argUser)
			assert.ErrorIs(t, tt.wantErr, err)
		})
	}
}

func TestRefresh(t *testing.T) {
	t.Parallel()
	mockUserStore := mocks.NewMockUserStore(t)
	ctx := context.Background()

	authService := New(mockUserStore, core.AuthServiceConfig{
		JWTSecret:            secret,
		AccessTokenLifetime:  2,
		RefreshTokenLifetime: 43800,
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
			mockUserStore.On("GetUserByID", ctx, tt.id).Return(tt.mockRetUser, tt.mockRetError).Once()

			accessToken, err := authService.Refresh(ctx, tt.id)
			assert.ErrorIs(t, err, tt.wantErr)
			if err == nil {
				id, username, err := validateToken(*accessToken)
				assert.ErrorIs(t, err, nil)
				assert.Equal(t, tt.mockRetUser.Username, *username)
				assert.Equal(t, tt.mockRetUser.ID, id)
			}
		})
	}
}

func TestAuthorizeVK(t *testing.T) {
	t.Parallel()
	t.Log("Need to think how to test...")
}
