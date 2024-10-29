package post

import (
	"context"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

// ToCorePostDetails creates a core.PostDetails object from a core.Post, core.Animal, and username string.
func ToCorePostDetails(post core.Post, animal core.Animal, userName string, photos []core.Photo) core.PostDetails {
	return core.PostDetails{
		Post:     post,
		Animal:   animal,
		Username: userName,
		Photos:   photos,
	}
}

// BuildPostDetailsList constructs a list of core.PostDetails from a list of core.Post objects.
// It fetches the associated animal and user information for each post.
func (s *service) BuildPostDetailsList(ctx context.Context, posts []core.Post, total int) ([]core.PostDetails, error) {
	postDetailsList := make([]core.PostDetails, 0, total)

	// Iterate through each post to build the post details
	for _, post := range posts {
		postDetails, err := s.BuildPostDetails(ctx, post)
		if err != nil {
			logger.Log().Error(ctx, err.Error())
			return nil, err
		}

		postDetailsList = append(postDetailsList, postDetails)
	}

	return postDetailsList, nil
}

// BuildPostDetails constructs a core.PostDetails object from a core.Post object.
// It fetches the associated animal and user information for the post.
func (s *service) BuildPostDetails(ctx context.Context, post core.Post) (core.PostDetails, error) {
	animal, err := s.animalStore.GetAnimalByID(ctx, post.AnimalID) // Fetch the animal details
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.PostDetails{}, err
	}

	user, err := s.userStore.GetUserByID(ctx, post.AuthorID) // Fetch the user details
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.PostDetails{}, err
	}

	photos, err := s.GetPhotosPost(ctx, post.ID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.PostDetails{}, err
	}

	postDetails := ToCorePostDetails(post, animal, user.Username, photos)

	return postDetails, nil
}

// UpdateRequestPost updates the post details based on core.UpdateRequestBodyPost
func UpdateRequestPost(postDetails core.PostDetails, updatePost core.UpdateRequestPost) core.PostDetails {
	if updatePost.Title != nil {
		postDetails.Post.Title = *updatePost.Title
	}

	if updatePost.Content != nil {
		postDetails.Post.Content = *updatePost.Content
	}

	if updatePost.AnimalType != nil {
		postDetails.Animal.AnimalType = *updatePost.AnimalType
	}

	if updatePost.Age != nil {
		postDetails.Animal.Age = *updatePost.Age
	}

	if updatePost.Color != nil {
		postDetails.Animal.Color = *updatePost.Color
	}

	if updatePost.Gender != nil {
		postDetails.Animal.Gender = *updatePost.Gender
	}

	if updatePost.Description != nil {
		postDetails.Animal.Description = *updatePost.Description
	}

	if updatePost.Status != nil {
		postDetails.Animal.Status = *updatePost.Status
	}

	return postDetails
}

func UpdateRequestPhotos(photos []core.Photo, updatePhotos core.UpdateRequestPhotos) []core.Photo {
	if updatePhotos.Photos == nil {
		return photos
	}

	updatedPhotos := *updatePhotos.Photos

	minLen := len(photos)
	if len(updatedPhotos) < minLen {
		minLen = len(updatedPhotos)
	}
	for i := 0; i < minLen; i++ {
		photos[i].Photo = updatedPhotos[i]
	}

	if len(updatedPhotos) > len(photos) {
		for i := len(photos); i < len(updatedPhotos); i++ {
			photos = append(photos, core.Photo{
				Photo: updatedPhotos[i],
				PostID: photos[0].PostID,
			})
		}
	}

	if len(updatedPhotos) < len(photos) {
		photos = photos[:len(updatedPhotos)]
	}

	return photos
}
