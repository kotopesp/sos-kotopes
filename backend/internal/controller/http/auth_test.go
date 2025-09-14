package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/user"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/validator"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	correctUsername       = "JackVorobey"                                                                                                                                                     // nolint:all
	correctPassword       = "VeryStrongPassword123"                                                                                                                                           // nolint:all
	invalidPassword       = "WrongPassword123"                                                                                                                                                // nolint:all
	emptyPassword         = ""                                                                                                                                                                // nolint:all
	smallPassword         = "pa$$"                                                                                                                                                            // nolint:all
	bigPassword           = "password_that_have_more_than_72_symbols,password_that_have_more_than_72_symbols,password_that_have_more_than_72_symbols,password_that_have_more_than_72_symbols" // nolint:all
	omitDigitPassword     = "StrongPasswordButWithoutDigit"                                                                                                                                   // nolint:all
	omitUppercasePassword = "strong_but_omit_uppercase123"                                                                                                                                    // nolint:all
	emptyUsername         = ""                                                                                                                                                                // nolint:all
	veryBigUsername       = "JackVorobeyJackVorobeyJackVorobeyJackVorobeyJackVorobeyJackVorobey"                                                                                              // nolint:all
	withSpecialsUsername  = "Jack_Vorobey"                                                                                                                                                    // nolint:all
)

var (
	firstname    = "Jack"                                                                                                                                      // nolint:all
	lastname     = "Vorobey"                                                                                                                                   // nolint:all
	description  = "Description"                                                                                                                               // nolint:all
	accessToken  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkphY2sgVm9yb2JleSJ9.-YQEVClcn8V-dwlqFGGV5NWSzHAYPytwfxQXy8sdz5M" // nolint:all
	refreshToken = "629e239d-7351-4440-b6c0-185d73f58a65"                                                                                                      // nolint:all
)

type TestAccessTokenSuccessResponse struct {
	Data string `json:"data"`
}

type TestValidationErrorResponse struct {
	Data validator.Response `json:"data"`
}

func stringPtr(s string) *string {
	return &s
}

