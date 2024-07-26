package comments

import (
	"context"

	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
)

type commentsService struct {
	CommentsStore core.CommentsStore
}

func NewCommentsService(store core.CommentsStore) core.CommentsService {
	return &commentsService{
		CommentsStore: store,
	}
}

func (s *commentsService) GetCommentsByPostID(ctx context.Context, params core.GetAllParamsComments, post_id int) ([]core.Comments, error) {

}

func (s *commentsService) CreateComment(ctx context.Context, comment core.Comments, post_id int) (core.Comments, error) {

}
func (s *commentsService) UpdateComments(ctx context.Context, comments core.Comments) (core.Comments, error) {

}
func (s *commentsService) DeleteComments(ctx context.Context, id int) error {

}
