package post

import (
	"context"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"

	"fmt"
)

func (s *service) GetPhotosPost(ctx context.Context, postID int) ([]core.Photo, error) {
    var photos []core.Photo

    photos, err := s.photoStore.GetPhotosPost(ctx, postID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, err
	}

    for i := range photos {
        photos[i].URL = generatePhotoURL(postID, photos[i].ID)
    }

    return photos, nil
}

func (s *service) GetPhotosPostByPhotoID(ctx context.Context, postID, photoID int) (core.Photo, error) {
	photo, err := s.photoStore.GetPhotosPostByPhotoID(ctx, postID, photoID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.Photo{}, err
	}

	photo.URL = generatePhotoURL(postID, photo.ID)

	return photo, nil
}

func generatePhotoURL(postID, photoID int) string {
    return fmt.Sprintf("posts/%d/photos/%d", postID, photoID)
}

func (s *service) AddPhotosPost(ctx context.Context, postID int, photos []core.Photo) ([]core.Photo, error) {
	photo, err := s.photoStore.AddPhotosPost(ctx, postID, photos)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, err
	}

	for i := range photo {
		photo[i].URL = generatePhotoURL(postID, photo[i].ID)
	}

	return photo, nil
}
