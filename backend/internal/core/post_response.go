package core

import (
	"context"
	"time"
)

type (
	PostResponse struct {
		ID        int       `gorm:"column:id"`
		PostID    int       `gorm:"column:post_id"`
		AuthorID  int       `gorm:"column:author_id"`
		Content   string    `gorm:"column:content"`
		IsDeleted bool      `gorm:"column:is_deleded"`
		DeletedAt time.Time `gorm:"column:deleted_at"`
		CreatedAt time.Time `gorm:"column:created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at"`
	}

	PostResponseStore interface {
		CreatePostResponse(ctx context.Context, response PostResponse) (PostResponse, error)
		GetResponsesByPostID(ctx context.Context, postID int) ([]PostResponse, error)
		UpdatePostResponse(ctx context.Context, response PostResponse) (PostResponse, error)
		DeletePostResponse(ctx context.Context, id int) error
	}

	PostResponseService interface {
		CreatePostResponse(ctx context.Context, response PostResponse) (PostResponse, error)
		GetResponsesByPostID(ctx context.Context, postID int) ([]PostResponse, error)
		UpdatePostResponse(ctx context.Context, response PostResponse) (PostResponse, error)
		DeletePostResponse(ctx context.Context, id int) error
	}
)

func (PostResponse) TableName() string {
	return "post_response"
}