func TestLoginBasic(t *testing.T) {
	t.Parallel()
	app, dependencies := newTestApp(t)

	const route = "/api/v1/auth/login"

	tests := []struct {
		name               string
		loginBasicArg2     user.User
		loginBasicRet1     string
		loginBasicRet2     string
		loginBasicRet3     error
		wantValidationErrs []validator.ResponseError
		wantCode           int
	}{
		{
			name: "success",
			loginBasicArg2: user.User{
				Username: correctUsername,
				Password: correctPassword,
			},
			loginBasicRet1: accessToken,
			loginBasicRet2: refreshToken,
			wantCode:       http.StatusOK,
		},
		{
			name: "empty password",
			loginBasicArg2: user.User{
				Username: correctUsername,
				Password: emptyPassword,
			},
			wantValidationErrs: []validator.ResponseError{
				{
					FailedField: "Password",
					Tag:         "required",
					Value:       emptyPassword,
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "small password",
			loginBasicArg2: user.User{
				Username: correctUsername,
				Password: smallPassword,
			},
			wantValidationErrs: []validator.ResponseError{
				{
					FailedField: "Password",
					Tag:         "min",
					Param:       "8",
					Value:       smallPassword,
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "very big password",
			loginBasicArg2: user.User{
				Username: correctUsername,
				Password: bigPassword,
			},
			wantValidationErrs: []validator.ResponseError{
				{
					FailedField: "Password",
					Tag:         "max",
					Param:       "72",
					Value:       bigPassword,
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "omit digit",
			loginBasicArg2: user.User{
				Username: "JackVorobey",
				Password: omitDigitPassword,
			},
			wantValidationErrs: []validator.ResponseError{
				{
					FailedField: "Password",
					Tag:         "contains_digit",
					Value:       omitDigitPassword,
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "omit uppercase",
			loginBasicArg2: user.User{
				Username: "JackVorobey",
				Password: omitUppercasePassword,
			},
			wantValidationErrs: []validator.ResponseError{
				{
					FailedField: "Password",
					Tag:         "contains_uppercase",
					Value:       omitUppercasePassword,
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty username",
			loginBasicArg2: user.User{
				Username: emptyUsername,
				Password: correctPassword,
			},
			wantValidationErrs: []validator.ResponseError{
				{
					FailedField: "Username",
					Tag:         "required",
					Value:       emptyUsername,
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "very big username",
			loginBasicArg2: user.User{
				Username: veryBigUsername,
				Password: correctPassword,
			},
			wantValidationErrs: []validator.ResponseError{
				{
					FailedField: "Username",
					Tag:         "max",
					Param:       "50",
					Value:       veryBigUsername,
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "username with specials",
			loginBasicArg2: user.User{
				Username: withSpecialsUsername,
				Password: correctPassword,
			},
			wantValidationErrs: []validator.ResponseError{
				{
					FailedField: "Username",
					Tag:         "no_specials",
					Value:       withSpecialsUsername,
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid credentials",
			loginBasicArg2: user.User{
				Username: correctPassword,
				Password: invalidPassword,
			},
			loginBasicRet3: core.ErrInvalidCredentials,
			wantCode:       http.StatusUnauthorized,
		},
		{
			name: "user is banned",
			loginBasicArg2: user.User{
				Username: correctUsername,
				Password: correctPassword,
			},
			loginBasicRet3: core.ErrUserIsBanned,
			wantCode:       http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantValidationErrs == nil {
				at := tt.loginBasicRet1
				rt := tt.loginBasicRet2
				dependencies.authService.On(
					"LoginBasic",
					mock.Anything,
					tt.loginBasicArg2.ToCoreUser(),
				).Return(&at, &rt, tt.loginBasicRet3).Once()
			}

			reqBody := new(bytes.Buffer)
			mp := multipart.NewWriter(reqBody)
			_ = mp.WriteField("username", tt.loginBasicArg2.Username)
			_ = mp.WriteField("password", tt.loginBasicArg2.Password)

			_ = mp.Close()

			req := httptest.NewRequest(http.MethodPost, route, reqBody)
			req.Header["Content-Type"] = []string{mp.FormDataContentType()}

			resp, _ := app.Test(req, -1)
			assert.Equal(t, tt.wantCode, resp.StatusCode)

			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode == http.StatusOK {
				var testTokenSuccessResponse TestAccessTokenSuccessResponse
				_ = json.Unmarshal(body, &testTokenSuccessResponse)

				assert.Equal(t, tt.loginBasicRet1, testTokenSuccessResponse.Data)

				cookies := resp.Cookies()
				hasToken := false
				for _, cookie := range cookies {
					if cookie.Name == "refresh_token" {
						hasToken = true
						assert.Equal(t, tt.loginBasicRet2, cookie.Value)
						break
					}
				}
				assert.Equal(t, true, hasToken)
			} else if resp.StatusCode == http.StatusUnprocessableEntity {
				var testValidationErrorResponse TestValidationErrorResponse
				_ = json.Unmarshal(body, &testValidationErrorResponse)

				assert.Equal(t, tt.wantValidationErrs, testValidationErrorResponse.Data.ValidationErrors)
			}
		})
	}
}

func TestSignup(t *testing.T) {
	t.Parallel()
	app, dependencies := newTestApp(t)

	const route = "/api/v1/auth/signup"

	tests := []struct {
		name               string
		signupBasicArg2    user.User
		signupBasicRet1    error
		invokeSignupBasic  bool
		wantValidationErrs []validator.ResponseError
		wantCode           int
	}{
		{
			name: "success",
			signupBasicArg2: user.User{
				Username:    correctUsername,
				Password:    correctPassword,
				Firstname:   &firstname,
				Lastname:    &lastname,
				Description: &description,
			},
			invokeSignupBasic: true,
			wantCode:          http.StatusCreated,
		},
		{
			name: "empty password",
			signupBasicArg2: user.User{
				Username: correctUsername,
				Password: emptyPassword,
			},
			wantValidationErrs: []validator.ResponseError{
				{
					FailedField: "Password",
					Tag:         "required",
					Value:       emptyPassword,
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "small password",
			signupBasicArg2: user.User{
				Username: correctUsername,
				Password: smallPassword,
			},
			wantValidationErrs: []validator.ResponseError{
				{
					FailedField: "Password",
					Tag:         "min",
					Param:       "8",
					Value:       smallPassword,
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "very big password",
			signupBasicArg2: user.User{
				Username: correctUsername,
				Password: bigPassword,
			},
			wantValidationErrs: []validator.ResponseError{
				{
					FailedField: "Password",
					Tag:         "max",
					Param:       "72",
					Value:       bigPassword,
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "omit digit",
			signupBasicArg2: user.User{
				Username: correctUsername,
				Password: omitDigitPassword,
			},
			wantValidationErrs: []validator.ResponseError{
				{
					FailedField: "Password",
					Tag:         "contains_digit",
					Value:       omitDigitPassword,
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "omit uppercase",
			signupBasicArg2: user.User{
				Username: correctUsername,
				Password: omitUppercasePassword,
			},
			wantValidationErrs: []validator.ResponseError{
				{
					FailedField: "Password",
					Tag:         "contains_uppercase",
					Value:       omitUppercasePassword,
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty username",
			signupBasicArg2: user.User{
				Username: emptyUsername,
				Password: correctPassword,
			},
			wantValidationErrs: []validator.ResponseError{
				{
					FailedField: "Username",
					Tag:         "required",
					Value:       emptyUsername,
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "very big username",
			signupBasicArg2: user.User{
				Username: veryBigUsername,
				Password: correctPassword,
			},
			wantValidationErrs: []validator.ResponseError{
				{
					FailedField: "Username",
					Tag:         "max",
					Param:       "50",
					Value:       veryBigUsername,
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "username with specials",
			signupBasicArg2: user.User{
				Username: withSpecialsUsername,
				Password: correctPassword,
			},
			wantValidationErrs: []validator.ResponseError{
				{
					FailedField: "Username",
					Tag:         "no_specials",
					Value:       withSpecialsUsername,
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "not unique username",
			signupBasicArg2: user.User{
				Username: correctUsername,
				Password: correctPassword,
			},
			invokeSignupBasic: true,
			wantValidationErrs: []validator.ResponseError{
				{
					FailedField: "Username",
					Tag:         "unique",
					Value:       correctUsername,
				},
			},
			signupBasicRet1: core.ErrNotUniqueUsername,
			wantCode:        http.StatusUnprocessableEntity,
		},
		{
			name: "internal server error",
			signupBasicArg2: user.User{
				Username: correctUsername,
				Password: correctPassword,
			},
			invokeSignupBasic: true,
			signupBasicRet1:   errors.New("internal server error"),
			wantCode:          http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.invokeSignupBasic {
				dependencies.authService.On(
					"SignupBasic",
					mock.Anything,
					tt.signupBasicArg2.ToCoreUser(),
				).Return(tt.signupBasicRet1).Once()
			}

			reqBody := new(bytes.Buffer)
			mp := multipart.NewWriter(reqBody)
			_ = mp.WriteField("username", tt.signupBasicArg2.Username)
			_ = mp.WriteField("password", tt.signupBasicArg2.Password)
			if tt.signupBasicArg2.Firstname != nil {
				_ = mp.WriteField("firstname", *tt.signupBasicArg2.Firstname)
			}
			if tt.signupBasicArg2.Lastname != nil {
				_ = mp.WriteField("lastname", *tt.signupBasicArg2.Lastname)
			}
			if tt.signupBasicArg2.Description != nil {
				_ = mp.WriteField("description", *tt.signupBasicArg2.Description)
			}

			_ = mp.Close()

			req := httptest.NewRequest(http.MethodPost, route, reqBody)
			req.Header["Content-Type"] = []string{mp.FormDataContentType()}

			resp, _ := app.Test(req, -1)
			assert.Equal(t, tt.wantCode, resp.StatusCode)

			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

			var testValidationErrorResponse TestValidationErrorResponse
			_ = json.Unmarshal(body, &testValidationErrorResponse)
			assert.Equal(t, tt.wantValidationErrs, testValidationErrorResponse.Data.ValidationErrors)
		})
	}
}

func TestRefresh(t *testing.T) {
	t.Parallel()
	app, dependencies := newTestApp(t)

	const route = "/api/v1/auth/token/refresh"

	tests := []struct {
		name               string
		cookieRefreshToken string
		refreshRet1        *string
		refreshRet2        *string
		refreshRet3        error
		wantCode           int
	}{
		{
			name:               "success",
			cookieRefreshToken: refreshToken,
			refreshRet1:        &accessToken,
			refreshRet2:        &refreshToken,
			wantCode:           http.StatusOK,
		},
		{
			name:               "invalid token",
			cookieRefreshToken: refreshToken,
			refreshRet3:        core.ErrUnauthorized,
			wantCode:           http.StatusUnauthorized,
		},
		{
			name:               "internal server error",
			cookieRefreshToken: refreshToken,
			refreshRet3:        errors.New("internal server error"),
			wantCode:           http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dependencies.authService.On(
				"Refresh",
				mock.Anything,
				core.RefreshSession{
					RefreshToken: tt.cookieRefreshToken,
				},
			).Return(tt.refreshRet1, tt.refreshRet2, tt.refreshRet3).Once()

			req := httptest.NewRequest(http.MethodPost, route, http.NoBody)
			req.AddCookie(&http.Cookie{
				Name:  "refresh_token",
				Value: tt.cookieRefreshToken,
			})

			resp, _ := app.Test(req, -1)
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

			assert.Equal(t, tt.wantCode, resp.StatusCode)

			var testAccessTokenSuccessResponse TestAccessTokenSuccessResponse
			_ = json.Unmarshal(body, &testAccessTokenSuccessResponse)
			if tt.wantCode == http.StatusOK {
				assert.Equal(t, *tt.refreshRet1, testAccessTokenSuccessResponse.Data)
			}
		})
	}
}
