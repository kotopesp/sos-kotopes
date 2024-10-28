package core

import (
	"context"
)

type (
	Photo struct {
		ID            int    `gorm:"column:id;primaryKey"`
		PostID        int    `gorm:"column:post_id"`
		URL           string `gorm:"column:url"`
		FileExtension string `gorm:"column:file_extension"`
		Photo         []byte `gorm:"column:photo"`
	}

	PhotoStore interface {
		GetPhotosPost(ctx context.Context, postID int) ([]Photo, error)
		GetPhotosPostByPhotoID(ctx context.Context, postID, photoID int) (Photo, error)
		AddPhotoPost(ctx context.Context, photo Photo) (Photo, error)
		AddPhotosPost(ctx context.Context, postID int, photos []Photo) ([]Photo, error)
		UpdatePhotosPost(ctx context.Context, photos []Photo) ([]Photo, error)
	}

	PostPhotoService interface {
		GetPhotosPost(ctx context.Context, postID int) ([]Photo, error)
		GetPhotosPostByPhotoID(ctx context.Context, postID, photoID int) (Photo, error)
		AddPhotosPost(ctx context.Context, postID int, photos []Photo) ([]Photo, error)
	}
)

func (Photo) TableName() string {
	return "photo_posts"
}
