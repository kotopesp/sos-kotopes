package core

import (
	"context"
	"time"
)

type (
	PostFavourite struct {
		ID        int       `gorm:"column:id"`         // Unique identifier for the favourite record
		UserID    int       `gorm:"column:user_id"`    // ID of the user who favourited the post
		PostID    int       `gorm:"column:post_id"`    // ID of the post that was favourited (the post contains the author_id)
		CreatedAt time.Time `gorm:"column:created_at"` // Timestamp when the post was favourited
	}

	PostFavouriteStore interface {
		GetFavouritePosts(ctx context.Context, userID, limit, offset int) ([]Post, int, error)
		AddToFavourites(ctx context.Context, postFavourite PostFavourite) error
		DeleteFromFavourites(ctx context.Context, postID, userID int) error
	}

	PostFavouriteService interface {
		GetFavouritePosts(ctx context.Context, userID int) ([]PostDetails, int, error)
		AddToFavourites(ctx context.Context, postFavourite PostFavourite) error
		DeleteFromFavourites(ctx context.Context, postID, userID int) error
	}
)

func (PostFavourite) TableName() string {
	return "favourite_posts"
}
