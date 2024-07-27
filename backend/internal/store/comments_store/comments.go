package comments

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
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

func (s *store) GetCommentsByPostID(ctx context.Context, params core.GetAllParamsComments, post_id int) ([]core.Comments, error) {

	var comments []core.Comments

	query := s.DB.WithContext(ctx).Model(&core.Comments{})

	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}

	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	if err := query.Order("replay_to_comment, reply_to_reply").Find(&comments).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return comments, nil

}

func (s *store) CreateComment(ctx context.Context, comment core.Comments, post_id int) (core.Comments, error) {
	/*if comment.Is_deleted == true {
		return comment, fiber.NewError(fiber.StatusBadRequest, "invalid is_deleted comment")
	}*/
	if err := s.DB.WithContext(ctx).Create(&comment).Error; err != nil {
		return comment, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return comment, nil
}

func (s *store) UpdateComments(ctx context.Context, comment core.Comments) (core.Comments, error) {
	if err := s.DB.WithContext(ctx).Save(&comment).Error; err != nil {
		return comment, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return comment, nil
}

func (s *store) DeleteComments(ctx context.Context, id int) error {
	deletedComment := core.Comments{
		Is_deleted: true,
		Deleted_at: time.Now(),
	}
	if err := s.DB.WithContext(ctx).Save(&deletedComment).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return nil
}
