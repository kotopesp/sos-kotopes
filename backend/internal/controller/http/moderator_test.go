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

	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/moderator"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetReportedPosts(t *testing.T) {
	t.Parallel()
	app, dependencies := newTestApp(t)

	const baseRoute = "/api/v1/moderation/posts"

	tests := []struct {
		name          string
		token         string
		queryParams   string
		mockBehaviour func()
		wantCode      int
		wantResponse  bool
	}{
		{
			name:        "success with ASC filter",
			token:       token,
			queryParams: "filter=ASC",
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()

				dependencies.moderatorService.EXPECT().
					GetPostsForModeration(mock.Anything, core.Filter("ASC")).
					Return([]core.PostForModeration{
						{Post: core.Post{ID: 1}, Reasons: []string{"spam"}},
					}, nil).Once()

				dependencies.postService.EXPECT().
					BuildPostDetailsList(mock.Anything, mock.Anything, mock.Anything).
					Return([]core.PostDetails{{Post: core.Post{ID: 1}}}, nil).Once()
			},
			wantCode:     http.StatusOK,
			wantResponse: true,
		},
		{
			name:        "success with DESC filter",
			token:       token,
			queryParams: "filter=DESC",
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()

				dependencies.moderatorService.EXPECT().
					GetPostsForModeration(mock.Anything, core.Filter("DESC")).
					Return([]core.PostForModeration{}, nil).Once()

				dependencies.postService.EXPECT().
					BuildPostDetailsList(mock.Anything, mock.Anything, mock.Anything).
					Return([]core.PostDetails{}, nil).Once()
			},
			wantCode:     http.StatusOK,
			wantResponse: true,
		},
		{
			name:          "unauthorized - missing token",
			token:         "",
			queryParams:   "filter=ASC",
			mockBehaviour: func() {},
			wantCode:      http.StatusUnauthorized,
		},
		{
			name:        "forbidden - not a moderator",
			token:       token,
			queryParams: "filter=ASC",
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, core.ErrNoSuchModerator).Once()
			},
			wantCode: http.StatusForbidden,
		},
		{
			name:        "validation error - missing filter",
			token:       token,
			queryParams: "",
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:        "validation error - invalid filter value",
			token:       token,
			queryParams: "filter=INVALID",
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:        "bad request - malformed query params",
			token:       token,
			queryParams: "%%%",
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:        "no content - no posts waiting",
			token:       token,
			queryParams: "filter=ASC",
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()

				dependencies.moderatorService.EXPECT().
					GetPostsForModeration(mock.Anything, core.Filter("ASC")).
					Return(nil, core.ErrNoPostsWaitingForModeration).Once()
			},
			wantCode: http.StatusNoContent,
		},
		{
			name:        "internal error - get posts failed",
			token:       token,
			queryParams: "filter=ASC",
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()

				dependencies.moderatorService.EXPECT().
					GetPostsForModeration(mock.Anything, core.Filter("ASC")).
					Return(nil, errors.New("db error")).Once()
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:        "internal error - build details failed",
			token:       token,
			queryParams: "filter=ASC",
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()

				dependencies.moderatorService.EXPECT().
					GetPostsForModeration(mock.Anything, core.Filter("ASC")).
					Return([]core.PostForModeration{{Post: core.Post{ID: 1}}}, nil).Once()

				dependencies.postService.EXPECT().
					BuildPostDetailsList(mock.Anything, mock.Anything, mock.Anything).
					Return(nil, errors.New("build error")).Once()
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()

			url := baseRoute
			if tt.queryParams != "" {
				url = fmt.Sprintf("%s?%s", baseRoute, tt.queryParams)
			}

			req := httptest.NewRequest(http.MethodGet, url, http.NoBody)
			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			}

			resp, err := app.Test(req)
			require.NoError(t, err)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			_ = resp.Body.Close()

			assert.Equal(t, tt.wantCode, resp.StatusCode)

			if tt.wantCode == http.StatusOK {
				var data model.Response
				err = json.Unmarshal(body, &data)
				require.NoError(t, err)

			}
		})
	}
}
func TestDeletePostByModerator(t *testing.T) {
	t.Parallel()
	app, dependencies := newTestApp(t)

	const route = "/api/v1/moderation/posts/%d"

	tests := []struct {
		name          string
		token         string
		postID        int
		mockBehaviour func()
		wantCode      int
	}{
		{
			name:   "success",
			token:  token,
			postID: 1,
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()

				dependencies.moderatorService.EXPECT().
					DeletePost(mock.Anything, mock.Anything).
					Return(nil).Once()
			},
			wantCode: http.StatusOK,
		},
		{
			name:   "unauthorized",
			token:  "",
			postID: 1,
			mockBehaviour: func() {
			},
			wantCode: http.StatusUnauthorized,
		},
		{
			name:   "forbidden",
			token:  token,
			postID: 1,
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, core.ErrNoSuchModerator).Once()
			},
			wantCode: http.StatusForbidden,
		},
		{
			name:   "validation error",
			token:  token,
			postID: -1,
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:   "post not found",
			token:  token,
			postID: 1,
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()

				dependencies.moderatorService.EXPECT().
					DeletePost(mock.Anything, mock.Anything).
					Return(core.ErrPostNotFound).Once()
			},
			wantCode: http.StatusNotFound,
		},
		{
			name:   "internal server error",
			token:  token,
			postID: 1,
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()

				dependencies.moderatorService.EXPECT().
					DeletePost(mock.Anything, mock.Anything).
					Return(errors.New("internal server error")).Once()
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf(route, tt.postID), http.NoBody)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))

			resp, err := app.Test(req, -1)
			require.NoError(t, err)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			err = resp.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.wantCode, resp.StatusCode)

			if tt.wantCode != http.StatusOK {
				var data model.Response
				err = json.Unmarshal(body, &data)
				require.NoError(t, err)
			}
		})
	}
}

