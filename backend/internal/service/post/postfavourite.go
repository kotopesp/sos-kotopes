package post

import (
	"context"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

func (s *service) GetFavouritePosts(ctx context.Context, userID int, params core.GetAllPostsParams) ([]core.PostDetails, int, error) {
	posts, total, err := s.postFavouriteStore.GetFavouritePosts(ctx, userID, params)
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

func (s *service) AddToFavourites(ctx context.Context, postFavourite core.PostFavourite) (core.PostDetails, error) {
	post, err := s.postFavouriteStore.AddToFavourites(ctx, postFavourite)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.PostDetails{}, err
	}

	postDetails, err := s.BuildPostDetails(ctx, post)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.PostDetails{}, err
	}
	
	return postDetails, nil
}

func (s *service) DeleteFromFavourites(ctx context.Context, post core.PostFavourite) error {
	dbFavourite, err := s.postFavouriteStore.GetPostFavouriteByPostAndUserID(ctx, post.PostID, post.UserID)
	if err != nil {
		return err
	}

	if dbFavourite.UserID != post.UserID {
		return core.ErrPostAuthorIDMismatch
	}

	err = s.postFavouriteStore.DeleteFromFavourites(ctx, post.PostID, post.UserID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}

	return nil
}

