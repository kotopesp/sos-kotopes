package post

import (
	"context"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

func (s *postService) GetFavouritePosts(ctx context.Context, userID, limit, offset int) ([]core.PostDetails, int, error) {
	posts, total, err := s.postFavouriteStore.GetFavouritePosts(ctx, userID, limit, offset)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, 0, err
	}

	postDetails := make([]core.PostDetails, total)
	for i, post := range posts {
		animal, err := s.animalStore.GetAnimalByID(ctx, post.AnimalID)
		if err != nil {
			logger.Log().Error(ctx, err.Error())
			return nil, 0, core.ErrAnimalNotFound
		}
		user, err := s.userStore.GetUserByID(ctx, post.AuthorID)
		if err != nil {
			logger.Log().Error(ctx, err.Error())
			return nil, 0, core.ErrUserNotFound
		}

		postDetails[i] = core.PostDetails{
			Post: post,
			Animal: animal,
			Username: user.Username,
		}
	}

	return postDetails, total, nil
}

func (s *postService) GetFavouritePostByID(ctx context.Context, userID, postID int) (core.PostDetails, error) {
	post, err := s.postFavouriteStore.GetFavouritePostByID(ctx, userID, postID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.PostDetails{}, err
	}

	animal, err := s.animalStore.GetAnimalByID(ctx, post.AnimalID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.PostDetails{}, err
	}

	user, err := s.userStore.GetUserByID(ctx, post.AuthorID)
		if err != nil {
			logger.Log().Error(ctx, err.Error())
			return core.PostDetails{}, core.ErrUserNotFound
		}

	postDetails := core.PostDetails{
		Post: post,
		Animal: animal,
		Username: user.Username,
	}

	return postDetails, nil
}

func (s *postService) AddToFavourites(ctx context.Context, postFavourite core.PostFavourite) (error) {
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
