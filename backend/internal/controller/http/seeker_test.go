package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/seeker"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

const (
	validUserID       = 1
	nonExistentUserID = 999
)

var mockSeeker = core.Seeker{
	UserID:           validUserID,
	AnimalType:       "cat",
	Description:      "Test description",
	Location:         "Moscow",
	EquipmentRental:  500,
	HaveMetalCage:    true,
	HavePlasticCage:  true,
	HaveNet:          true,
	HaveLadder:       true,
	HaveOther:        " ",
	Price:            100,
	HaveCar:          true,
	WillingnessCarry: "yes",
}

func TestHttp_GetSeeker(t *testing.T) {
	t.Parallel()
	app, dependencies := newTestApp(t)

	const route = "/api/v1/seekers/"

	tests := []struct {
		name         string
		userID       int
		request      core.Seeker
		requestError error
		wantCode     int
	}{
		{
			name:         "success",
			userID:       validUserID,
			request:      mockSeeker,
			requestError: nil,
			wantCode:     http.StatusOK,
		},
		{
			name:         "not found",
			userID:       nonExistentUserID,
			request:      core.Seeker{},
			requestError: core.ErrSeekerNotFound,
			wantCode:     http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dependencies.seekerService.On(
				"GetSeeker",
				mock.Anything,
				mock.Anything,
			).Return(mockSeeker, tt.requestError).Once()

			req := httptest.NewRequest(http.MethodGet, route+strconv.Itoa(tt.userID), nil)

			resp, err := app.Test(req)
			require.NoError(t, err, "Request failed")
			defer resp.Body.Close()

			assert.Equal(t, tt.wantCode, resp.StatusCode)
		})
	}
}

var mockCreateSeeker = seeker.CreateSeeker{
	AnimalType:       "cat",
	Description:      "Test description",
	Location:         "Moscow",
	EquipmentRental:  500,
	HaveMetalCage:    true,
	HavePlasticCage:  true,
	HaveNet:          true,
	HaveLadder:       true,
	HaveOther:        " ",
	Price:            100,
	HaveCar:          true,
	WillingnessCarry: "yes",
}

func TestHttp_CreateSeeker(t *testing.T) {
	t.Parallel()
	app, dependencies := newTestApp(t)

	const route = "/api/v1/seekers"

	tests := []struct {
		name      string
		request   seeker.CreateSeeker
		token     string
		setupMock func()
		wantCode  int
	}{
		{
			name:    "success",
			request: mockCreateSeeker,
			token:   token,
			setupMock: func() {
				dependencies.seekerService.On(
					"CreateSeeker",
					mock.Anything,
					mock.Anything,
				).Return(mockSeeker, nil).Once()
			},
			wantCode: http.StatusOK,
		},
		{
			name:      "unprocessable entity",
			request:   seeker.CreateSeeker{},
			token:     token,
			setupMock: func() {},
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name: "missing location",
			request: func() seeker.CreateSeeker {
				req := mockCreateSeeker
				req.Location = ""
				return req
			}(),
			token:     token,
			setupMock: func() {},
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name: "missing animal type",
			request: func() seeker.CreateSeeker {
				req := mockCreateSeeker
				req.AnimalType = " "
				return req
			}(),
			token:     token,
			setupMock: func() {},
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name: "missing equipment rental",
			request: func() seeker.CreateSeeker {
				req := mockCreateSeeker
				req.EquipmentRental = -100
				return req
			}(),
			token:     token,
			setupMock: func() {},
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name: "missing willingness carry",
			request: func() seeker.CreateSeeker {
				req := mockCreateSeeker
				req.WillingnessCarry = " "
				return req
			}(),
			token:     token,
			setupMock: func() {},
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name: "negative price",
			request: func() seeker.CreateSeeker {
				req := mockCreateSeeker
				req.Price = -100
				return req
			}(),
			token:     token,
			setupMock: func() {},
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "unauthorized missing token",
			request:   mockCreateSeeker,
			token:     "",
			setupMock: func() {},
			wantCode:  http.StatusUnauthorized,
		},
		{
			name:      "unauthorized invalid token",
			request:   mockCreateSeeker,
			token:     "invalid_token",
			setupMock: func() {},
			wantCode:  http.StatusUnauthorized,
		},
		{
			name: "user not found",
			request: func() seeker.CreateSeeker {
				req := mockCreateSeeker
				return req
			}(),
			token: token,
			setupMock: func() {
				dependencies.seekerService.On(
					"CreateSeeker",
					mock.Anything,
					mock.Anything,
				).Return(core.Seeker{}, core.ErrNoSuchUser).Once()
			},
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			body, err := json.Marshal(tt.request)
			require.NoError(t, err, "Failed to marshal request")

			req := httptest.NewRequest(
				http.MethodPost,
				route,
				bytes.NewReader(body),
			)

			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err, "Request failed")
			defer resp.Body.Close()

			assert.Equal(t, tt.wantCode, resp.StatusCode)
		})
	}
}

