package comments

import (
	"context"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

type commentsService struct {
	CommentsStore core.CommentsStore
}

func NewCommentsService(store core.CommentsStore) core.CommentsService {
	return &commentsService{
		CommentsStore: store,
	}
}

func (s *commentsService) GetCommentByID(ctx context.Context, commentID int) (core.Comments, error) {
	return s.CommentsStore.GetCommentByID(ctx, commentID)
}

func (s *commentsService) GetCommentsByPostID(ctx context.Context, params core.GetAllParamsComments, postID int) (data []core.Comments, total int, err error) {
	comments, total, err := s.CommentsStore.GetCommentsByPostID(ctx, params, postID)
	return comments, total, err
}

func (s *commentsService) CreateComment(ctx context.Context, comment core.Comments) (core.Comments, error) {
	return s.CommentsStore.CreateComment(ctx, comment)
}

func (s *commentsService) UpdateComments(ctx context.Context, comment core.Comments) (core.Comments, error) {
	// checking if the author id and user id are equal
	dbComment, err := s.GetCommentByID(ctx, comment.ID)
	if err != nil {
		return comment, err
	}

	if dbComment.AuthorID != comment.AuthorID {
		return comment, core.ErrCommentAuthorIDMismatch
	} else if dbComment.IsDeleted {
		return comment, core.ErrCommentIsDeleted
	}

	return s.CommentsStore.UpdateComments(ctx, comment)
}

func (s *commentsService) DeleteComments(ctx context.Context, comments core.Comments) error {
	dbComment, err := s.GetCommentByID(ctx, comments.ID)
	if err != nil {
		return err
	}

	if dbComment.AuthorID != comments.AuthorID {
		return core.ErrCommentAuthorIDMismatch
	} else if dbComment.IsDeleted {
		return core.ErrCommentIsDeleted
	}

	return s.CommentsStore.DeleteComments(ctx, comments)
}
