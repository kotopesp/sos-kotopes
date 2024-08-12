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
		GetFavouritePosts(ctx context.Context, userID int, params GetAllPostsParams) ([]Post, int, error)
		GetPostFavouriteByPostAndUserID(ctx context.Context, postID, userID int) (PostFavourite, error) // needed for if post.UserID == userID in DeleteFromFavourites
		AddToFavourites(ctx context.Context, postFavourite PostFavourite) (Post, error)
		DeleteFromFavourites(ctx context.Context, postID, userID int) error
	}

	PostFavouriteService interface {
		GetFavouritePosts(ctx context.Context, userID int, params GetAllPostsParams) ([]PostDetails, int, error)
		AddToFavourites(ctx context.Context, postFavourite PostFavourite) (PostDetails, error)
		DeleteFromFavourites(ctx context.Context, post PostFavourite) error
	}
)

func (PostFavourite) TableName() string {
	return "favourite_posts"
}
