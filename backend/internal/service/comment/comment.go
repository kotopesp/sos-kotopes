package commentservice

import (
	"context"

	"github.com/kotopesp/sos-kotopes/internal/core"
)

type service struct {
	commentStore core.CommentStore
	postStore    core.PostStore
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

func (s *service) GetAllComments(ctx context.Context, params core.GetAllCommentsParams, userID int) (data []core.Comment, total int, err error) {
	if _, err := s.postStore.GetPostByID(ctx, params.PostID, userID); err != nil {
		return nil, 0, err
	}

	return s.commentStore.GetAllComments(ctx, params)
}

func (s *service) CreateComment(ctx context.Context, comment core.Comment, userID int) (data core.Comment, err error) {
	if _, err := s.postStore.GetPostByID(ctx, comment.PostID, userID); err != nil {
		return comment, err
	}

	if comment.ParentID != nil {
		dbComment, err := s.commentStore.GetCommentByID(ctx, *comment.ParentID)
		if err != nil {
			return comment, core.ErrParentCommentNotFound
		}

		if dbComment.PostID != comment.PostID {
			return comment, core.ErrReplyToCommentOfAnotherPost
		}

		if dbComment.ParentID != nil {
			return comment, core.ErrInvalidCommentParentID
		}
	}

	if comment.ReplyID != nil {
		if comment.ParentID == nil {
			return comment, core.ErrNullCommentParentID
		}

		dbComment, err := s.commentStore.GetCommentByID(ctx, *comment.ReplyID)
		if err != nil {
			return comment, core.ErrReplyCommentNotFound
		}

		if dbComment.PostID != comment.PostID {
			return comment, core.ErrReplyToCommentOfAnotherPost
		}

		if dbComment.ParentID == nil || *dbComment.ParentID != *comment.ParentID {
			return comment, core.ErrInvalidCommentReplyID
		}
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
