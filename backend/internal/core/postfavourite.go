package core

import (
	"context"
	"time"
)
type (
	PostFavourite struct {
		ID        int       `gorm:"column:id"`
		UserID  int       `gorm:"column:user_id"`
		PostID    int       `gorm:"column:post_id"`
		CreatedAt time.Time `gorm:"column:created_at"`
	}

	PostFavouriteStore interface {
		GetFavouritePosts(ctx context.Context, userID, limit, offset int) (data []Post, total int, err error)
		GetFavouritePostByID(ctx context.Context, userID, postID int) (data Post, err error)
		AddToFavourites(ctx context.Context, postFavourite PostFavourite) (err error)
		DeleteFromFavourites(ctx context.Context, postID, userID int) (err error)
	}

	PostFavouriteService interface {
		GetFavouritePosts(ctx context.Context, userID, limit, offset int) (data []PostDetails, total int, err error)
		GetFavouritePostByID(ctx context.Context, userID, postID int) (data PostDetails, err error)
		AddToFavourites(ctx context.Context, postFavourite PostFavourite) (err error)
		DeleteFromFavourites(ctx context.Context, postID, userID int) (err error)
	}
)

func (PostFavourite) TableName() string {
	return "favourite_posts"
}
