package auth

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/internal/core/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

var (
	secret                  = []byte("secret")
	errInvalidToken         = errors.New("invalid token")
	errCountSessions        = errors.New("count sessions error")
	errUpdateRefreshSession = errors.New("update refresh session error")
	errGetUserByUsername    = errors.New("get user by username error")
	errGerUserByID          = errors.New("get user by id error")
)

const (
	accessTokenLifetime      = 2
	refreshTokenLifetime     = 10
	longRefreshTokenLifetime = 200
	username                 = "JackVorobey"
	password                 = "0sLGcJAm96L6b01AeGbJ"
	passwordHash             = "$2a$12$u3U2peGqPmD4yk0bJ0h5VOU1woza0F9uauPfAgHcU5gI/NYflKvtm"
	invalidPassword          = "invalid"
	refreshToken1            = "b49b5443-f9d3-44b5-829c-7b524fdc92d4"
	refreshToken2            = "c0c9723c-3723-4fdb-9f67-b6365839e526"
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
	mockUserStore := mocks.NewUserStore(t)
	mockRefreshSessionStore := mocks.NewRefreshSessionStore(t)

	ctx := context.Background()

	authService := New(mockUserStore, mockRefreshSessionStore, core.AuthServiceConfig{
		JWTSecret:            secret,
		AccessTokenLifetime:  accessTokenLifetime,
		RefreshTokenLifetime: refreshTokenLifetime,
	})

	tests := []struct {
		name                         string
		getUserByUsernameArg2        string
		getUserByUsernameRet1        core.User
		getUserByUsernameRet2        error
		invokeGetUserByUsername      bool
		countSessionsAndDeleteArg2   int
		countSessionsAndDeleteRet1   error
		invokeCountSessionsAndDelete bool
		updateRefreshSessionArg2     core.RefreshSession
		updateRefreshSessionRet1     error
		invokeUpdateRefreshSession   bool
		loginBasicArg2               core.User
		wantErr                      error
	}{
		{
			name:                  "success",
			getUserByUsernameArg2: username,
			getUserByUsernameRet1: core.User{
				ID:           1,
				Username:     username,
				PasswordHash: passwordHash,
			},
			invokeGetUserByUsername:      true,
			countSessionsAndDeleteArg2:   1,
			invokeCountSessionsAndDelete: true,
			updateRefreshSessionArg2: core.RefreshSession{
				UserID:    1,
				ExpiresAt: time.Now().Add(time.Minute * time.Duration(refreshTokenLifetime)),
			},
			invokeUpdateRefreshSession: true,
			loginBasicArg2: core.User{
				Username:     username,
				PasswordHash: password,
			},
		},
		{
			name:                    "username not exists",
			getUserByUsernameArg2:   username,
			getUserByUsernameRet2:   core.ErrNoSuchUser,
			invokeGetUserByUsername: true,
			loginBasicArg2: core.User{
				Username: username,
			},
			wantErr: core.ErrInvalidCredentials,
		},
		{
			name:                  "invalid password",
			getUserByUsernameArg2: username,
			getUserByUsernameRet1: core.User{
				ID:           1,
				Username:     username,
				PasswordHash: passwordHash,
			},
			invokeGetUserByUsername: true,
			loginBasicArg2: core.User{
				Username:     username,
				PasswordHash: invalidPassword,
			},
			wantErr: core.ErrInvalidCredentials,
		},
		{
			name:                    "get user by username error",
			getUserByUsernameArg2:   username,
			getUserByUsernameRet2:   errGetUserByUsername,
			invokeGetUserByUsername: true,
			loginBasicArg2: core.User{
				Username:     username,
				PasswordHash: password,
			},
			wantErr: errGetUserByUsername,
		},
		{
			name:                  "count sessions error",
			getUserByUsernameArg2: username,
			getUserByUsernameRet1: core.User{
				ID:           1,
				Username:     username,
				PasswordHash: passwordHash,
			},
			invokeGetUserByUsername:      true,
			countSessionsAndDeleteArg2:   1,
			countSessionsAndDeleteRet1:   errCountSessions,
			invokeCountSessionsAndDelete: true,
			loginBasicArg2: core.User{
				Username:     username,
				PasswordHash: password,
			},
			wantErr: errCountSessions,
		},
		{
			name:                  "update refresh session error",
			getUserByUsernameArg2: username,
			getUserByUsernameRet1: core.User{
				ID:           1,
				Username:     username,
				PasswordHash: passwordHash,
			},
			invokeGetUserByUsername:      true,
			countSessionsAndDeleteArg2:   1,
			invokeCountSessionsAndDelete: true,
			updateRefreshSessionArg2: core.RefreshSession{
				UserID:    1,
				ExpiresAt: time.Now().Add(time.Minute * time.Duration(refreshTokenLifetime)),
			},
			updateRefreshSessionRet1:   errUpdateRefreshSession,
			invokeUpdateRefreshSession: true,
			loginBasicArg2: core.User{
				Username:     username,
				PasswordHash: password,
			},
			wantErr: errUpdateRefreshSession,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.invokeGetUserByUsername {
				mockUserStore.
					On("GetUserByUsername", ctx, tt.getUserByUsernameArg2).
					Return(tt.getUserByUsernameRet1, tt.getUserByUsernameRet2).Once()
			}
			if tt.invokeCountSessionsAndDelete {
				mockRefreshSessionStore.
					On("CountSessionsAndDelete", ctx, tt.countSessionsAndDeleteArg2).
					Return(tt.countSessionsAndDeleteRet1).Once()
			}
			if tt.invokeUpdateRefreshSession {
				mockRefreshSessionStore.
					On(
						"UpdateRefreshSession",
						ctx,
						mock.AnythingOfType("core.UpdateRefreshSessionParam"),
						mock.MatchedBy(func(rs core.RefreshSession) bool {
							return rs.UserID == tt.updateRefreshSessionArg2.UserID &&
								rs.ExpiresAt.Sub(tt.updateRefreshSessionArg2.ExpiresAt) < time.Second
						}),
					).Return(tt.updateRefreshSessionRet1).Once()
			}

			accessToken, _, err := authService.LoginBasic(ctx, tt.loginBasicArg2)
			assert.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr == nil {
				id, _, err := validateToken(*accessToken)
				assert.NoError(t, err)
				assert.Equal(t, tt.getUserByUsernameRet1.ID, id)
			}
		})
	}
}

