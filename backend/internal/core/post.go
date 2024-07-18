package core

import (
	"context"
	"time"
)

type (
	Post struct {
		ID        int       `gorm:"column:id" json:"id"`
		Title     string    `gorm:"column:title" json:"title"`
		Body      string    `gorm:"column:body" json:"body"`
		UserID    int       `gorm:"column:user_id" json:"user_id"`
		CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
		AnimalID  int       `gorm:"column:animal_id" json:"animal_id"`
	}

	PostStore interface {
		GetAll(ctx context.Context, params GetAllPostsParams) (data []Post, err error)
		GetByID(ctx context.Context, id int) (data Post, err error)
		Create(ctx context.Context, post Post) (data Post, err error)
		Update(ctx context.Context, post Post) (data Post, err error)
		Delete(ctx context.Context, id int) (err error)
	}		
	PostService interface {
		GetAll(ctx context.Context, params GetAllPostsParams) (data []Post, total int, err error)
		GetByID(ctx context.Context, id int) (data Post, err error)
		Create(ctx context.Context, post Post) (data Post, err error)
		Update(ctx context.Context, post Post) (data Post, err error)
		Delete(ctx context.Context, id int) (err error)
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