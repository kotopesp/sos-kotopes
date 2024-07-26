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
	comments, err := s.CommentsStore.GetCommentsByPostID(ctx, params, post_id)
	if err != nil {
		//TODO обработать
	}
	return comments, nil
}

func (s *commentsService) CreateComment(ctx context.Context, comment core.Comments, post_id int) (core.Comments, error) {
	return s.CommentsStore.CreateComment(ctx, comment, post_id)
}
func (s *commentsService) UpdateComments(ctx context.Context, comments core.Comments) (core.Comments, error) {
	return s.CommentsStore.UpdateComments(ctx, comments)
}
func (s *commentsService) DeleteComments(ctx context.Context, id int) error {
	return s.CommentsStore.DeleteComments(ctx, id)
}
