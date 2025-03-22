package moderator

import (
	"context"
	"fmt"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

type service struct {
	moderatorStore core.ModeratorStore
	postStore      core.PostStore
	reportStore    core.ReportStore
}

func New(moderatorStore core.ModeratorStore, postStore core.PostStore, reportStore core.ReportStore) core.ModeratorService {
	return &service{moderatorStore: moderatorStore, postStore: postStore, reportStore: reportStore}
}

// GetModerator - returns moderator struct by its id.
func (s *service) GetModerator(ctx context.Context, id int) (moderator core.Moderator, err error) {
	moderator, err = s.moderatorStore.GetModeratorByID(ctx, id)
	if err != nil {
		logger.Log().Error(ctx, err.Error())

		return core.Moderator{}, err
	}

	return moderator, nil
}

// GetPostsForModeration - returns a list of posts that were sorted by the time of the report in order core.Filter.
func (s *service) GetPostsForModeration(ctx context.Context, filter core.Filter) (moderationPosts []core.PostForModeration, err error) {
	posts, err := s.postStore.GetPostsForModeration(ctx, filter)
	if err != nil {
		logger.Log().Error(ctx, err.Error())

		return nil, err
	}

	for _, post := range posts {
		reasons, err := s.reportStore.GetReportReasonsForPost(ctx, post.ID)
		if err != nil {
			logger.Log().Error(ctx, fmt.Sprintf("Error getting report reasons for, %d: ", post.ID)+err.Error())

			continue
		}

		postWithReasons := core.PostForModeration{
			Post:    post,
			Reasons: reasons,
		}

		moderationPosts = append(moderationPosts, postWithReasons)
	}

	return moderationPosts, nil
}

// DeletePost - method that allows moderator to delete posts.
func (s *service) DeletePost(ctx context.Context, id int) (err error) {
	err = s.postStore.DeletePost(ctx, id)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}

	return nil
}

func (s *service) ApprovePost(ctx context.Context, postID int) (err error) {
	err = s.reportStore.DeleteAllReportsForPost(ctx, postID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}

	return nil
}
