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

	UpdateRequestBodyPost struct {
		Title   string `gorm:"column:title"`
		Content string `gorm:"column:content"`
		Photo   []byte `gorm:"column:photo"`
	}

	PostStore interface {
		GetAllPosts(ctx context.Context, limit, offset int) (data []Post, total int, err error)
		GetPostByID(ctx context.Context, id int) (data Post, err error)
		CreatePost(ctx context.Context, post Post) (err error)
		UpdatePost(ctx context.Context, post Post) (err error)
		DeletePost(ctx context.Context, id int) (err error)

		GetAuthorUsernameByID(ctx context.Context, authorID int) (string, error)
	}
	PostService interface {
		GetAllPosts(ctx context.Context, limit, offset int) (data []Post, total int, err error)
		GetPostByID(ctx context.Context, id int) (Post, Animal, error)
		CreatePost(ctx context.Context, post Post, fileHeader *multipart.FileHeader, animal Animal) (err error)
		UpdatePost(ctx context.Context, post Post, animal Animal) (err error)
		DeletePost(ctx context.Context, id int) (err error)

		GetFavouritePosts(ctx context.Context, userID, limit, offset int) (data []Post, total int, err error)
		GetFavouritePostByID(ctx context.Context, userID, postID int) (data Post, animal Animal, err error)
		AddToFavourites(ctx context.Context, postFavourite PostFavourite) (err error)
		DeleteFromFavourites(ctx context.Context, postID, userID int) (err error)

		GetAuthorUsernameByID(ctx context.Context, authorID int) (string, error)

		GetAnimalByID(ctx context.Context, id int) (animal Animal, err error)
	}	
)

func (Post) TableName() string {
	return "posts"
}
