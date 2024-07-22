package postfavouriteservice

import (
	"context"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
)

type postFavoriteService struct {
	PostFavoriteStore core.PostFavoriteStore
}

func NewPostFavoriteService(store core.PostFavoriteStore) core.PostFavoriteService {
    return &postFavoriteService{
        PostFavoriteStore: store,
    }
}

func (s *postFavoriteService) GetFavoritePosts(ctx context.Context, userID int, params core.GetAllPostsParams) ([]core.Post, int, error) {
    return s.PostFavoriteStore.GetFavoritePosts(ctx, userID, params)
}

func (s *postFavoriteService) GetFavoritePostByID(ctx context.Context, userID, postID int) (core.Post, error) {
    return s.PostFavoriteStore.GetFavoritePostByID(ctx, userID, postID)
}

func (s *postFavoriteService) AddToFavorites(ctx context.Context, postFavorite core.PostFavorite) (core.PostFavorite, error) {
    return s.PostFavoriteStore.AddToFavorites(ctx, postFavorite)
}

func (s *postFavoriteService) DeleteFromFavorites(ctx context.Context, postID, userID int) error {
    return s.PostFavoriteStore.DeleteFromFavorites(ctx, postID, userID)
}