func TestApprovePostByModerator(t *testing.T) {
	t.Parallel()
	app, dependencies := newTestApp(t)

	const route = "/api/v1/moderation/posts/%d"

	tests := []struct {
		name          string
		token         string
		postID        int
		mockBehaviour func()
		wantCode      int
	}{
		{
			name:   "success",
			token:  token,
			postID: 1,
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()

				dependencies.moderatorService.EXPECT().
					ApprovePost(mock.Anything, mock.Anything).
					Return(nil).Once()
			},
			wantCode: http.StatusOK,
		},
		{
			name:   "unauthorized",
			token:  "",
			postID: 1,
			mockBehaviour: func() {
			},
			wantCode: http.StatusUnauthorized,
		},
		{
			name:   "forbidden",
			token:  token,
			postID: 1,
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, core.ErrNoSuchModerator).Once()
			},
			wantCode: http.StatusForbidden,
		},
		{
			name:   "validation error",
			token:  token,
			postID: -1,
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:   "approve post failed",
			token:  token,
			postID: 1,
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()

				dependencies.moderatorService.EXPECT().
					ApprovePost(mock.Anything, mock.Anything).
					Return(errors.New("internal error")).Once()
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:   "post not found",
			token:  token,
			postID: 1,
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()

				dependencies.moderatorService.EXPECT().
					ApprovePost(mock.Anything, mock.Anything).
					Return(core.ErrPostNotFound).Once()
			},
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()

			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf(route, tt.postID), http.NoBody)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))

			resp, err := app.Test(req, -1)
			require.NoError(t, err)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			err = resp.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.wantCode, resp.StatusCode)

			if tt.wantCode != http.StatusOK {
				var data model.Response
				err = json.Unmarshal(body, &data)
				require.NoError(t, err)
			}
		})
	}
}

