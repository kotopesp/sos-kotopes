package post

import (
	"context"
	"mime/multipart"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/internal/controller/http"
)

type (
	postService struct {
		postStore core.PostStore
		postFavouriteStore core.PostFavouriteStore
		animalStore core.AnimalStore
		userStore core.UserStore
	}
)

func New(postStore core.PostStore, postFavouriteStore core.PostFavouriteStore, animalStore core.AnimalStore, userStore core.UserStore) core.PostService {
	return &postService{
		postStore: 			postStore,
		postFavouriteStore: postFavouriteStore,
		animalStore: 		animalStore,
		userStore:          userStore,
	}
}

func (s *postService) GetAllPosts(ctx context.Context, params core.GetAllPostsParams) ([]core.PostDetails, int, error) {

	posts, total, err := s.postStore.GetAllPosts(ctx, params)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, 0, err
	}

	postDetails := make([]core.PostDetails, total)
	for i, post := range posts {
		animal, err := s.animalStore.GetAnimalByID(ctx, post.AnimalID)
		if err != nil {
			logger.Log().Error(ctx, err.Error())
			return nil, 0, err
		}
		user, err := s.userStore.GetUserByID(ctx, post.AuthorID)
		if err != nil {
			logger.Log().Error(ctx, err.Error())
			return nil, 0, err
		}

		postDetails[i] = core.PostDetails{
			Post: post,
			Animal: animal,
			Username: user.Username,
		}
	}

	return postDetails, total, nil
}

func (s *postService) GetPostByID(ctx context.Context, id int) (core.PostDetails, error) {
	post, err := s.postStore.GetPostByID(ctx, id)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.PostDetails{}, err
	}

	animal, err := s.animalStore.GetAnimalByID(ctx, post.AnimalID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.PostDetails{}, err
	}

	user, err := s.userStore.GetUserByID(ctx, post.AuthorID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.PostDetails{}, err
	}

	postDetails := core.PostDetails{
		Post: post,
		Animal: animal,
		Username: user.Username,
	}

	return postDetails, nil
}

func (s *postService) CreatePost(ctx context.Context, postDetails core.PostDetails, fileHeader *multipart.FileHeader) (core.PostDetails, error) {
	photoBytes, err := http.GetPhotoBytes(fileHeader)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.PostDetails{}, err
	}

	postDetails.Post.Photo = *photoBytes

	animal, err := s.animalStore.CreateAnimal(ctx, postDetails.Animal)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.PostDetails{}, err
	}

	postDetails.Post.AnimalID = animal.ID
	post, err := s.postStore.CreatePost(ctx, postDetails.Post)

	user, err := s.userStore.GetUserByID(ctx, post.AuthorID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.PostDetails{}, err
	}

	var createPostDetails core.PostDetails

	createPostDetails.Post = post
	createPostDetails.Animal = animal
	createPostDetails.Username = user.Username

	return createPostDetails, err
}

func (s *postService) UpdatePost(ctx context.Context, postDetails core.PostDetails) (core.PostDetails, error) {
	post, err := s.postStore.UpdatePost(ctx, postDetails.Post)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.PostDetails{}, err
	}
	
	animal, err := s.animalStore.UpdateAnimal(ctx, postDetails.Animal)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.PostDetails{}, err
	}

	var updatePostDetails core.PostDetails

	updatePostDetails.Post = post
	updatePostDetails.Animal = animal

	return updatePostDetails, nil
}

func (s *postService) DeletePost(ctx context.Context, id int) error {
	err := s.postStore.DeletePost(ctx, id)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}
	return nil
}
