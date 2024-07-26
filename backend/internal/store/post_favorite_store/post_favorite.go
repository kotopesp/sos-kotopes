package postfavouritestore

import (
	"context"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/internal/store/errors"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
)

type store struct {
	*postgres.Postgres
}

func NewFavoritePostStore(pg *postgres.Postgres) core.PostFavoriteStore {
	return &store{pg}
}

func (s *store) GetFavoritePosts(ctx context.Context, userID int, params core.GetAllPostsParams) ([]core.Post, int, error) {
	var posts []core.Post
	query := s.DB.WithContext(ctx).Table("posts").
		Joins("JOIN favourite_posts ON posts.id = favourite_posts.post_id").
		Where("favourite_posts.user_id = ?", userID)

	if params.SortBy != nil && params.SortOrder != nil {
		query = query.Order(*params.SortBy + " " + *params.SortOrder)
	}

	if params.SearchTerm != nil {
		query = query.Where("posts.title ILIKE ?", "%"+*params.SearchTerm+"%")
	}

	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}

	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Select("posts.*").Find(&posts).Error; err != nil {
		return nil, 0, err
	}
	return posts, int(total), nil
}

func (s *store) GetFavoritePostByID(ctx context.Context, userID, postID int) (core.Post, error) {
	var post core.Post
	if err := s.DB.WithContext(ctx).Table("posts").
		Joins("JOIN favourite_posts ON posts.id = favourite_posts.post_id").
		Where("favourite_posts.user_id = ? AND favourite_posts.post_id = ?", userID, postID).
		Select("posts.*").First(&post).Error; err != nil {
		return core.Post{}, err
	}
	return post, nil
}

func (s *store) AddToFavorites(ctx context.Context, postFavourite core.PostFavorite) (core.PostFavorite, error) {
	var existing core.PostFavorite
	if err := s.DB.WithContext(ctx).
		Where("post_id = ? AND user_id = ?", postFavourite.PostID, postFavourite.UserID).
		First(&existing).Error; err == nil {
		return core.PostFavorite{}, errors.ErrPostAlreadyInFavorites
	}

	postFavourite.CreatedAt = time.Now()

	if err := s.DB.WithContext(ctx).Create(&postFavourite).Error; err != nil {
		return core.PostFavorite{}, err
	}
	return postFavourite, nil
}

func (s *store) DeleteFromFavorites(ctx context.Context, postID, userID int) error {
	if err := s.DB.WithContext(ctx).Where("post_id = ? AND user_id = ?", postID, userID).Delete(&core.PostFavorite{}).Error; err != nil {
		return err
	}
	return nil
}
