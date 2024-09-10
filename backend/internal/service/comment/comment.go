package commentservice

import (
	"context"

	"github.com/kotopesp/sos-kotopes/internal/core"
)

type service struct {
	commentStore core.CommentStore
	postStore    core.PostStore
	userStore    core.UserStore
}

func New(
	commentStore core.CommentStore,
	postStore core.PostStore,
) core.CommentService {
	return &service{
		commentStore: commentStore,
		postStore:    postStore,
	}
}

func (s *service) GetAllComments(ctx context.Context, params core.GetAllCommentsParams) (data []core.Comment, total int, err error) {
	if _, err := s.postStore.GetPostByID(ctx, params.PostID); err != nil {
		return nil, 0, err
	}

	return s.commentStore.GetAllComments(ctx, params)
}

func (s *service) CreateComment(ctx context.Context, comment core.Comment) (data core.Comment, err error) {
	if _, err := s.postStore.GetPostByID(ctx, comment.PostID); err != nil {
		return comment, err
	}

	return s.commentStore.CreateComment(ctx, comment)
}

func (s *service) UpdateComment(ctx context.Context, comment core.Comment) (data core.Comment, err error) {
	dbComment, err := s.commentStore.GetCommentByID(ctx, comment.ID)
	if err != nil {
		return comment, err
	}

	if dbComment.AuthorID != comment.AuthorID {
		return comment, core.ErrCommentAuthorIDMismatch
	} else if dbComment.PostID != comment.PostID {
		return comment, core.ErrCommentPostIDMismatch
	} else if dbComment.IsDeleted {
		return comment, core.ErrCommentIsDeleted
	}

	return s.commentStore.UpdateComment(ctx, comment)
}

func (s *service) DeleteComment(ctx context.Context, comment core.Comment) error {
	dbComment, err := s.commentStore.GetCommentByID(ctx, comment.ID)
	if err != nil {
		return err
	}

	if dbComment.AuthorID != comment.AuthorID {
		return core.ErrCommentAuthorIDMismatch
	} else if dbComment.PostID != comment.PostID {
		return core.ErrCommentPostIDMismatch
	} else if dbComment.IsDeleted {
		return core.ErrCommentIsDeleted
	}

	return s.commentStore.DeleteComment(ctx, comment)
}
