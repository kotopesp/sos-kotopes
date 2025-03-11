package moderator

import (
	"context"
	"errors"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
	"gorm.io/gorm"
	"time"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.ModeratorStore {
	return &store{pg}
}

// GetModeratorByID - retrieves the moderator structure by their id.
func (s *store) GetModeratorByID(ctx context.Context, id int) (moderator core.Moderator, err error) {
	err = s.DB.WithContext(ctx).
		Table(moderator.TableName()).
		Where("id = ?", id).
		First(&moderator).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return moderator, core.ErrNoSuchModerator
		}
		logger.Log().Debug(ctx, err.Error())
		return moderator, err
	}

	return moderator, nil
}

// CreateModerator creates moderator
func (s *store) CreateModerator(ctx context.Context, moderator core.Moderator) (err error) {
	moderator.CreatedAt = time.Now().UTC()

	err = s.DB.WithContext(ctx).
		Create(&moderator).Error

	if err != nil {
		return err
	}

	return nil
}

// GetReasonsForReportedPost returns list of reasons why post was banned.
func (s *store) GetReasonsForReportedPost(ctx context.Context, postID int) (reasons []string, err error) {
	err = s.DB.WithContext(ctx).
		Table(core.Report{}.TableName()).
		Where("post_id = ?", postID).
		Pluck("reason", &reasons).Error
	if err != nil {
		logger.Log().Debug(ctx, err.Error())

		return nil, core.ErrGettingReportResponse
	}

	return reasons, nil
}

// GetPostsForModeration - takes first 10 records from posts table which status is "on_moderation"
func (s *store) GetPostsForModeration(ctx context.Context) ([]core.PostForModeration, error) {
	var posts []core.Post
	err := s.DB.WithContext(ctx).
		Where("status = ?", string(core.OnModeration)).
		Order("updated_at ASC").
		Limit(core.AmountOfPostsForModeration).
		Find(&posts).Error

	if err != nil {
		logger.Log().Error(ctx, err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrNoPostsWaitingForModeration
		}
		return nil, err
	}

	var postsForModeration []core.PostForModeration

	for _, post := range posts {
		reasons, err := s.GetReasonsForReportedPost(ctx, post.ID)
		if err != nil {
			logger.Log().Error(ctx, err.Error())

			continue
		}

		postsForModeration = append(postsForModeration, core.PostForModeration{
			Post:    post,
			Reasons: reasons,
		})
	}

	return postsForModeration, nil
}