func TestGetReportedComments(t *testing.T) {
	t.Parallel()
	app, dependencies := newTestApp(t)

	const baseRoute = "/api/v1/moderation/comments"

	tests := []struct {
		name          string
		token         string
		queryParams   string
		mockBehaviour func()
		wantCode      int
		wantResponse  bool
	}{
		{
			name:        "success with ASC filter",
			token:       token,
			queryParams: "filter=ASC",
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()

				dependencies.moderatorService.EXPECT().
					GetCommentsForModeration(mock.Anything, core.Filter("ASC")).
					Return([]core.CommentForModeration{
						{Comment: core.Comment{ID: 1, Content: "test"}, Reasons: []string{"spam"}},
					}, nil).Once()
			},
			wantCode:     http.StatusOK,
			wantResponse: true,
		},
		{
			name:        "success with DESC filter",
			token:       token,
			queryParams: "filter=DESC",
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()

				dependencies.moderatorService.EXPECT().
					GetCommentsForModeration(mock.Anything, core.Filter("DESC")).
					Return([]core.CommentForModeration{}, nil).Once()
			},
			wantCode:     http.StatusOK,
			wantResponse: true,
		},
		{
			name:          "unauthorized - missing token",
			token:         "",
			queryParams:   "filter=ASC",
			mockBehaviour: func() {},
			wantCode:      http.StatusUnauthorized,
		},
		{
			name:        "forbidden - not a moderator",
			token:       token,
			queryParams: "filter=ASC",
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, core.ErrNoSuchModerator).Once()
			},
			wantCode: http.StatusForbidden,
		},
		{
			name:        "validation error - missing filter",
			token:       token,
			queryParams: "",
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:        "validation error - invalid filter value",
			token:       token,
			queryParams: "filter=INVALID",
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:        "bad request - malformed query params",
			token:       token,
			queryParams: "%%%",
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:        "no content - no comments waiting",
			token:       token,
			queryParams: "filter=ASC",
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()

				dependencies.moderatorService.EXPECT().
					GetCommentsForModeration(mock.Anything, core.Filter("ASC")).
					Return(nil, core.ErrNoCommentsWaitingForModeration).Once()
			},
			wantCode: http.StatusNoContent,
		},
		{
			name:        "internal error - get comments failed",
			token:       token,
			queryParams: "filter=ASC",
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()

				dependencies.moderatorService.EXPECT().
					GetCommentsForModeration(mock.Anything, core.Filter("ASC")).
					Return(nil, errors.New("db error")).Once()
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()

			url := baseRoute
			if tt.queryParams != "" {
				url = fmt.Sprintf("%s?%s", baseRoute, tt.queryParams)
			}

			req := httptest.NewRequest(http.MethodGet, url, http.NoBody)
			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			}

			resp, err := app.Test(req)
			require.NoError(t, err)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			_ = resp.Body.Close()

			assert.Equal(t, tt.wantCode, resp.StatusCode)

			if tt.wantCode == http.StatusOK {
				var data model.Response
				err = json.Unmarshal(body, &data)
				require.NoError(t, err)
			}
		})
	}
}

func TestDeleteCommentByModerator(t *testing.T) {
	t.Parallel()
	app, dependencies := newTestApp(t)

	const route = "/api/v1/moderation/comments/%d"

	tests := []struct {
		name          string
		token         string
		commentID     int
		mockBehaviour func()
		wantCode      int
	}{
		{
			name:      "success",
			token:     token,
			commentID: 1,
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()

				dependencies.moderatorService.EXPECT().
					DeleteComment(mock.Anything, 1).
					Return(nil).Once()
			},
			wantCode: http.StatusOK,
		},
		{
			name:          "unauthorized",
			token:         "",
			commentID:     1,
			mockBehaviour: func() {},
			wantCode:      http.StatusUnauthorized,
		},
		{
			name:      "forbidden",
			token:     token,
			commentID: 1,
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, core.ErrNoSuchModerator).Once()
			},
			wantCode: http.StatusForbidden,
		},
		{
			name:      "validation error",
			token:     token,
			commentID: -1,
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:      "comment not found",
			token:     token,
			commentID: 1,
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()

				dependencies.moderatorService.EXPECT().
					DeleteComment(mock.Anything, 1).
					Return(core.ErrNoSuchComment).Once()
			},
			wantCode: http.StatusNotFound,
		},
		{
			name:      "internal server error",
			token:     token,
			commentID: 1,
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()

				dependencies.moderatorService.EXPECT().
					DeleteComment(mock.Anything, 1).
					Return(errors.New("internal server error")).Once()
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf(route, tt.commentID), http.NoBody)
			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			}

			resp, err := app.Test(req, -1)
			require.NoError(t, err)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			err = resp.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.wantCode, resp.StatusCode)

			if tt.wantCode != http.StatusOK {
				var data model.Response
				err = json.Unmarshal(body, &data)
				require.NoError(t, err)
			}
		})
	}
}

