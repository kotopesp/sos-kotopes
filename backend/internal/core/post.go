package core

import (
	"context"
	"time"
)

type (
	Post struct {
		ID        int       `gorm:"column:id"`
		Title     string    `gorm:"column:title"`
		Body      string    `gorm:"column:body"`
		UserID    int       `gorm:"column:user_id"`
		CreatedAt time.Time `gorm:"column:created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at"`
		AnimalID  int       `gorm:"column:animal_id"`
		Photo     []byte    `gorm:"column:photo"`
	}

	PostFavorite struct {
		ID        int       `gorm:"column:id"`
		UserID    int       `gorm:"column:user_id"`
		PostID    int       `gorm:"column:post_id"`
		CreatedAt time.Time `gorm:"column:created_at"`
	}

	PostStore interface {
		GetAllPosts(ctx context.Context, params GetAllPostsParams) (data []Post, total int, err error)
		GetPostByID(ctx context.Context, ID int) (data Post, err error)
		CreatePost(ctx context.Context, post Post) (data Post, err error)
		UpdatePost(ctx context.Context, post Post) (data Post, err error)
		DeletePost(ctx context.Context, ID int) (err error)
	}
	PostService interface {
		GetAllPosts(ctx context.Context, params GetAllPostsParams) (data []Post, total int, err error)
		GetPostByID(ctx context.Context, ID int) (data Post, err error)
		CreatePost(ctx context.Context, post Post) (data Post, err error)
		UpdatePost(ctx context.Context, post Post) (data Post, err error)
		DeletePost(ctx context.Context, ID int) (err error)
	}

	PostFavoriteStore interface {
		GetFavoritePosts(ctx context.Context, userID int, params GetAllPostsParams) (data []Post, total int, err error)
		GetFavoritePostByID(ctx context.Context, userID, postID int) (data Post, err error)
		AddToFavorites(ctx context.Context, postFavourite PostFavorite) (data PostFavorite, err error)
		DeleteFromFavorites(ctx context.Context, postID, userID int) (err error)
	}

	PostFavoriteService interface {
		GetFavoritePosts(ctx context.Context, userID int, params GetAllPostsParams) (data []Post, total int, err error)
		GetFavoritePostByID(ctx context.Context, userID, postID int) (data Post, err error)
		AddToFavorites(ctx context.Context, postFavourite PostFavorite) (data PostFavorite, err error)
		DeleteFromFavorites(ctx context.Context, postID, userID int) (err error)
	}

	GetAllPostsParams struct {
		SortBy     *string
		SortOrder  *string
		SearchTerm *string
		Limit      *int
		Offset     *int
	}
)

func (Post) TableName() string {
	return "posts"
}

func (PostFavorite) TableName() string {
	return "favourite_posts"
}