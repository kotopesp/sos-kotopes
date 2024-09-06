package commentservice

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/internal/core/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func generateTestComments() []core.Comment {
	// Example users
	user1 := core.User{ID: 1, Username: "User1"}
	user2 := core.User{ID: 2, Username: "User2"}

	// Example timestamps
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)
	twoHoursAgo := now.Add(-2 * time.Hour)

	// Test comments
	comments := []core.Comment{
		{
			ID:        1,
			ParentID:  nil,
			ReplyID:   nil,
			PostID:    1,
			AuthorID:  user1.ID,
			Author:    user1,
			Content:   "This is the first comment.",
			IsDeleted: false,
			DeletedAt: nil,
			CreatedAt: twoHoursAgo,
			UpdatedAt: twoHoursAgo,
		},
		{
			ID:        2,
			ParentID:  nil,
			ReplyID:   &[]int{1}[0], // Reply to the first comment
			PostID:    1,
			AuthorID:  user2.ID,
			Author:    user2,
			Content:   "This is a reply to the first comment.",
			IsDeleted: false,
			DeletedAt: nil,
			CreatedAt: oneHourAgo,
			UpdatedAt: oneHourAgo,
		},
		{
			ID:        3,
			ParentID:  nil,
			ReplyID:   nil,
			PostID:    1,
			AuthorID:  user2.ID,
			Author:    user2,
			Content:   "This is another comment.",
			IsDeleted: true, // Deleted comment
			DeletedAt: &now,
			CreatedAt: oneHourAgo,
			UpdatedAt: now,
		},
		{
			ID:        4,
			ParentID:  &[]int{1}[0], // This is a child of the first comment
			ReplyID:   nil,
			PostID:    1,
			AuthorID:  user1.ID,
			Author:    user1,
			Content:   "This is a nested comment.",
			IsDeleted: false,
			DeletedAt: nil,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	return comments
}

func generateTestPost() core.Post {
	// Example timestamps
	now := time.Now()
	twoDaysAgo := now.Add(-48 * time.Hour)

	// Example photo as a byte slice (you would replace this with actual photo data in a real scenario)
	photo := []byte{137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13, 73, 72, 68, 82} // Example PNG file header

	// Generate a test post
	post := core.Post{
		ID:        1,
		Title:     "Test Post Title",
		Content:   "This is the content of the test post. It discusses various interesting topics.",
		AuthorID:  1,
		AnimalID:  101,
		IsDeleted: false,
		DeletedAt: time.Time{}, // Zero value means it is not deleted
		CreatedAt: twoDaysAgo,
		UpdatedAt: now,
		Photo:     photo,
	}

	return post
}

func generateTestGetAllCommentsParams() core.GetAllCommentsParams {
	// Example values for limit and offset
	limit := 10
	offset := 20

	// Generate test parameters
	params := core.GetAllCommentsParams{
		PostID: 1,       // Example PostID
		Limit:  &limit,  // Set limit to 10 comments
		Offset: &offset, // Set offset to skip the first 20 comments
	}

	return params
}

func TestGetAllComments(t *testing.T) {
	ctx := context.Background()

	commentStore := mocks.NewCommentStore(t)
	postStore := mocks.NewPostStore(t)

	var (
		comments     = generateTestComments()
		post         = generateTestPost()
		emptyPost    = core.Post{}
		postError    = core.ErrPostNotFound
		commentError = core.ErrNoSuchComment
		params       = generateTestGetAllCommentsParams()
	)

	commentService := New(
		commentStore,
		postStore,
	)

	tests := []struct {
		name                 string
		retComments          []core.Comment
		retTotal             int
		retCommentError      error
		retPost              core.Post
		retPostError         error
		wantComments         []core.Comment
		invokeGetAllComments bool
		invokeGetPostByID    bool
		wantTotal            int
		wantError            error
	}{
		{
			name:                 "success",
			retComments:          comments,
			retTotal:             len(comments),
			retPost:              post,
			invokeGetAllComments: true,
			invokeGetPostByID:    true,
			wantComments:         comments,
			wantTotal:            len(comments),
		},
		{
			name:                 "post store fail",
			retPost:              emptyPost,
			retPostError:         postError,
			invokeGetAllComments: false,
			invokeGetPostByID:    true,
			wantError:            postError,
		},
		{
			name:                 "comment store fail",
			retPost:              post,
			invokeGetAllComments: true,
			invokeGetPostByID:    true,
			retCommentError:      commentError,
			wantError:            commentError,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if tt.invokeGetAllComments {
				commentStore.On(
					"GetAllComments",
					mock.Anything,
					params,
				).Return(tt.retComments, tt.retTotal, tt.retCommentError).Once()
			}

			if tt.invokeGetPostByID {
				postStore.On(
					"GetPostByID",
					mock.Anything,
					params.PostID,
				).Return(tt.retPost, tt.retPostError).Once()
			}

			res, total, err := commentService.GetAllComments(ctx, params)

			assert.ErrorIs(t, err, tt.wantError, "errors should match")
			if tt.wantError == nil {
				assert.Equal(t, tt.wantComments, res, "comments should be equal")
				assert.Equal(t, tt.wantTotal, total, "total amount of comments should be equal")
			}
		})
	}
}