func TestApproveCommentByModerator(t *testing.T) {
	t.Parallel()
	app, dependencies := newTestApp(t)

	const route = "/api/v1/moderation/comments/%d"

	tests := []struct {
		name          string
		token         string
		commentID     int
		mockBehaviour func()
		wantCode      int
	}{
		{
			name:      "success",
			token:     token,
			commentID: 1,
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()

				dependencies.moderatorService.EXPECT().
					ApproveComment(mock.Anything, 1).
					Return(nil).Once()
			},
			wantCode: http.StatusOK,
		},
		{
			name:          "unauthorized",
			token:         "",
			commentID:     1,
			mockBehaviour: func() {},
			wantCode:      http.StatusUnauthorized,
		},
		{
			name:      "forbidden",
			token:     token,
			commentID: 1,
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, core.ErrNoSuchModerator).Once()
			},
			wantCode: http.StatusForbidden,
		},
		{
			name:      "validation error",
			token:     token,
			commentID: -1,
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:      "approve comment failed",
			token:     token,
			commentID: 1,
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()

				dependencies.moderatorService.EXPECT().
					ApproveComment(mock.Anything, 1).
					Return(errors.New("internal error")).Once()
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:      "comment not found",
			token:     token,
			commentID: 1,
			mockBehaviour: func() {
				dependencies.moderatorService.EXPECT().
					GetModerator(mock.Anything, mock.Anything).
					Return(core.Moderator{}, nil).Once()

				dependencies.moderatorService.EXPECT().
					ApproveComment(mock.Anything, 1).
					Return(core.ErrNoSuchComment).Once()
			},
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()

			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf(route, tt.commentID), http.NoBody)
			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			}

			resp, err := app.Test(req, -1)
			require.NoError(t, err)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			err = resp.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.wantCode, resp.StatusCode)

			if tt.wantCode != http.StatusOK {
				var data model.Response
				err = json.Unmarshal(body, &data)
				require.NoError(t, err)
			}
		})
	}
}

