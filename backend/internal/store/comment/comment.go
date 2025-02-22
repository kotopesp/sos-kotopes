package commentstore

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.CommentStore {
	return &store{pg}
}

func (s *store) GetCommentByID(ctx context.Context, commentID int) (core.Comment, error) {
	var comment core.Comment
	if err := s.DB.WithContext(ctx).
		First(&comment, "id=?", commentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return comment, core.ErrNoSuchComment
		}
		return comment, err
	}

	return comment, nil
}

func (s *store) GetAllComments(ctx context.Context, params core.GetAllCommentsParams) (data []core.Comment, total int, err error) {
	var comments []core.Comment

	query := s.DB.WithContext(ctx).
		Model(&core.Comment{}).
		Where("posts_id=?", params.PostID)

	query = query.Order(
		"COALESCE(parent_id, id), (parent_id IS NULL)::int DESC, id",
	) // sorting by comment and replies to it

	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}

	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	if err := query.
		Preload("Author").
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	query = s.DB.WithContext(ctx).
		Model(&core.Comment{}).
		Where("posts_id=?", params.PostID)

	var totalInt64 int64
	if err := query.Count(&totalInt64).Error; err != nil {
		return nil, 0, err
	}

	return comments, int(totalInt64), nil
}

func (s *store) CreateComment(ctx context.Context, comment core.Comment) (core.Comment, error) {
	if err := s.DB.WithContext(ctx).Create(&comment).Error; err != nil {
		return comment, err
	}

	if err := s.DB.WithContext(ctx).Preload("Author").First(&comment).Error; err != nil {
		return comment, err
	}

	return comment, nil
}

func (s *store) UpdateComment(ctx context.Context, comment core.Comment) (core.Comment, error) {
	if err := s.DB.WithContext(ctx).Updates(comment).Error; err != nil {
		return comment, err
	}

	// unfortunately, updates does not update `comment_service` variable
	if err := s.DB.WithContext(ctx).
		Preload("Author").
		First(&comment, "id=?", comment.ID).Error; err != nil {
		return comment, err
	}

	return comment, nil
}

func (s *store) DeleteComment(ctx context.Context, comment core.Comment) error {
	comment.IsDeleted = true
	comment.DeletedAt = time.Now().UTC()

	if err := s.DB.WithContext(ctx).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}
