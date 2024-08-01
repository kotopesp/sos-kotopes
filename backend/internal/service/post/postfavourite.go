package post

import (
	"context"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

func (s *postService) GetFavouritePosts(ctx context.Context, userID int) ([]core.PostDetails, int, error) {
	// TODO: add params processing
	limit := 0
	offset := 0

	posts, total, err := s.postFavouriteStore.GetFavouritePosts(ctx, userID, limit, offset)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, 0, err
	}

	postDetails, err := s.BuildPostDetailsList(ctx, posts, total)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, 0, err
	}

	return postDetails, total, nil
}

func (s *postService) AddToFavourites(ctx context.Context, postFavourite core.PostFavourite) error {
	// TODO: return PostDetails
    err := s.postFavouriteStore.AddToFavourites(ctx, postFavourite)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}
	
	return nil
}

func (s *postService) DeleteFromFavourites(ctx context.Context, postID, userID int) error {
    err := s.postFavouriteStore.DeleteFromFavourites(ctx, postID, userID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}

	return nil
}
