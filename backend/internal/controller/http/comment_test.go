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

	"github.com/brianvoe/gofakeit/v7"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/comment"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/pagination"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	amountOfComments = 10
	token            = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MX0.tjVEMiS5O2yNzclwLdaZ-FuzrhyqOT7UwM9Hfc0ZQ8Q" // nolint
	authorID         = 1
)

func generateComments(t *testing.T) []core.Comment {
	t.Helper()

	comments := make([]core.Comment, 0, amountOfComments)
	for range amountOfComments {
		var comm core.Comment
		err := gofakeit.Struct(&comm)
		require.NoError(t, err)

		comments = append(comments, comm)
	}

	return comments
}

func TestGetComments(t *testing.T) {
	t.Parallel()
	app, dependencies := newTestApp(t)

	const route = "/api/v1/posts/%d/comments?limit=%d&offset=%d"

	var (
		limit         = amountOfComments
		invalidLimit  = 0
		offset        = 1
		comments      = generateComments(t)
		totalComments = 200
	)

	type Data struct {
		Comments []comment.Comment     `json:"comments"`
		Meta     pagination.Pagination `json:"meta"`
	}

	type GetCommentsResponse struct {
		Data Data `json:"data"`
	}

	tests := []struct {
		name          string
		limit         int
		offset        int
		postID        int
		mockBehaviour func()
		wantData      GetCommentsResponse
		wantCode      int
	}{
		{
			name:   "success",
			postID: 1,
			limit:  limit,
			offset: offset,
			mockBehaviour: func() {
				dependencies.commentService.EXPECT().
					GetAllComments(mock.Anything, core.GetAllCommentsParams{
						PostID: 1,
						Limit:  &limit,
						Offset: &offset,
					}).Return(comments, totalComments, nil).Once()
			},
			wantData: GetCommentsResponse{
				Data: Data{
					Comments: comment.ToModelCommentsSlice(comments),
					Meta: pagination.Pagination{
						Total:       totalComments,
						TotalPages:  (totalComments + limit - 1) / limit,
						CurrentPage: offset/limit + 1,
						PerPage:     limit,
					},
				},
			},
			wantCode: http.StatusOK,
		},
		{
			name:   "post not found",
			postID: 2,
			limit:  limit,
			offset: offset,
			mockBehaviour: func() {
				dependencies.commentService.EXPECT().
					GetAllComments(mock.Anything, core.GetAllCommentsParams{
						PostID: 2,
						Limit:  &limit,
						Offset: &offset,
					}).
					Return(nil, 0, core.ErrPostNotFound).Once()
			},
			wantCode: http.StatusNotFound,
		},
		{
			name:   "internal error",
			postID: 1,
			limit:  limit,
			offset: offset,
			mockBehaviour: func() {
				dependencies.commentService.EXPECT().
					GetAllComments(mock.Anything, core.GetAllCommentsParams{
						PostID: 1,
						Limit:  &limit,
						Offset: &offset,
					}).
					Return(nil, 0, errors.New("internal error")).Once()
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:          "invalid limit",
			postID:        1,
			limit:         invalidLimit,
			offset:        offset,
			mockBehaviour: func() {},
			wantCode:      http.StatusUnprocessableEntity,
		},
		{
			name:          "invalid offset",
			postID:        1,
			limit:         limit,
			offset:        -1,
			mockBehaviour: func() {},
			wantCode:      http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf(route, tt.postID, tt.limit, tt.offset), http.NoBody)

			resp, err := app.Test(req, -1)
			require.NoError(t, err)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			err = resp.Body.Close()
			require.NoError(t, err)

			if tt.wantCode == http.StatusOK {
				var data GetCommentsResponse

				err = json.Unmarshal(body, &data)
				require.NoError(t, err)

				assert.Equal(t, tt.wantData.Data.Meta, data.Data.Meta)
				assert.Equal(t, tt.wantData.Data.Comments, data.Data.Comments)
			}

			assert.Equal(t, tt.wantCode, resp.StatusCode)
		})
	}
}

