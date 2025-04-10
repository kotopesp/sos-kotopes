package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/report"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateReport(t *testing.T) {
	t.Parallel()
	app, dependencies := newTestApp(t)

	const route = "/api/v1/reports/%d"

	tests := []struct {
		name          string
		token         string
		postID        int
		requestBody   interface{}
		mockBehaviour func()
		wantCode      int
	}{
		{
			name:   "success",
			token:  token,
			postID: 1,
			requestBody: report.CreateRequestBodyReport{
				Reason: "spam",
			},
			mockBehaviour: func() {
				dependencies.reportService.EXPECT().
					CreateReport(mock.Anything, mock.Anything).
					Return(nil).Once()
			},
			wantCode: http.StatusCreated,
		},
		{
			name:   "unauthorized",
			token:  "",
			postID: 1,
			requestBody: report.CreateRequestBodyReport{
				Reason: "spam",
			},
			mockBehaviour: func() {},
			wantCode:      http.StatusUnauthorized,
		},
		{
			name:   "post not found",
			token:  token,
			postID: 1,
			requestBody: report.CreateRequestBodyReport{
				Reason: "spam",
			},
			mockBehaviour: func() {
				dependencies.reportService.EXPECT().
					CreateReport(mock.Anything, mock.Anything).
					Return(core.ErrPostNotFound).Once()
			},
			wantCode: http.StatusNotFound,
		},
		{
			name:   "duplicate report",
			token:  token,
			postID: 1,
			requestBody: report.CreateRequestBodyReport{
				Reason: "spam",
			},
			mockBehaviour: func() {
				dependencies.reportService.EXPECT().
					CreateReport(mock.Anything, mock.Anything).
					Return(core.ErrDuplicateReport).Once()
			},
			wantCode: http.StatusConflict,
		},
		{
			name:   "validation error - missing reason",
			token:  token,
			postID: 1,
			requestBody: report.CreateRequestBodyReport{
				Reason: "",
			},
			mockBehaviour: func() {},
			wantCode:      http.StatusUnprocessableEntity,
		},
		{
			name:   "validation error - invalid reason",
			token:  token,
			postID: 1,
			requestBody: report.CreateRequestBodyReport{
				Reason: "invalid_reason",
			},
			mockBehaviour: func() {},
			wantCode:      http.StatusUnprocessableEntity,
		},
		{
			name:   "internal server error",
			token:  token,
			postID: 1,
			requestBody: report.CreateRequestBodyReport{
				Reason: "spam",
			},
			mockBehaviour: func() {
				dependencies.reportService.EXPECT().
					CreateReport(mock.Anything, mock.Anything).
					Return(errors.New("internal error")).Once()
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Мокируем поведение в зависимости от теста
			tt.mockBehaviour()

			// Подготавливаем тело запроса
			body, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			// Исправлено создание запроса
			req := httptest.NewRequest(
				http.MethodPost,
				fmt.Sprintf(route, tt.postID),
				bytes.NewReader(body),
			)
			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			}
			req.Header.Set("Content-Type", "application/json")

			// Выполняем запрос
			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Проверяем статус ответа
			assert.Equal(t, tt.wantCode, resp.StatusCode)

			// Для ошибок можно дополнительно проверить тело ответа
			if tt.wantCode != http.StatusCreated {
				respBody, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				assert.NotEmpty(t, respBody)
			}
		})
	}
}