func TestSignupBasic(t *testing.T) {
	mockUserStore := mocks.NewUserStore(t)
	mockRefreshSessionStore := mocks.NewRefreshSessionStore(t)

	ctx := context.Background()

	authService := New(mockUserStore, mockRefreshSessionStore, core.AuthServiceConfig{
		JWTSecret:            secret,
		AccessTokenLifetime:  accessTokenLifetime,
		RefreshTokenLifetime: refreshTokenLifetime,
	})

	tests := []struct {
		name            string
		addUserArg2     core.User
		addUserRet2     error
		invokeAddUser   bool
		signupBasicArg2 core.User
		wantErr         error
	}{
		{
			name: "success",
			addUserArg2: core.User{
				Username: username,
			},
			invokeAddUser: true,
			signupBasicArg2: core.User{
				Username:     username,
				PasswordHash: password,
			},
		},
		{
			name: "not unique username",
			addUserArg2: core.User{
				Username: username,
			},
			addUserRet2:   core.ErrNotUniqueUsername,
			invokeAddUser: true,
			signupBasicArg2: core.User{
				Username:     username,
				PasswordHash: password,
			},
			wantErr: core.ErrNotUniqueUsername,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.invokeAddUser {
				mockUserStore.On(
					"AddUser",
					ctx,
					mock.MatchedBy(func(user core.User) bool {
						return user.Username == tt.addUserArg2.Username &&
							bcrypt.CompareHashAndPassword(
								[]byte(user.PasswordHash),
								[]byte(tt.signupBasicArg2.PasswordHash),
							) == nil
					}),
				).Return(1, tt.addUserRet2).Once()
			}

			err := authService.SignupBasic(ctx, tt.signupBasicArg2)
			assert.ErrorIs(t, tt.wantErr, err)
		})
	}
}

