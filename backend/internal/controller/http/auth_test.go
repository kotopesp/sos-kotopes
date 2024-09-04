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

type TestAccessTokenSuccessResponse struct {
	Data string `json:"data"`
}

type TestValidationErrorResponse struct {
	Data struct {
		ValidationErrors []validator.ResponseError `json:"validation_errors"`
	} `json:"data"`
}

func stringPtr(s string) *string {
	return &s
}

func TestLoginBasic(t *testing.T) {
	app, dependencies := newTestApp(t)

	const route = "/api/v1/auth/login"

	tests := []struct {
		name                string
		mockArgUser         user.User
		mockRetAccessToken  string
		mockRetRefreshToken string
		mockRetError        error
		wantErrs            []validator.ResponseError
		wantCode            int
	}{
		{
			name: "success",
			mockArgUser: user.User{
				Username: "JackVorobey",
				Password: "VeryStrongPassword123",
			},
			mockRetAccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkphY2sgVm9yb2JleSJ9.-YQEVClcn8V-dwlqFGGV5NWSzHAYPytwfxQXy8sdz5M",
			mockRetRefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIyMzAzODkwMzg5MDMiLCJuYW1lIjoiSmFjayBWb3JvYmV5In0.5TH7KYwBR_zkFf9VUsB_51U4ThD81vQTW-wMCp-hnco",
			wantCode:            http.StatusOK,
		},
		{
			name: "empty password",
			mockArgUser: user.User{
				Username: "JackVorobey",
				Password: "",
			},
			wantErrs: []validator.ResponseError{
				{
					FailedField: "Password",
					Tag:         "required",
					Value:       "",
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "small password",
			mockArgUser: user.User{
				Username: "JackVorobey",
				Password: "pa$$",
			},
			wantErrs: []validator.ResponseError{
				{
					FailedField: "Password",
					Tag:         "min",
					Param:       "8",
					Value:       "pa$$",
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "very big password",
			mockArgUser: user.User{
				Username: "JackVorobey",
				Password: "password_that_have_more_than_72_symbols,password_that_have_more_than_72_symbols,password_that_have_more_than_72_symbols,password_that_have_more_than_72_symbols",
			},
			wantErrs: []validator.ResponseError{
				{
					FailedField: "Password",
					Tag:         "max",
					Param:       "72",
					Value:       "password_that_have_more_than_72_symbols,password_that_have_more_than_72_symbols,password_that_have_more_than_72_symbols,password_that_have_more_than_72_symbols",
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "omit digit",
			mockArgUser: user.User{
				Username: "JackVorobey",
				Password: "StrongPasswordButWithoutDigit",
			},
			wantErrs: []validator.ResponseError{
				{
					FailedField: "Password",
					Tag:         "contains_digit",
					Value:       "StrongPasswordButWithoutDigit",
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "omit uppercase",
			mockArgUser: user.User{
				Username: "JackVorobey",
				Password: "strong_but_omit_uppercase123",
			},
			wantErrs: []validator.ResponseError{
				{
					FailedField: "Password",
					Tag:         "contains_uppercase",
					Value:       "strong_but_omit_uppercase123",
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty username",
			mockArgUser: user.User{
				Username: "",
				Password: "VeryStrongPassword123",
			},
			wantErrs: []validator.ResponseError{
				{
					FailedField: "Username",
					Tag:         "required",
					Value:       "",
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "very big username",
			mockArgUser: user.User{
				Username: "JackVorobeyJackVorobeyJackVorobeyJackVorobeyJackVorobeyJackVorobey",
				Password: "VeryStrongPassword123",
			},
			wantErrs: []validator.ResponseError{
				{
					FailedField: "Username",
					Tag:         "max",
					Param:       "50",
					Value:       "JackVorobeyJackVorobeyJackVorobeyJackVorobeyJackVorobeyJackVorobey",
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "username with specials",
			mockArgUser: user.User{
				Username: "Jack_Vorobey",
				Password: "VeryStrongPassword123",
			},
			wantErrs: []validator.ResponseError{
				{
					FailedField: "Username",
					Tag:         "no_specials",
					Value:       "Jack_Vorobey",
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid credentials",
			mockArgUser: user.User{
				Username: "JackVorobeyWrongCredentials",
				Password: "WrongPassword123",
			},
			mockRetError: core.ErrInvalidCredentials,
			wantCode:     http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErrs == nil {
				at := tt.mockRetAccessToken
				rt := tt.mockRetRefreshToken
				dependencies.authService.On("LoginBasic", mock.Anything, tt.mockArgUser.ToCoreUser()).Return(&at, &rt, tt.mockRetError).Once()
			}

			// request body (multipart)
			reqBody := new(bytes.Buffer)
			mp := multipart.NewWriter(reqBody)
			_ = mp.WriteField("username", tt.mockArgUser.Username)
			_ = mp.WriteField("password", tt.mockArgUser.Password)

			err := mp.Close()
			if err != nil {
				t.Fatal(err)
			}

			// request
			req := httptest.NewRequest(http.MethodPost, route, reqBody)
			req.Header["Content-Type"] = []string{mp.FormDataContentType()}

			// doing request
			resp, _ := app.Test(req, -1)
			assert.Equal(t, tt.wantCode, resp.StatusCode)

			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode == http.StatusOK {
				var testTokenSuccessResponse TestAccessTokenSuccessResponse
				_ = json.Unmarshal(body, &testTokenSuccessResponse)

				assert.Equal(t, tt.mockRetAccessToken, testTokenSuccessResponse.Data)

				cookies := resp.Cookies()
				hasToken := false
				for _, cookie := range cookies {
					if cookie.Name == "refresh_token" {
						hasToken = true
						assert.Equal(t, tt.mockRetRefreshToken, cookie.Value)
						break
					}
				}
				assert.Equal(t, true, hasToken)
			} else if resp.StatusCode == http.StatusUnprocessableEntity {
				var testValidationErrorResponse TestValidationErrorResponse
				_ = json.Unmarshal(body, &testValidationErrorResponse)

				assert.Equal(t, tt.wantErrs, testValidationErrorResponse.Data.ValidationErrors)
			}
		})
	}
}

func TestSignup(t *testing.T) {
	app, dependencies := newTestApp(t)

	const route = "/api/v1/auth/signup"

	var (
		FirstName   = "Jack"
		LastName    = "Jackson"
		Description = "Elephants are the largest land animals on Earth."
	)

	tests := []struct {
		name         string
		mockArgUser  user.User
		mockRetError error
		wantErrs     []validator.ResponseError
		wantCode     int
	}{
		{
			name: "success",
			mockArgUser: user.User{
				Username:    "JackVorobey",
				Password:    "VeryStrongPassword123",
				Firstname:   &FirstName,
				Lastname:    &LastName,
				Description: &Description,
			},
			wantCode: http.StatusCreated,
		},
		{
			name: "empty password",
			mockArgUser: user.User{
				Username: "JackVorobey",
				Password: "",
			},
			wantErrs: []validator.ResponseError{
				{
					FailedField: "Password",
					Tag:         "required",
					Value:       "",
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "small password",
			mockArgUser: user.User{
				Username: "JackVorobey",
				Password: "pa$$",
			},
			wantErrs: []validator.ResponseError{
				{
					FailedField: "Password",
					Tag:         "min",
					Param:       "8",
					Value:       "pa$$",
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "very big password",
			mockArgUser: user.User{
				Username: "JackVorobey",
				Password: "password_that_have_more_than_72_symbols,password_that_have_more_than_72_symbols,password_that_have_more_than_72_symbols,password_that_have_more_than_72_symbols",
			},
			wantErrs: []validator.ResponseError{
				{
					FailedField: "Password",
					Tag:         "max",
					Param:       "72",
					Value:       "password_that_have_more_than_72_symbols,password_that_have_more_than_72_symbols,password_that_have_more_than_72_symbols,password_that_have_more_than_72_symbols",
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "omit digit",
			mockArgUser: user.User{
				Username: "JackVorobey",
				Password: "StrongPasswordButWithoutDigit",
			},
			wantErrs: []validator.ResponseError{
				{
					FailedField: "Password",
					Tag:         "contains_digit",
					Value:       "StrongPasswordButWithoutDigit",
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "omit uppercase",
			mockArgUser: user.User{
				Username: "JackVorobey",
				Password: "strong_but_omit_uppercase123",
			},
			wantErrs: []validator.ResponseError{
				{
					FailedField: "Password",
					Tag:         "contains_uppercase",
					Value:       "strong_but_omit_uppercase123",
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty username",
			mockArgUser: user.User{
				Username: "",
				Password: "VeryStrongPassword123",
			},
			wantErrs: []validator.ResponseError{
				{
					FailedField: "Username",
					Tag:         "required",
					Value:       "",
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "very big username",
			mockArgUser: user.User{
				Username: "JackVorobeyJackVorobeyJackVorobeyJackVorobeyJackVorobeyJackVorobey",
				Password: "VeryStrongPassword123",
			},
			wantErrs: []validator.ResponseError{
				{
					FailedField: "Username",
					Tag:         "max",
					Param:       "50",
					Value:       "JackVorobeyJackVorobeyJackVorobeyJackVorobeyJackVorobeyJackVorobey",
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "username with specials",
			mockArgUser: user.User{
				Username: "Jack_Vorobey",
				Password: "VeryStrongPassword123",
			},
			wantErrs: []validator.ResponseError{
				{
					FailedField: "Username",
					Tag:         "no_specials",
					Value:       "Jack_Vorobey",
				},
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "not unique username",
			mockArgUser: user.User{
				Username: "JackVorobeyNotUnique",
				Password: "VeryStrongPassword123",
			},
			wantErrs: []validator.ResponseError{
				{
					FailedField: "Username",
					Tag:         "unique",
					Value:       "JackVorobeyNotUnique",
				},
			},
			mockRetError: core.ErrNotUniqueUsername,
			wantCode:     http.StatusUnprocessableEntity,
		},
		{
			name: "internal server error",
			mockArgUser: user.User{
				Username: "JackVorobey",
				Password: "VeryStrongPassword123",
			},
			mockRetError: errors.New("internal server error"),
			wantCode:     http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErrs == nil || errors.Is(tt.mockRetError, core.ErrNotUniqueUsername) {
				dependencies.authService.On("SignupBasic", mock.Anything, tt.mockArgUser.ToCoreUser()).Return(tt.mockRetError).Once()
			}

			// request body (multipart)
			reqBody := new(bytes.Buffer)
			mp := multipart.NewWriter(reqBody)
			_ = mp.WriteField("username", tt.mockArgUser.Username)
			_ = mp.WriteField("password", tt.mockArgUser.Password)
			if tt.mockArgUser.Firstname != nil {
				_ = mp.WriteField("firstname", *tt.mockArgUser.Firstname)
			}
			if tt.mockArgUser.Lastname != nil {
				_ = mp.WriteField("lastname", *tt.mockArgUser.Lastname)
			}
			if tt.mockArgUser.Description != nil {
				_ = mp.WriteField("description", *tt.mockArgUser.Description)
			}

			err := mp.Close()
			if err != nil {
				t.Fatal(err.Error())
			}

			// request
			req := httptest.NewRequest(http.MethodPost, route, reqBody)
			req.Header["Content-Type"] = []string{mp.FormDataContentType()}

			// doing request
			resp, _ := app.Test(req, -1)
			assert.Equal(t, tt.wantCode, resp.StatusCode)

			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

			var testValidationErrorResponse TestValidationErrorResponse
			_ = json.Unmarshal(body, &testValidationErrorResponse)
			assert.Equal(t, tt.wantErrs, testValidationErrorResponse.Data.ValidationErrors)
		})
	}
}

func TestRefresh(t *testing.T) {
	app, dependencies := newTestApp(t)

	const route = "/api/v1/auth/token/refresh"

	var (
		userID1 = 1
		userID2 = 2
	)

	tests := []struct {
		name               string
		cookieRefreshToken *string
		mockUserID         *int
		mockRetAccessToken string
		mockRetError       error
		wantCode           int
	}{
		{
			name:               "success",
			cookieRefreshToken: stringPtr("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MX0.tjVEMiS5O2yNzclwLdaZ-FuzrhyqOT7UwM9Hfc0ZQ8Q"),
			mockUserID:         &userID1,
			mockRetAccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEyMzQ1Njc4OTAifQ.Tn0GaHKhAZZB9Y3fZwP4QDG-yvjUXMx3dzAbLKjCX9M",
			wantCode:           http.StatusOK,
		},
		{
			name:               "invalid refresh token",
			cookieRefreshToken: stringPtr("invalid token"),
			wantCode:           http.StatusUnauthorized,
		},
		{
			name:               "internal server error",
			cookieRefreshToken: stringPtr("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Mn0.-ScBrpAXat0bA0Q-kJnL7xnst1-dd_SsIzseTUPT2wE"),
			mockUserID:         &userID2,
			mockRetError:       errors.New("internal server error"),
			wantCode:           http.StatusInternalServerError,
		},
		{
			name:               "invalid refresh token (invalid id but correct sign)",
			cookieRefreshToken: stringPtr("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjEwMH0.sgK75UAPVpqB8iQG3wNw2zlevle3OiOkpqWJLcHAllA"),
			wantCode:           http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockUserID != nil {
				at := tt.mockRetAccessToken
				dependencies.authService.On("Refresh", mock.Anything, *tt.mockUserID).Return(&at, tt.mockRetError).Once()
			}

			req := httptest.NewRequest(http.MethodPost, route, http.NoBody)
			if tt.cookieRefreshToken != nil {
				req.AddCookie(&http.Cookie{
					Name:  "refresh_token",
					Value: *tt.cookieRefreshToken,
				})
			}

			resp, _ := app.Test(req, -1)
			body, _ := io.ReadAll(resp.Body)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantCode, resp.StatusCode)

			var testAccessTokenSuccessResponse TestAccessTokenSuccessResponse
			_ = json.Unmarshal(body, &testAccessTokenSuccessResponse)
			if tt.wantCode == http.StatusOK {
				assert.Equal(t, tt.mockRetAccessToken, testAccessTokenSuccessResponse.Data)
			}
		})
	}
}
