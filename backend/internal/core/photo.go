package core

import (
	"context"
)

type (
	Photo struct {
		ID          int    `gorm:"column:id;primaryKey"`
		PostID      int    `gorm:"column:post_id"`
		Photo 		[]byte `gorm:"column:photo"`
	}

	PhotoStore interface {
		GetPhotosPost(ctx context.Context, postID int) ([]Photo, error)
		AddPhotoPost(ctx context.Context, photo Photo) (Photo, error)
		AddPhotosPost(ctx context.Context, postID int, photos []Photo) ([]Photo, error)
		UpdatePhotosPost(ctx context.Context, photos []Photo) ([]Photo, error)
	}
)

func (Photo) TableName() string {
	return "photo_posts"
}