var mockUpdateSeeker = seeker.UpdateSeeker{
	AnimalType: stringPtr("dog"),
}

func TestHttp_UpdateSeeker(t *testing.T) {
	t.Parallel()
	app, dependencies := newTestApp(t)

	const route = "/api/v1/seekers"
	mockSeeker.AnimalType = "dog"

	tests := []struct {
		name      string
		request   seeker.UpdateSeeker
		token     string
		setupMock func()
		wantCode  int
	}{
		{
			name:    "success",
			request: mockUpdateSeeker,
			token:   token,
			setupMock: func() {
				dependencies.seekerService.On(
					"UpdateSeeker",
					mock.Anything,
					mock.Anything,
				).Return(mockSeeker, nil).Once()
			},
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			body, err := json.Marshal(tt.request)
			require.NoError(t, err, "Failed to marshal request")

			req := httptest.NewRequest(
				http.MethodPatch,
				route+"/"+strconv.Itoa(validUserID),
				bytes.NewReader(body),
			)

			req.Header.Set("Authorization", "Bearer "+tt.token)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err, "Request failed")
			defer resp.Body.Close()

			assert.Equal(t, tt.wantCode, resp.StatusCode)
		})
	}
}

func TestHttp_DeleteSeeker(t *testing.T) {
	t.Parallel()
	app, dependencies := newTestApp(t)

	const route = "/api/v1/seekers"

	tests := []struct {
		name      string
		seekerID  string
		token     string
		setupMock func()
		wantCode  int
	}{
		{
			name:     "success",
			seekerID: strconv.Itoa(validUserID),
			token:    token,
			setupMock: func() {
				dependencies.seekerService.On(
					"DeleteSeeker",
					mock.Anything,
					mock.Anything,
				).Return(nil).Once()
			},
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("%s/%s", route, tt.seekerID),
				nil,
			)

			req.Header.Set("Authorization", "Bearer "+tt.token)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err, "Request failed")
			defer resp.Body.Close()

			assert.Equal(t, tt.wantCode, resp.StatusCode)

			bodyBytes, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "Read response failed")
			var response model.Response
			json.Unmarshal(bodyBytes, &response)

			assert.Equal(t, "Delete", response.Data)
		})
	}
}

func TestHttp_GetSeekers(t *testing.T) {
	t.Parallel()
	app, dependencies := newTestApp(t)

	const route = "/api/v1/seekers"

	mockSeekers := []core.Seeker{
		{
			ID:         validUserID,
			AnimalType: "dog",
			Location:   "Moscow",
			Price:      1000,
		},
		{
			ID:         validUserID + 1,
			AnimalType: "cat",
			Location:   "St. Petersburg",
			Price:      800,
		},
	}

	tests := []struct {
		name          string
		queryParams   map[string]string
		mockBehaviour func()
		wantCode      int
	}{
		{
			name: "success with params",
			queryParams: map[string]string{
				"limit":       "10",
				"offset":      "0",
				"price_min":   "500",
				"animal_type": "dog",
			},
			mockBehaviour: func() {
				dependencies.seekerService.On(
					"GetAllSeekers",
					mock.Anything,
					mock.Anything,
				).Return(mockSeekers, nil).Once()

			},
			wantCode: http.StatusOK,
		},
		{
			name: "invalid limit",
			queryParams: map[string]string{
				"limit": "-1",
			},
			mockBehaviour: func() {},
			wantCode:      http.StatusUnprocessableEntity,
		},
		{
			name: "invalid offset",
			queryParams: map[string]string{
				"offset": "-1",
			},
			mockBehaviour: func() {},
			wantCode:      http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()

			newRoute := route

			params := url.Values{}
			for k, v := range tt.queryParams {
				params.Add(k, v)
			}
			if len(params) > 0 {
				newRoute += "?" + params.Encode()
			}

			req := httptest.NewRequest(http.MethodGet, newRoute, http.NoBody)
			t.Logf("Request URL: %s", req.URL)
			resp, err := app.Test(req)
			require.NoError(t, err, "Request failed")
			defer resp.Body.Close()

			assert.Equal(t, tt.wantCode, resp.StatusCode)

		})
	}
}
