package post

import (
	"context"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

// ToCorePostDetails creates a core.PostDetails object from a core.Post, core.Animal, and username string.
func ToCorePostDetails(post core.Post, animal core.Animal, userName string) core.PostDetails {
	return core.PostDetails{
		Post:     post,
		Animal:   animal,
		Username: userName,
	}
}

// BuildPostDetailsList constructs a list of core.PostDetails from a list of core.Post objects.
// It fetches the associated animal and user information for each post.
func (s *service) BuildPostDetailsList(ctx context.Context, posts []core.Post, total int) ([]core.PostDetails, error) {
	postDetails := make([]core.PostDetails, total)

	// Iterate through each post to build the post details
	for i, post := range posts {
		animal, err := s.animalStore.GetAnimalByID(ctx, post.AnimalID) // Fetch the animal details
		if err != nil {
			logger.Log().Error(ctx, err.Error())
			return nil, err
		}
		user, err := s.userStore.GetUserByID(ctx, post.AuthorID) // Fetch the user details
		if err != nil {
			logger.Log().Error(ctx, err.Error())
			return nil, err
		}

		postDetails[i] = ToCorePostDetails(post, animal, user.Username)
	}

	return postDetails, nil
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

	postDetails := ToCorePostDetails(post, animal, user.Username)

	return postDetails, nil
}

// FuncUpdateRequestBodyPost updates the post details based on UpdateRequestBodyPost
func FuncUpdateRequestBodyPost(postDetails core.PostDetails, updatePost core.UpdateRequestBodyPost) core.PostDetails {
	if updatePost.Photo != nil {
		postDetails.Post.Photo = *updatePost.Photo
	}

	if updatePost.AnimalType != nil {
		postDetails.Animal.AnimalType = *updatePost.AnimalType
	}

	if updatePost.Color != nil {
		postDetails.Animal.Color = *updatePost.Color
	}

	if updatePost.Gender != nil {
		postDetails.Animal.Gender = *updatePost.Gender
	}

	if updatePost.Description != nil {
		postDetails.Post.Description = *updatePost.Description
	}

	if updatePost.Status != nil {
		postDetails.Animal.Status = *updatePost.Status
	}

	return postDetails
}
