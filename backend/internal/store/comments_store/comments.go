package comments

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"strings"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
)

type store struct {
	*postgres.Postgres
}

func NewCommentsStore(pg *postgres.Postgres) core.CommentsStore {
	return &store{
		pg,
	}
}

func (s *store) GetCommentByID(ctx context.Context, commentID int) (core.Comments, error) {
	var comment core.Comments
	if err := s.DB.WithContext(ctx).First(&comment, "id=?", commentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return comment, core.ErrNoSuchComment
		}
		return comment, err
	}

	return comment, nil
}

func (s *store) GetCommentsByPostID(ctx context.Context, params core.GetAllParamsComments, postID int) (data []core.Comments, total int, err error) {
	var comments []core.Comments

	query := s.DB.WithContext(ctx).Order(
		"COALESCE(parent_id, id), (parent_id IS NULL)::int DESC, id",
	) // sorting by comment and replies to it

	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}

	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	if err := query.Find(&comments, "posts_id=?", postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, core.ErrNoSuchComment
		}
		return nil, 0, err
	}

	var totalInt64 int64
	if err := s.DB.WithContext(ctx).
		Model(&core.Comments{}).
		Where("posts_id=?", postID).
		Count(&totalInt64).Error; err != nil {
		return nil, 0, err
	}

	return comments, int(totalInt64), nil
}

func (s *store) CreateComment(ctx context.Context, comment core.Comments) (core.Comments, error) {
	if err := s.DB.WithContext(ctx).Create(&comment).Error; err != nil {
		if strings.Contains(err.Error(), "comments_posts_id_fkey") {
			return comment, core.ErrNoSuchPost
		}
		return comment, err
	}

	return comment, nil
}

func (s *store) UpdateComments(ctx context.Context, comment core.Comments) (core.Comments, error) {
	if err := s.DB.WithContext(ctx).Updates(comment).Error; err != nil {
		return comment, err
	}

	// unfortunately, updates does not update `comment` variable
	if err := s.DB.WithContext(ctx).First(&comment, "id=?", comment.ID).Error; err != nil {
		return comment, err
	}

	return comment, nil
}

func (s *store) DeleteComments(ctx context.Context, comment core.Comments) error {
	comment.IsDeleted = true
	now := time.Now()
	comment.DeletedAt = &now

	if err := s.DB.WithContext(ctx).Updates(&comment).Error; err != nil {
		return err
	}

	return nil
}
