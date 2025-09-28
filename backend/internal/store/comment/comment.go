package commentstore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
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
		First(&comment).
		Where("id = ?", commentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return comment, core.ErrNoSuchComment
		}
		return core.Comment{}, err
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

	comment.Status = core.Deleted
	comment.UpdatedAt = time.Now().UTC()

	if err := s.DB.WithContext(ctx).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (s *store) SendToModeration(ctx context.Context, commentID int) error {
	return s.DB.WithContext(ctx).
		Model(&core.Comment{}).
		Where("id = ?", commentID).
		Updates(map[string]interface{}{
			"status":     core.OnModeration,
			"updated_at": time.Now().UTC(),
		}).Error
}

// GetCommentsForModeration - takes amount of records limited by the constant core.AmountOfCommentsForModeration
// from comments table which status is "on_moderation"
func (s *store) GetCommentsForModeration(ctx context.Context, filter core.Filter) ([]core.Comment, error) {
	var comments []core.Comment

	err := s.DB.WithContext(ctx).
		Where("status = ?", core.OnModeration).
		Order("updated_at " + filter).
		Limit(core.AmountOfCommentsForModeration).
		Find(&comments).Error

	if err != nil {
		logger.Log().Error(ctx, "Failed to get comments for moderation: "+err.Error())
		if errors.Is(err, core.ErrRecordNotFound) {
			return nil, core.ErrNoCommentsWaitingForModeration
		}
		return nil, err
	}

	return comments, nil
}

// ApproveCommentFromModeration - changes comment status from "on_moderation" to "published"
func (s *store) ApproveCommentFromModeration(ctx context.Context, commentID int) error {
	updates := map[string]interface{}{
		"status":     core.Published,
		"updated_at": time.Now().UTC(),
	}

	result := s.DB.WithContext(ctx).
		Model(&core.Comment{}).
		Where("id = ? AND status = ?", commentID, core.OnModeration).
		Updates(updates)

	if result.Error != nil {
		logger.Log().Error(ctx, "Failed to approve comment: "+result.Error.Error())
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return core.ErrNoSuchComment
		}
		return result.Error
	}

	if result.RowsAffected == 0 {
		logger.Log().Error(ctx, fmt.Sprintf("Comment not found or not on moderation: %d", commentID))
		return core.ErrNoSuchComment
	}

	return nil
}