func TestCreateComment(t *testing.T) {
	ctx := context.Background()

	commentStore := mocks.NewCommentStore(t)
	postStore := mocks.NewPostStore(t)

	var (
		comment      = generateTestComments()[0]
		emptyComment = core.Comment{}
		post         = generateTestPost()
		emptyPost    = core.Post{}
		postError    = core.ErrPostNotFound
		commentError = errors.New("some comment error")
		params       = generateTestGetAllCommentsParams()
	)

	commentService := New(
		commentStore,
		postStore,
	)

	tests := []struct {
		name                string
		retComment          core.Comment
		retCommentError     error
		retPost             core.Post
		retPostError        error
		invokeCreateComment bool
		invokeGetPostByID   bool
		wantError           error
		wantComment         core.Comment
	}{
		{
			name:                "success",
			retComment:          comment,
			retPost:             post,
			invokeCreateComment: true,
			invokeGetPostByID:   true,
			wantComment:         comment,
		},
		{
			name:                "post store fail",
			retPost:             emptyPost,
			retPostError:        postError,
			invokeCreateComment: false,
			invokeGetPostByID:   true,
			wantError:           postError,
		},
		{
			name:                "comment store fail",
			retPost:             post,
			invokeCreateComment: true,
			invokeGetPostByID:   true,
			retComment:          emptyComment,
			retCommentError:     commentError,
			wantError:           commentError,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if tt.invokeCreateComment {
				commentStore.On(
					"CreateComment",
					mock.Anything,
					comment,
				).Return(tt.retComment, tt.retCommentError).Once()
			}

			if tt.invokeGetPostByID {
				postStore.On(
					"GetPostByID",
					mock.Anything,
					params.PostID,
				).Return(tt.retPost, tt.retPostError).Once()
			}

			res, err := commentService.CreateComment(ctx, comment)

			assert.ErrorIs(t, err, tt.wantError, "errors should match")
			if err == nil {
				assert.Equal(t, tt.wantComment, res, "comments should be equal")
			}
		})
	}
}

func TestUpdateComment(t *testing.T) {
	ctx := context.Background()

	commentStore := mocks.NewCommentStore(t)
	postStore := mocks.NewPostStore(t)

	var (
		comment                 = generateTestComments()[0]
		commentAuthorIDMismatch = comment
		commentPostIDMisMatch   = comment
		commentDeleted          = comment
		commentError            = errors.New("some comment error")
	)

	// author id mismatch
	commentAuthorIDMismatch.AuthorID = comment.ID + 1

	// post id mismatch
	commentPostIDMisMatch.PostID = comment.PostID + 1

	// deleted comment
	commentDeleted.IsDeleted = true

	commentService := New(
		commentStore,
		postStore,
	)

	tests := []struct {
		name                     string
		retGetCommentByIDComment core.Comment
		retGetCommentByIDError   error
		invokeGetCommentByID     bool
		retUpdateCommentComment  core.Comment
		retUpdateCommentError    error
		invokeUpdateComment      bool
		wantComment              core.Comment
		wantError                error
	}{
		{
			name:                     "success",
			retGetCommentByIDComment: comment,
			invokeGetCommentByID:     true,
			retUpdateCommentComment:  comment,
			invokeUpdateComment:      true,
			wantComment:              comment,
		},
		{
			name:                     "get comment by id error",
			retGetCommentByIDComment: comment,
			retGetCommentByIDError:   commentError,
			invokeGetCommentByID:     true,
			invokeUpdateComment:      false,
			wantComment:              comment,
			wantError:                commentError,
		},
		{
			name:                     "author id mismatch",
			retGetCommentByIDComment: commentAuthorIDMismatch,
			invokeGetCommentByID:     true,
			invokeUpdateComment:      false,
			wantComment:              comment,
			wantError:                core.ErrCommentAuthorIDMismatch,
		},
		{
			name:                     "post id mismatch",
			retGetCommentByIDComment: commentPostIDMisMatch,
			invokeGetCommentByID:     true,
			invokeUpdateComment:      false,
			wantComment:              comment,
			wantError:                core.ErrCommentPostIDMismatch,
		},
		{
			name:                     "deleted comment",
			retGetCommentByIDComment: commentDeleted,
			invokeGetCommentByID:     true,
			invokeUpdateComment:      false,
			wantComment:              comment,
			wantError:                core.ErrCommentIsDeleted,
		},
		{
			name:                     "update comment error",
			retGetCommentByIDComment: comment,
			invokeGetCommentByID:     true,
			retUpdateCommentComment:  comment,
			retUpdateCommentError:    commentError,
			invokeUpdateComment:      true,
			wantComment:              comment,
			wantError:                commentError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.invokeGetCommentByID {
				commentStore.On(
					"GetCommentByID",
					mock.Anything,
					comment.ID,
				).Return(tt.retGetCommentByIDComment, tt.retGetCommentByIDError).Once()
			}

			if tt.invokeUpdateComment {
				commentStore.On(
					"UpdateComment",
					mock.Anything,
					comment,
				).Return(tt.retUpdateCommentComment, tt.retUpdateCommentError).Once()
			}

			res, err := commentService.UpdateComment(ctx, comment)

			assert.ErrorIs(t, err, tt.wantError, "errors should match")
			if err == nil {
				assert.Equal(t, tt.wantComment, res, "comments should be equal")
			}
		})
	}
}

