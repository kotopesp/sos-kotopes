package photostore

import (
	"context"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
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

func (s *store) GetPhotosPostByPhotoID(ctx context.Context, postID, photoID int) (core.Photo, error) {
	var photo core.Photo

	if err := s.DB.WithContext(ctx).Where("post_id = ? AND id = ?", postID, photoID).First(&photo).Error; err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.Photo{}, err
	}

	return photo, nil
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

func (s *store) updatePhotoPost(ctx context.Context, photo core.Photo) (core.Photo, error) {
	var updatedPhoto core.Photo

	if err := s.DB.WithContext(ctx).
		Where("id = ?", photo.ID).
		Save(&photo).
		First(&updatedPhoto, photo.ID).
		Error; err != nil {
			
		logger.Log().Error(ctx, err.Error())
		return core.Photo{}, err
	}

	return updatedPhoto, nil
}

func (s *store) UpdatePhotosPost(ctx context.Context, photos []core.Photo) ([]core.Photo, error) {
	var updatedPhotos []core.Photo

	photosPost, err := s.GetPhotosPost(ctx, photos[0].PostID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, err
	}

	existingCount := len(photosPost)
	newCount := len(photos)

	for i := 0; i < existingCount && i < newCount; i++ {
		p := photos[i]
		updatedPhoto, err := s.updatePhotoPost(ctx, p)
		if err != nil {
			logger.Log().Error(ctx, err.Error())
			return nil, err
		}
		updatedPhotos = append(updatedPhotos, updatedPhoto)
	}

	if newCount > existingCount {
		newPhotos := photos[existingCount:]
		createdPhotos, err := s.AddPhotosPost(ctx, photos[0].PostID, newPhotos)
		if err != nil {
			logger.Log().Error(ctx, err.Error())
			return nil, err
		}
		updatedPhotos = append(updatedPhotos, createdPhotos...)
	}

	if newCount < existingCount {
		photosToDelete := photosPost[newCount:]
		err := s.deletePhotosPost(ctx, photosToDelete)
		if err != nil {
			logger.Log().Error(ctx, err.Error())
			return nil, err
		}
	}

	return updatedPhotos, nil
}

func (s *store) deletePhotosPost(ctx context.Context, photos []core.Photo) error {
	var idsToDelete []int

	for _, photo := range photos {
		idsToDelete = append(idsToDelete, photo.ID)
	}

	if err := s.DB.WithContext(ctx).Where("id IN (?)", idsToDelete).Delete(&core.Photo{}).Error; err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}

	return nil
}
