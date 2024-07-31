package core

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type (
	Post struct {
		ID        int       `gorm:"column:id"`
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

	PostStore interface {
		GetPostByID(ctx context.Context, id int) (data Post, err error)
	}
)

var (
	ErrPostNotFound   = errors.New("post not found")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

func (Post) TableName() string {
	return "posts"
}
