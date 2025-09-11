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
	commentStore   core.CommentStore
	userStore      core.UserStore
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
		reasons, err := s.reportStore.GetReportReasons(ctx, post.ID, core.ReportableTypePost)
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
	err = s.postStore.ApprovePostFromModeration(ctx, postID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}

	err = s.reportStore.DeleteAllReports(ctx, postID, core.ReportableTypePost)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}

	return nil
}

func (s *service) GetCommentsForModeration(ctx context.Context, filter core.Filter) ([]core.CommentForModeration, error) {
	comments, err := s.commentStore.GetCommentsForModeration(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(comments) == 0 {
		return nil, core.ErrNoCommentsWaitingForModeration
	}

	var result []core.CommentForModeration
	for _, comment := range comments {
		reasons, err := s.reportStore.GetReportReasons(ctx, comment.ID, core.ReportableTypeComment)
		if err != nil {
			return nil, err
		}

		result = append(result, core.CommentForModeration{
			Comment: comment,
			Reasons: reasons,
		})
	}

	return result, nil
}

func (s *service) DeleteComment(ctx context.Context, commentID int) error {
	comment, err := s.commentStore.GetCommentByID(ctx, commentID)
	if err != nil {
		return core.ErrNoSuchComment
	}

	if err := s.commentStore.DeleteComment(ctx, comment); err != nil {
		return err
	}

	if err := s.reportStore.DeleteAllReports(ctx, commentID, core.ReportableTypeComment); err != nil {
		logger.Log().Error(ctx, "Failed to delete reports for comment: "+err.Error())
	}

	return nil
}

func (s *service) ApproveComment(ctx context.Context, commentID int) error {
	_, err := s.commentStore.GetCommentByID(ctx, commentID)
	if err != nil {
		return core.ErrNoSuchComment
	}

	if err := s.commentStore.ApproveCommentFromModeration(ctx, commentID); err != nil {
		return err
	}

	if err := s.reportStore.DeleteAllReports(ctx, commentID, core.ReportableTypeComment); err != nil {
		logger.Log().Error(ctx, "Failed to delete reports for comment: "+err.Error())
	}

	return nil
}

func (s *service) BanUser(ctx context.Context, userID, reportID int) error {
	_, err := s.userStore.GetUserByID(ctx, userID)
	if err != nil {
		return core.ErrNoSuchUser
	}

	return nil
}