func TestRefresh(t *testing.T) {
	mockUserStore := mocks.NewUserStore(t)
	mockRefreshSessionStore := mocks.NewRefreshSessionStore(t)

	ctx := context.Background()

	authService := New(mockUserStore, mockRefreshSessionStore, core.AuthServiceConfig{
		JWTSecret:            secret,
		AccessTokenLifetime:  accessTokenLifetime,
		RefreshTokenLifetime: refreshTokenLifetime,
	})

	tests := []struct {
		name                           string
		getRefreshSessionByTokenArg2   string
		getRefreshSessionByTokenRet1   core.RefreshSession
		getRefreshSessionByTokenRet2   error
		invokeGetRefreshSessionByToken bool
		getUserByIDRet1                core.User
		getUserByIDRet2                error
		invokeGetUserByID              bool
		updateRefreshSessionArg3       core.RefreshSession
		updateRefreshSessionRet1       error
		invokeUpdateRefreshSession     bool
		refreshArg2                    core.RefreshSession
		wantErr                        error
	}{
		{
			name:                         "success",
			getRefreshSessionByTokenArg2: refreshToken1,
			getRefreshSessionByTokenRet1: core.RefreshSession{
				ID:           1,
				UserID:       1,
				RefreshToken: refreshToken2,
				ExpiresAt:    time.Now().Add(time.Minute * time.Duration(longRefreshTokenLifetime)),
			},
			invokeGetRefreshSessionByToken: true,
			getUserByIDRet1: core.User{
				ID:       1,
				Username: username,
			},
			invokeGetUserByID: true,
			updateRefreshSessionArg3: core.RefreshSession{
				UserID: 1,
			},
			invokeUpdateRefreshSession: true,
			refreshArg2: core.RefreshSession{
				RefreshToken: refreshToken1,
			},
		},
		{
			name:                           "token does not exists",
			getRefreshSessionByTokenArg2:   refreshToken1,
			getRefreshSessionByTokenRet2:   errors.New("token not found"),
			invokeGetRefreshSessionByToken: true,
			refreshArg2: core.RefreshSession{
				RefreshToken: refreshToken1,
			},
			wantErr: core.ErrUnauthorized,
		},
		{
			name:                         "get user by id error",
			getRefreshSessionByTokenArg2: refreshToken1,
			getRefreshSessionByTokenRet1: core.RefreshSession{
				ID:           1,
				UserID:       1,
				RefreshToken: refreshToken2,
				ExpiresAt:    time.Now().Add(time.Minute * time.Duration(longRefreshTokenLifetime)),
			},
			invokeGetRefreshSessionByToken: true,
			getUserByIDRet2:                errGerUserByID,
			invokeGetUserByID:              true,
			refreshArg2: core.RefreshSession{
				RefreshToken: refreshToken1,
			},
			wantErr: errGerUserByID,
		},
		{
			name:                         "expired token",
			getRefreshSessionByTokenArg2: refreshToken1,
			getRefreshSessionByTokenRet1: core.RefreshSession{
				ID:           1,
				UserID:       1,
				RefreshToken: refreshToken2,
				ExpiresAt:    time.Now().Add(-time.Minute),
			},
			invokeGetRefreshSessionByToken: true,
			invokeGetUserByID:              true,
			refreshArg2: core.RefreshSession{
				RefreshToken: refreshToken1,
			},
			wantErr: core.ErrUnauthorized,
		},
		{
			name:                         "update refresh session error",
			getRefreshSessionByTokenArg2: refreshToken1,
			getRefreshSessionByTokenRet1: core.RefreshSession{
				ID:           1,
				UserID:       1,
				RefreshToken: refreshToken2,
				ExpiresAt:    time.Now().Add(time.Minute * time.Duration(longRefreshTokenLifetime)),
			},
			invokeGetRefreshSessionByToken: true,
			getUserByIDRet1: core.User{
				ID:       1,
				Username: username,
			},
			invokeGetUserByID: true,
			updateRefreshSessionArg3: core.RefreshSession{
				UserID: 1,
			},
			updateRefreshSessionRet1:   errUpdateRefreshSession,
			invokeUpdateRefreshSession: true,
			refreshArg2: core.RefreshSession{
				RefreshToken: refreshToken1,
			},
			wantErr: errUpdateRefreshSession,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.invokeGetUserByID {
				mockUserStore.On(
					"GetUserByID",
					ctx,
					tt.getRefreshSessionByTokenRet1.UserID,
				).Return(tt.getUserByIDRet1, tt.getUserByIDRet2).Once()
			}
			if tt.invokeGetRefreshSessionByToken {
				mockRefreshSessionStore.On(
					"GetRefreshSessionByToken",
					ctx,
					tt.getRefreshSessionByTokenArg2,
				).Return(tt.getRefreshSessionByTokenRet1, tt.getRefreshSessionByTokenRet2).Once()
			}
			if tt.invokeUpdateRefreshSession {
				mockRefreshSessionStore.On(
					"UpdateRefreshSession",
					ctx,
					mock.AnythingOfType("core.UpdateRefreshSessionParam"),
					mock.MatchedBy(func(rs core.RefreshSession) bool {
						return rs.UserID == tt.updateRefreshSessionArg3.UserID &&
							rs.RefreshToken != tt.refreshArg2.RefreshToken
					}),
				).Return(tt.updateRefreshSessionRet1).Once()
			}

			accessToken, _, err := authService.Refresh(ctx, tt.refreshArg2)
			assert.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr == nil {
				id, _, err := validateToken(*accessToken)
				assert.ErrorIs(t, err, nil)
				assert.Equal(t, tt.getUserByIDRet1.ID, id)
			}
		})
	}
}
