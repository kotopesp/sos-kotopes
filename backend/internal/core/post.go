package core

import (
	"context"
	"time"
	"mime/multipart"
)

type (

	Post struct {
		ID        int       `gorm:"column:id;primaryKey"`
		Title     string    `gorm:"column:title"`
		Content   string    `gorm:"column:content"`
		AuthorID  int       `gorm:"column:author_id"`
		AnimalID  int       `gorm:"column:animal_id"`
		IsDeleted bool      `gorm:"column:is_deleted"`
		DeletedAt time.Time `gorm:"column:deleted_at"`
		CreatedAt time.Time `gorm:"column:created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at"`
		Photo     []byte    `gorm:"column:photo"`
	}

	PostDetails struct {
		Post 	 Post
		Animal   Animal
		Username string
	}

	UpdateRequestBodyPost struct {
		Title   string `gorm:"column:title"`
		Content string `gorm:"column:content"`
		Photo   []byte `gorm:"column:photo"`
	}

	GetAllPostsParams struct {
		Limit      *int
		Offset     *int
		Status     *string
		AnimalType *string
		Gender     *string
		Color      *string
		Location   *string
	}

	PostStore interface {
		GetAllPosts(ctx context.Context, params GetAllPostsParams) (data []Post, total int, err error)
		GetPostByID(ctx context.Context, id int) (data Post, err error)
		CreatePost(ctx context.Context, post Post) (data Post, err error)
		UpdatePost(ctx context.Context, post Post) (data Post, err error)
		DeletePost(ctx context.Context, id int) (err error)
	}

	PostService interface {
		GetAllPosts(ctx context.Context, params GetAllPostsParams) (data []PostDetails, total int, err error)
		GetPostByID(ctx context.Context, id int) (PostDetails, error)
		CreatePost(ctx context.Context, postDetails PostDetails, fileHeader *multipart.FileHeader) (PostDetails, error)
		UpdatePost(ctx context.Context, postDetails PostDetails) (PostDetails, error)
		DeletePost(ctx context.Context, id int) (err error)

		PostFavouriteService
	}
	
	
)

func (Post) TableName() string {
	return "posts"
}
