package photostore

import (
	"context"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
	"gorm.io/gorm"
	"errors"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.PhotoStore {
	return &store{pg}
}

func (s *store) GetPhotosPost(ctx context.Context, postID int) ([]core.Photo, error) {
	var photos []core.Photo

	if err := s.DB.WithContext(ctx).Where("post_id = ?", postID).Find(&photos).Error; err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, err
	}

	return photos, nil
}

func (s *store) AddPhotosPost(ctx context.Context, postID int, photos []core.Photo) ([]core.Photo, error) {
	var createdPhotos []core.Photo

	for _, photo := range photos {
		photo.PostID = postID
		createdPhoto, err := s.AddPhotoPost(ctx, photo)
		if err != nil {
			return nil, err
		}
		createdPhotos = append(createdPhotos, createdPhoto)
	}

	return createdPhotos, nil
}

func (s *store) AddPhotoPost(ctx context.Context, photo core.Photo) (core.Photo, error) {
	var createdPhoto core.Photo

	if err := s.DB.WithContext(ctx).Create(&photo).First(&createdPhoto, photo.ID).Error; err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.Photo{}, err
	}

	return createdPhoto, nil
}

func (s *store) UpdatePhotoPost(ctx context.Context, photo core.Photo) (core.Photo, error) {
	panic("implement me")
}

func (s *store) UpdatePhotosPost(ctx context.Context, photos []core.Photo) ([]core.Photo, error) {
	var updatedPhotos []core.Photo

    err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        for _, photo := range photos {

			p := photo

            var updatedPhoto core.Photo

            if err := tx.
                Where("id = ?", p.ID).
                Save(&p).
                First(&updatedPhoto, p.ID).
                Error; err != nil {

                if errors.Is(err, gorm.ErrRecordNotFound) {
                    logger.Log().Error(ctx, core.ErrPhotoNotFound.Error())
                    return core.ErrPhotoNotFound
                }

                logger.Log().Error(ctx, err.Error())
                return err
            }

            updatedPhotos = append(updatedPhotos, updatedPhoto)
        }
        return nil
    })

    if err != nil {
        return nil, err
    }

    return updatedPhotos, nil
}