func TestCreateComment(t *testing.T) {
	t.Parallel()
	app, dependencies := newTestApp(t)

	const route = "/api/v1/posts/%d/comments"

	tests := []struct {
		name          string
		postID        int
		token         string
		comment       comment.Create
		mockBehaviour func(comment.Create)
		wantCode      int
	}{
		{
			name:   "success",
			postID: 1,
			token:  token,
			comment: comment.Create{
				Content:  gofakeit.Sentence(10),
				ParentID: &[]int{gofakeit.Number(1, 10)}[0],
				ReplyID:  &[]int{gofakeit.Number(1, 10)}[0],
			},
			mockBehaviour: func(comment comment.Create) {
				coreComment := comment.ToCoreComment()
				coreComment.AuthorID = authorID
				coreComment.PostID = 1
				dependencies.commentService.EXPECT().
					CreateComment(mock.Anything, coreComment).
					Return(coreComment, nil).Once()
			},
			wantCode: http.StatusCreated,
		},
		{
			name:   "post not found",
			postID: 2,
			token:  token,
			comment: comment.Create{
				Content:  gofakeit.Sentence(10),
				ParentID: &[]int{gofakeit.Number(1, 10)}[0],
				ReplyID:  &[]int{gofakeit.Number(1, 10)}[0],
			},
			mockBehaviour: func(comment comment.Create) {
				coreComment := comment.ToCoreComment()
				coreComment.AuthorID = authorID
				coreComment.PostID = 2
				dependencies.commentService.EXPECT().
					CreateComment(mock.Anything, coreComment).
					Return(core.Comment{}, core.ErrPostNotFound).Once()
			},
			wantCode: http.StatusNotFound,
		},
		{
			name:   "internal error",
			postID: 1,
			token:  token,
			comment: comment.Create{
				Content:  gofakeit.Sentence(10),
				ParentID: &[]int{gofakeit.Number(1, 10)}[0],
				ReplyID:  &[]int{gofakeit.Number(1, 10)}[0],
			},
			mockBehaviour: func(comment comment.Create) {
				coreComment := comment.ToCoreComment()
				coreComment.AuthorID = authorID
				coreComment.PostID = 1
				dependencies.commentService.EXPECT().
					CreateComment(mock.Anything, coreComment).
					Return(core.Comment{}, errors.New("internal error")).Once()
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:   "invalid content",
			postID: 1,
			token:  token,
			comment: comment.Create{
				Content:  "",
				ParentID: &[]int{gofakeit.Number(1, 10)}[0],
				ReplyID:  &[]int{gofakeit.Number(1, 10)}[0],
			},
			mockBehaviour: func(comment comment.Create) {},
			wantCode:      http.StatusUnprocessableEntity,
		},
		{
			name:   "invalid post id",
			postID: 0,
			token:  token,
			comment: comment.Create{
				Content:  gofakeit.Sentence(10),
				ParentID: &[]int{gofakeit.Number(1, 10)}[0],
				ReplyID:  &[]int{gofakeit.Number(1, 10)}[0],
			},
			mockBehaviour: func(comment comment.Create) {},
			wantCode:      http.StatusUnprocessableEntity,
		},
		{
			name:          "invalid token",
			postID:        1,
			token:         "",
			mockBehaviour: func(comment comment.Create) {},
			wantCode:      http.StatusUnauthorized,
		},
		{
			name:   "parent comment not found",
			postID: 1,
			token:  token,
			comment: comment.Create{
				Content:  gofakeit.Sentence(10),
				ParentID: &[]int{1000}[0],
				ReplyID:  &[]int{gofakeit.Number(1, 10)}[0],
			},
			mockBehaviour: func(comment comment.Create) {
				coreComment := comment.ToCoreComment()
				coreComment.AuthorID = authorID
				coreComment.PostID = 1
				dependencies.commentService.EXPECT().
					CreateComment(mock.Anything, coreComment).
					Return(core.Comment{}, core.ErrParentCommentNotFound).Once()
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:   "reply comment not found",
			postID: 1,
			token:  token,
			comment: comment.Create{
				Content:  gofakeit.Sentence(10),
				ParentID: &[]int{gofakeit.Number(1, 10)}[0],
				ReplyID:  &[]int{1000}[0],
			},
			mockBehaviour: func(comment comment.Create) {
				coreComment := comment.ToCoreComment()
				coreComment.AuthorID = authorID
				coreComment.PostID = 1
				dependencies.commentService.EXPECT().
					CreateComment(mock.Anything, coreComment).
					Return(core.Comment{}, core.ErrReplyCommentNotFound).Once()
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:   "reply to comment of another post",
			postID: 1,
			token:  token,
			comment: comment.Create{
				Content:  gofakeit.Sentence(10),
				ParentID: &[]int{gofakeit.Number(1, 10)}[0],
				ReplyID:  &[]int{gofakeit.Number(1, 10)}[0],
			},
			mockBehaviour: func(comment comment.Create) {
				coreComment := comment.ToCoreComment()
				coreComment.AuthorID = authorID
				coreComment.PostID = 1
				dependencies.commentService.EXPECT().
					CreateComment(mock.Anything, coreComment).
					Return(core.Comment{}, core.ErrReplyToCommentOfAnotherPost).Once()
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:   "invalid reply comment",
			postID: 1,
			token:  token,
			comment: comment.Create{
				Content:  gofakeit.Sentence(10),
				ParentID: &[]int{gofakeit.Number(1, 10)}[0],
				ReplyID:  &[]int{gofakeit.Number(1, 10)}[0],
			},
			mockBehaviour: func(comment comment.Create) {
				coreComment := comment.ToCoreComment()
				coreComment.AuthorID = authorID
				coreComment.PostID = 1
				dependencies.commentService.EXPECT().
					CreateComment(mock.Anything, coreComment).
					Return(core.Comment{}, core.ErrInvalidCommentReplyID).Once()
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:   "invalid parent comment",
			postID: 1,
			token:  token,
			comment: comment.Create{
				Content:  gofakeit.Sentence(10),
				ParentID: &[]int{gofakeit.Number(1, 10)}[0],
				ReplyID:  &[]int{gofakeit.Number(1, 10)}[0],
			},
			mockBehaviour: func(comment comment.Create) {
				coreComment := comment.ToCoreComment()
				coreComment.AuthorID = authorID
				coreComment.PostID = 1
				dependencies.commentService.EXPECT().
					CreateComment(mock.Anything, coreComment).
					Return(core.Comment{}, core.ErrInvalidCommentParentID).Once()
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:   "null parent comment id",
			postID: 1,
			token:  token,
			comment: comment.Create{
				Content:  gofakeit.Sentence(10),
				ParentID: nil,
				ReplyID:  &[]int{gofakeit.Number(1, 10)}[0],
			},
			mockBehaviour: func(comment comment.Create) {
				coreComment := comment.ToCoreComment()
				coreComment.AuthorID = authorID
				coreComment.PostID = 1
				dependencies.commentService.EXPECT().
					CreateComment(mock.Anything, coreComment).
					Return(core.Comment{}, core.ErrNullCommentParentID).Once()
			},
			wantCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.comment)

			body, err := json.Marshal(tt.comment)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf(route, tt.postID), bytes.NewBuffer(body))

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))

			resp, err := app.Test(req, -1)
			require.NoError(t, err)

			body, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
			err = resp.Body.Close()
			require.NoError(t, err)

			if tt.wantCode == http.StatusOK {
				var data comment.Comment

				err = json.Unmarshal(body, &data)
				require.NoError(t, err)

				assert.Equal(t, tt.comment, data)
			}

			assert.Equal(t, tt.wantCode, resp.StatusCode)
		})
	}
}

func TestUpdateComment(t *testing.T) {
	t.Parallel()
	app, dependencies := newTestApp(t)

	const route = "/api/v1/posts/%d/comments/%d"

	type UpdateCommentResponse struct {
		Data comment.Comment `json:"data"`
	}

	tests := []struct {
		name          string
		postID        int
		commentID     int
		token         string
		comment       comment.Update
		mockBehaviour func(comment.Update) core.Comment
		wantCode      int
	}{
		{
			name:      "success",
			postID:    1,
			commentID: 1,
			token:     token,
			comment: comment.Update{
				Content: gofakeit.Sentence(10),
			},
			mockBehaviour: func(comment comment.Update) core.Comment {
				coreComment := comment.ToCoreComment()
				coreComment.AuthorID = authorID
				coreComment.PostID = 1
				coreComment.ID = 1
				dependencies.commentService.EXPECT().
					UpdateComment(mock.Anything, coreComment).
					Return(coreComment, nil).Once()
				return coreComment
			},
			wantCode: http.StatusOK,
		},
		{
			name:      "post id mismatch",
			postID:    2,
			commentID: 1,
			token:     token,
			comment: comment.Update{
				Content: gofakeit.Sentence(10),
			},
			mockBehaviour: func(comment comment.Update) core.Comment {
				coreComment := comment.ToCoreComment()
				coreComment.AuthorID = authorID
				coreComment.PostID = 2
				coreComment.ID = 1
				dependencies.commentService.EXPECT().
					UpdateComment(mock.Anything, coreComment).
					Return(core.Comment{}, core.ErrCommentPostIDMismatch).Once()
				return coreComment
			},
			wantCode: http.StatusNotFound,
		},
		{
			name:      "comment not found",
			postID:    1,
			commentID: 2,
			token:     token,
			comment: comment.Update{
				Content: gofakeit.Sentence(10),
			},
			mockBehaviour: func(comment comment.Update) core.Comment {
				coreComment := comment.ToCoreComment()
				coreComment.AuthorID = authorID
				coreComment.PostID = 1
				coreComment.ID = 2
				dependencies.commentService.EXPECT().
					UpdateComment(mock.Anything, coreComment).
					Return(core.Comment{}, core.ErrNoSuchComment).Once()
				return coreComment
			},
			wantCode: http.StatusNotFound,
		},
		{
			name:      "internal error",
			postID:    1,
			commentID: 1,
			token:     token,
			comment: comment.Update{
				Content: gofakeit.Sentence(10),
			},
			mockBehaviour: func(comment comment.Update) core.Comment {
				coreComment := comment.ToCoreComment()
				coreComment.AuthorID = authorID
				coreComment.PostID = 1
				coreComment.ID = 1
				dependencies.commentService.EXPECT().
					UpdateComment(mock.Anything, coreComment).
					Return(core.Comment{}, errors.New("internal error")).Once()
				return coreComment
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:      "comment author id mismatch",
			postID:    1,
			commentID: 1,
			token:     token,
			comment: comment.Update{
				Content: gofakeit.Sentence(10),
			},
			mockBehaviour: func(comment comment.Update) core.Comment {
				coreComment := comment.ToCoreComment()
				coreComment.AuthorID = authorID
				coreComment.PostID = 1
				coreComment.ID = 1
				dependencies.commentService.EXPECT().
					UpdateComment(mock.Anything, coreComment).
					Return(core.Comment{}, core.ErrCommentAuthorIDMismatch).Once()
				return coreComment
			},
			wantCode: http.StatusForbidden,
		},
		{
			name:      "invalid content",
			postID:    1,
			commentID: 1,
			token:     token,
			comment: comment.Update{
				Content: "",
			},
			mockBehaviour: func(comment comment.Update) core.Comment { return core.Comment{} },
			wantCode:      http.StatusUnprocessableEntity,
		},
		{
			name:      "invalid token",
			postID:    1,
			commentID: 1,
			token:     "",
			comment: comment.Update{
				Content: gofakeit.Sentence(10),
			},
			mockBehaviour: func(comment comment.Update) core.Comment { return core.Comment{} },
			wantCode:      http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coreComment := tt.mockBehaviour(tt.comment)

			body, err := json.Marshal(tt.comment)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf(route, tt.postID, tt.commentID), bytes.NewBuffer(body))

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))

			resp, err := app.Test(req, -1)
			require.NoError(t, err)

			body, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
			err = resp.Body.Close()
			require.NoError(t, err)

			if tt.wantCode == http.StatusOK {
				var data UpdateCommentResponse

				err = json.Unmarshal(body, &data)
				require.NoError(t, err)

				assert.Equal(t, comment.ToModelComment(coreComment), data.Data)
			}

			assert.Equal(t, tt.wantCode, resp.StatusCode)
		})
	}
}

