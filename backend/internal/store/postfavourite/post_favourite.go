package postfavourite

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.PostFavouriteStore {
	return &store{pg}
}

// GetFavouritePosts retrieves the favourite posts of a user from the database based on the GetAllPostsParams
// TODO: add GetAllPostsParams
func (s *store) GetFavouritePosts(ctx context.Context, userID, limit, offset int) ([]core.Post, int, error) {
	var posts []core.Post

	query := s.DB.WithContext(ctx).Model(&core.Post{}).
		Joins("JOIN favourite_posts ON posts.id = favourite_posts.post_id").
		Where("favourite_posts.user_id = ?", userID).
		Offset(offset).
		Limit(limit)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return nil, 0, core.ErrPostNotFound
		}
		logger.Log().Error(ctx, err.Error())
		return nil, 0, err
	}

	if err := query.Select("posts.*").Find(&posts).Error; err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return nil, 0, core.ErrPostNotFound
		}
		logger.Log().Error(ctx, err.Error())
		return nil, 0, err
	}

	return posts, int(total), nil
}

// AddToFavourites adds a post to the user's favourites in the database
// TODO: return added post
func (s *store) AddToFavourites(ctx context.Context, postFavourite core.PostFavourite) error {
	// Check if the post exists
	var post core.Post
	if err := s.DB.WithContext(ctx).Where("id = ?", postFavourite.PostID).First(&post).Error; err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.ErrPostNotFound
		}
		logger.Log().Error(ctx, err.Error())
		return err
	}

	// Check if the post is already in the user's favourites
	var existing core.PostFavourite
	if err := s.DB.WithContext(ctx).
		Where("post_id = ? AND user_id = ?", postFavourite.PostID, postFavourite.UserID).
		First(&existing).Error; err == nil {
		return core.ErrPostAlreadyInFavourites
	}

	postFavourite.CreatedAt = time.Now().UTC()

	logger.Log().Debug(ctx, fmt.Sprintf("%v", postFavourite))
	if err := s.DB.WithContext(ctx).Create(&postFavourite).Error; err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}

	return nil
}

// DeleteFromFavourites removes a post from the user's favourites in the database
func (s *store) DeleteFromFavourites(ctx context.Context, postID, userID int) error {
	if err := s.DB.WithContext(ctx).Where("post_id = ? AND user_id = ?", postID, userID).Delete(&core.PostFavourite{}).Error; err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.ErrPostNotFound
		}
		logger.Log().Error(ctx, err.Error())
		return err
	}

	return nil
}
