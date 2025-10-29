package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/report"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateReport(t *testing.T) {
	t.Parallel()
	app, dependencies := newTestApp(t)

	const route = "/api/v1/reports"

	tests := []struct {
		name          string
		token         string
		requestBody   interface{}
		mockBehaviour func()
		wantCode      int
	}{
		{
			name:  "success post report",
			token: token,
			requestBody: report.CreateRequestBodyReport{
				TargetID:   1,
				TargetType: "post",
				Reason:     "spam",
			},
			mockBehaviour: func() {
				dependencies.reportService.EXPECT().
					CreateReport(mock.Anything, mock.MatchedBy(func(r core.Report) bool {
						return r.ReportableID == 1 && r.ReportableType == core.ReportableTypePost
					})).
					Return(nil).Once()
			},
			wantCode: http.StatusCreated,
		},
		{
			name:  "success comment report",
			token: token,
			requestBody: report.CreateRequestBodyReport{
				TargetID:   2,
				TargetType: "comment",
				Reason:     "violent_content",
			},
			mockBehaviour: func() {
				dependencies.reportService.EXPECT().
					CreateReport(mock.Anything, mock.MatchedBy(func(r core.Report) bool {
						return r.ReportableID == 2 && r.ReportableType == core.ReportableTypeComment
					})).
					Return(nil).Once()
			},
			wantCode: http.StatusCreated,
		},
		{
			name:  "unauthorized",
			token: "",
			requestBody: report.CreateRequestBodyReport{
				TargetID:   1,
				TargetType: "post",
				Reason:     "spam",
			},
			mockBehaviour: func() {},
			wantCode:      http.StatusUnauthorized,
		},
		{
			name:  "target not found",
			token: token,
			requestBody: report.CreateRequestBodyReport{
				TargetID:   1,
				TargetType: "post",
				Reason:     "spam",
			},
			mockBehaviour: func() {
				dependencies.reportService.EXPECT().
					CreateReport(mock.Anything, mock.Anything).
					Return(core.ErrTargetNotFound).Once()
			},
			wantCode: http.StatusNotFound,
		},
		{
			name:  "duplicate report",
			token: token,
			requestBody: report.CreateRequestBodyReport{
				TargetID:   1,
				TargetType: "post",
				Reason:     "spam",
			},
			mockBehaviour: func() {
				dependencies.reportService.EXPECT().
					CreateReport(mock.Anything, mock.Anything).
					Return(core.ErrDuplicateReport).Once()
			},
			wantCode: http.StatusConflict,
		},
		{
			name:  "invalid reportable type - validation error",
			token: token,
			requestBody: report.CreateRequestBodyReport{
				TargetID:   1,
				TargetType: "invalid_type",
				Reason:     "spam",
			},
			mockBehaviour: func() {},
			wantCode:      http.StatusUnprocessableEntity,
		},
		{
			name:  "validation error - missing reason",
			token: token,
			requestBody: report.CreateRequestBodyReport{
				TargetID:   1,
				TargetType: "post",
				Reason:     "",
			},
			mockBehaviour: func() {},
			wantCode:      http.StatusUnprocessableEntity,
		},
		{
			name:  "validation error - missing target_id",
			token: token,
			requestBody: map[string]interface{}{
				"target_type": "post",
				"reason":      "spam",
			},
			mockBehaviour: func() {},
			wantCode:      http.StatusUnprocessableEntity,
		},
		{
			name:  "validation error - missing target_type",
			token: token,
			requestBody: map[string]interface{}{
				"target_id": 1,
				"reason":    "spam",
			},
			mockBehaviour: func() {},
			wantCode:      http.StatusUnprocessableEntity,
		},
		{
			name:  "validation error - invalid reason",
			token: token,
			requestBody: report.CreateRequestBodyReport{
				TargetID:   1,
				TargetType: "post",
				Reason:     "invalid_reason",
			},
			mockBehaviour: func() {},
			wantCode:      http.StatusUnprocessableEntity,
		},
		{
			name:  "internal server error",
			token: token,
			requestBody: report.CreateRequestBodyReport{
				TargetID:   1,
				TargetType: "post",
				Reason:     "spam",
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
			tt.mockBehaviour()

			body, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(
				http.MethodPost,
				route,
				bytes.NewReader(body),
			)
			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantCode, resp.StatusCode)

			if tt.wantCode != http.StatusCreated {
				respBody, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				assert.NotEmpty(t, respBody)
			}
		})
	}
}
