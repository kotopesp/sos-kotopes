package postresponseservice

import (
	"context"

	"github.com/kotopesp/sos-kotopes/internal/core"
)

type (
	responseService struct {
		responseStore core.PostResponseStore
	}
)

func New(store core.PostResponseStore) core.PostResponseService {
	return &responseService{store}
}

func (s *responseService) CreatePostResponse(ctx context.Context, response core.PostResponse) (core.PostResponse, error) {
	return s.responseStore.CreatePostResponse(ctx, response)
}

func (s *responseService) GetResponsesByPostID(ctx context.Context, postID int) ([]core.PostResponse, error) {
	return s.responseStore.GetResponsesByPostID(ctx, postID)
}

func (s *responseService) UpdatePostResponse(ctx context.Context, response core.PostResponse) (core.PostResponse, error) {
	return s.responseStore.UpdatePostResponse(ctx, response)
}

func (s *responseService) DeletePostResponse(ctx context.Context, id int) error {
	return s.responseStore.DeletePostResponse(ctx, id)
}