func TestDeleteComment(t *testing.T) {
	t.Parallel()
	app, dependencies := newTestApp(t)

	const route = "/api/v1/posts/%d/comments/%d"

	tests := []struct {
		name          string
		postID        int
		commentID     int
		token         string
		mockBehaviour func()
		wantCode      int
	}{
		{
			name:      "success",
			postID:    1,
			commentID: 1,
			token:     token,
			mockBehaviour: func() {
				dependencies.commentService.EXPECT().
					DeleteComment(mock.Anything, core.Comment{
						ID:       1,
						AuthorID: authorID,
						PostID:   1,
					}).
					Return(nil).Once()
			},
			wantCode: http.StatusNoContent,
		},
		{
			name:      "post id mismatch",
			postID:    2,
			commentID: 1,
			token:     token,
			mockBehaviour: func() {
				dependencies.commentService.EXPECT().
					DeleteComment(mock.Anything, core.Comment{
						ID:       1,
						AuthorID: authorID,
						PostID:   2,
					}).
					Return(core.ErrCommentPostIDMismatch).Once()
			},
			wantCode: http.StatusNoContent,
		},
		{
			name:      "comment not found",
			postID:    1,
			commentID: 2,
			token:     token,
			mockBehaviour: func() {
				dependencies.commentService.EXPECT().
					DeleteComment(mock.Anything, core.Comment{
						ID:       2,
						AuthorID: authorID,
						PostID:   1,
					}).
					Return(core.ErrNoSuchComment).Once()
			},
			wantCode: http.StatusNoContent,
		},
		{
			name:      "internal error",
			postID:    1,
			commentID: 1,
			token:     token,
			mockBehaviour: func() {
				dependencies.commentService.EXPECT().
					DeleteComment(mock.Anything, core.Comment{
						ID:       1,
						AuthorID: authorID,
						PostID:   1,
					}).
					Return(errors.New("internal error")).Once()
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:      "comment author id mismatch",
			postID:    1,
			commentID: 1,
			token:     token,
			mockBehaviour: func() {
				dependencies.commentService.EXPECT().
					DeleteComment(mock.Anything, core.Comment{
						ID:       1,
						AuthorID: authorID,
						PostID:   1,
					}).
					Return(core.ErrCommentAuthorIDMismatch).Once()
			},
			wantCode: http.StatusForbidden,
		},
		{
			name:          "invalid token",
			postID:        1,
			commentID:     1,
			token:         "",
			mockBehaviour: func() {},
			wantCode:      http.StatusUnauthorized,
		},
		{
			name:          "invalid post id",
			postID:        0,
			commentID:     1,
			token:         token,
			mockBehaviour: func() {},
			wantCode:      http.StatusUnprocessableEntity,
		},
		{
			name:          "invalid comment id",
			postID:        1,
			commentID:     0,
			token:         token,
			mockBehaviour: func() {},
			wantCode:      http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf(route, tt.postID, tt.commentID), http.NoBody)

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			err = resp.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.wantCode, resp.StatusCode)
		})
	}
}