func TestDeleteComment(t *testing.T) {
	ctx := context.Background()

	commentStore := mocks.NewCommentStore(t)
	postStore := mocks.NewPostStore(t)

	var (
		comment                 = generateTestComments()[0]
		commentAuthorIDMismatch = comment
		commentPostIDMisMatch   = comment
		commentDeleted          = comment
		commentError            = errors.New("some comment error")
	)

	// author id mismatch
	commentAuthorIDMismatch.AuthorID = comment.ID + 1

	// post id mismatch
	commentPostIDMisMatch.PostID = comment.PostID + 1

	// deleted comment
	commentDeleted.IsDeleted = true

	commentService := New(
		commentStore,
		postStore,
	)

	tests := []struct {
		name                     string
		retGetCommentByIDComment core.Comment
		retGetCommentByIDError   error
		invokeGetCommentByID     bool
		retDeleteCommentError    error
		invokeDeleteComment      bool
		wantError                error
	}{
		{
			name:                     "success",
			retGetCommentByIDComment: comment,
			invokeGetCommentByID:     true,
			invokeDeleteComment:      true,
		},
		{
			name:                     "get comment by id error",
			retGetCommentByIDComment: comment,
			retGetCommentByIDError:   commentError,
			invokeGetCommentByID:     true,
			invokeDeleteComment:      false,
			wantError:                commentError,
		},
		{
			name:                     "author id mismatch",
			retGetCommentByIDComment: commentAuthorIDMismatch,
			invokeGetCommentByID:     true,
			invokeDeleteComment:      false,
			wantError:                core.ErrCommentAuthorIDMismatch,
		},
		{
			name:                     "post id mismatch",
			retGetCommentByIDComment: commentPostIDMisMatch,
			invokeGetCommentByID:     true,
			invokeDeleteComment:      false,
			wantError:                core.ErrCommentPostIDMismatch,
		},
		{
			name:                     "deleted comment",
			retGetCommentByIDComment: commentDeleted,
			invokeGetCommentByID:     true,
			invokeDeleteComment:      false,
			wantError:                core.ErrCommentIsDeleted,
		},
		{
			name:                     "delete comment error",
			retGetCommentByIDComment: comment,
			invokeGetCommentByID:     true,
			retDeleteCommentError:    commentError,
			invokeDeleteComment:      true,
			wantError:                commentError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.invokeGetCommentByID {
				commentStore.On(
					"GetCommentByID",
					mock.Anything,
					comment.ID,
				).Return(tt.retGetCommentByIDComment, tt.retGetCommentByIDError).Once()
			}

			if tt.invokeDeleteComment {
				commentStore.On(
					"DeleteComment",
					mock.Anything,
					comment,
				).Return(tt.retDeleteCommentError).Once()
			}

			err := commentService.DeleteComment(ctx, comment)

			assert.ErrorIs(t, err, tt.wantError, "errors should match")
		})
	}
}
