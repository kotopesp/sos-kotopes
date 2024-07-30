package core

import (
	"context"
	"time"
)

type (
	Post struct {
		ID        int       `gorm:"column:id"`
		Title     string    `gorm:"column:title"`
		Content   string    `gorm:"column:content"`
		UserID    int       `gorm:"column:author_id"`
		IsDeleted bool      `gorm:"column:is_deleted"`
		DeletedAt time.Time `gorm:"column:deleted_at"`
		CreatedAt time.Time `gorm:"column:created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at"`
		AnimalID  int       `gorm:"column:animal_id"`
	}
	PostDetails struct {
		ID        int       `gorm:"column:id"`
		Title     string    `gorm:"column:title"`
		Content   string    `gorm:"column:content"`
		Username  string    `gorm:"-"`
		CreatedAt time.Time `gorm:"column:created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at"`
		AnimalID  int       `gorm:"column:animal_id"`
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

func (p *Post) ToPostDetails(name string) PostDetails {
	if p == nil {
		return PostDetails{}
	}
	return PostDetails{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		Username:  name,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		AnimalID:  p.AnimalID,
	}
}