func TestBanUser(t *testing.T) {
	t.Parallel()
	app, dependencies := newTestApp(t)

	const route = "/api/v1/moderation/users/ban"

	tests := []struct {
		name          string
		token         string
		requestBody   interface{}
		mockBehaviour func()
		wantCode      int
		wantError     string
		checkResponse func(t *testing.T, body []byte)
	}{
		{
			name:  "success - ban user with report",
			token: token,
			requestBody: moderator.BanUserRequest{
				UserID:   1,
				ReportID: func() *int { i := 5; return &i }(),
			},
			mockBehaviour: func() {
				dependencies.moderatorService.On("GetModerator", mock.Anything, 1).Return(core.Moderator{}, nil).Once()
				dependencies.moderatorService.On("BanUser", mock.Anything, mock.MatchedBy(func(record core.BannedUserRecord) bool {
					return record.UserID == 1 && record.ModeratorID == 1 && *record.ReportID == 5
				})).Return(nil).Once()
			},
			wantCode: fiber.StatusOK,
		},
		{
			name:  "success - ban user without report",
			token: token,
			requestBody: moderator.BanUserRequest{
				UserID:   1,
				ReportID: nil,
			},
			mockBehaviour: func() {
				dependencies.moderatorService.On("GetModerator", mock.Anything, 1).Return(core.Moderator{}, nil).Once()
				dependencies.moderatorService.On("BanUser", mock.Anything, mock.MatchedBy(func(record core.BannedUserRecord) bool {
					return record.UserID == 1 && record.ModeratorID == 1 && record.ReportID == nil
				})).Return(nil).Once()
			},
			wantCode: fiber.StatusOK,
		},
		{
			name:        "unauthorized - missing token",
			token:       "",
			requestBody: moderator.BanUserRequest{UserID: 1},
			mockBehaviour: func() {
			},
			wantCode:  fiber.StatusUnauthorized,
			wantError: "Unauthorized",
		},
		{
			name:        "forbidden - not a moderator",
			token:       token,
			requestBody: moderator.BanUserRequest{UserID: 1},
			mockBehaviour: func() {
				dependencies.moderatorService.On("GetModerator", mock.Anything, 1).Return(core.Moderator{}, core.ErrNoSuchModerator).Once()
			},
			wantCode:  fiber.StatusForbidden,
			wantError: core.ErrNoSuchModerator.Error(),
		},
		{
			name:        "validation error - invalid user_id",
			token:       token,
			requestBody: map[string]interface{}{"user_id": 0, "report_id": nil},
			mockBehaviour: func() {
				dependencies.moderatorService.On("GetModerator", mock.Anything, 1).Return(core.Moderator{}, nil).Once()
			},
			wantCode:  fiber.StatusUnprocessableEntity,
			wantError: "Invalid request body",
		},
		{
			name:        "validation error - negative report_id",
			token:       token,
			requestBody: map[string]interface{}{"user_id": 1, "report_id": -1},
			mockBehaviour: func() {
				dependencies.moderatorService.On("GetModerator", mock.Anything, 1).Return(core.Moderator{}, nil).Once()
			},
			wantCode:  fiber.StatusUnprocessableEntity,
			wantError: "Invalid request body",
		},
		{
			name:        "user not found",
			token:       token,
			requestBody: moderator.BanUserRequest{UserID: 999},
			mockBehaviour: func() {
				dependencies.moderatorService.On("GetModerator", mock.Anything, 1).Return(core.Moderator{}, nil).Once()
				dependencies.moderatorService.On("BanUser", mock.Anything, mock.Anything).Return(core.ErrNoSuchUser).Once()
			},
			wantCode:  fiber.StatusNotFound,
			wantError: core.ErrNoSuchUser.Error(),
		},
		{
			name:        "user already banned",
			token:       token,
			requestBody: moderator.BanUserRequest{UserID: 1},
			mockBehaviour: func() {
				dependencies.moderatorService.On("GetModerator", mock.Anything, 1).Return(core.Moderator{}, nil).Once()
				dependencies.moderatorService.On("BanUser", mock.Anything, mock.Anything).Return(core.ErrUserAlreadyBanned).Once()
			},
			wantCode:  fiber.StatusConflict,
			wantError: core.ErrUserAlreadyBanned.Error(),
		},
		{
			name:        "internal server error",
			token:       token,
			requestBody: moderator.BanUserRequest{UserID: 1},
			mockBehaviour: func() {
				dependencies.moderatorService.On("GetModerator", mock.Anything, 1).Return(core.Moderator{}, nil).Once()
				dependencies.moderatorService.On("BanUser", mock.Anything, mock.Anything).Return(errors.New("database error")).Once()
			},
			wantCode:  fiber.StatusInternalServerError,
			wantError: "database error",
		},
		{
			name:        "malformed JSON",
			token:       token,
			requestBody: "{invalid json",
			mockBehaviour: func() {
				dependencies.moderatorService.On("GetModerator", mock.Anything, 1).Return(core.Moderator{}, nil).Once()
			},
			wantCode:  fiber.StatusUnprocessableEntity,
			wantError: "Invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()

			var requestBody []byte
			switch body := tt.requestBody.(type) {
			case string:
				requestBody = []byte(body)
			default:
				requestBody, _ = json.Marshal(body)
			}

			req := httptest.NewRequest(http.MethodPost, route, bytes.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.wantCode, resp.StatusCode)

			if tt.wantError != "" {
				var errorResp model.Response
				err = json.Unmarshal(body, &errorResp)
				require.NoError(t, err)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, body)
			}

			dependencies.moderatorService.AssertExpectations(t)
		})
	}
}
