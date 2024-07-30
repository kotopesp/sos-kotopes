package postservice

import (
	"context"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

func (s *postService) GetFavouritePosts(ctx context.Context, userID, limit, offset int) ([]core.Post, int, error) {
	posts, total, err := s.PostFavouriteStore.GetFavouritePosts(ctx, userID, limit, offset)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, 0, core.ErrPostNotFound
	}

	return posts, total, nil
}

func (s *postService) GetFavouritePostByID(ctx context.Context, userID, postID int) (core.Post, core.Animal, error) {
	post, err := s.PostFavouriteStore.GetFavouritePostByID(ctx, userID, postID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.Post{}, core.Animal{}, core.ErrPostNotFound
	}

	animal, err := s.AnimalStore.GetAnimalByID(ctx, post.AnimalID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.Post{}, core.Animal{}, core.ErrAnimalNotFound
	}

	return post, animal, nil
}

func (s *postService) AddToFavourites(ctx context.Context, postFavorite core.PostFavourite) (error) {
    return s.PostFavouriteStore.AddToFavourites(ctx, postFavorite)
}

func (s *postService) DeleteFromFavourites(ctx context.Context, postID, userID int) error {
    err := s.PostFavouriteStore.DeleteFromFavourites(ctx, postID, userID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}
	return nil
}
