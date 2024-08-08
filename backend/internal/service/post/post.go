package post

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/kotopesp/sos-kotopes/internal/controller/http"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

type (
	postService struct {
		postStore          core.PostStore
		postFavouriteStore core.PostFavouriteStore
		animalStore        core.AnimalStore
		userStore          core.UserStore
	}
)

// New initializes a new instance of postService
func New(postStore core.PostStore, postFavouriteStore core.PostFavouriteStore, animalStore core.AnimalStore, userStore core.UserStore) core.PostService {
	return &postService{
		postStore:          postStore,
		postFavouriteStore: postFavouriteStore,
		animalStore:        animalStore,
		userStore:          userStore,
	}
}

// GetAllPosts retrieves all posts with the given parameters
func (s *postService) GetAllPosts(ctx context.Context, params core.GetAllPostsParams) ([]core.PostDetails, int, error) {

	posts, total, err := s.postStore.GetAllPosts(ctx, params)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, 0, err
	}

	postDetails, err := s.BuildPostDetailsList(ctx, posts, total)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, 0, err
	}

	return postDetails, total, nil
}

// GetPostByID retrieves a post by its ID
func (s *postService) GetPostByID(ctx context.Context, id int) (core.PostDetails, error) {
	post, err := s.postStore.GetPostByID(ctx, id)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.PostDetails{}, err
	}

	postDetails, err := s.BuildPostDetails(ctx, post)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.PostDetails{}, err
	}

	return postDetails, nil
}

// CreatePost creates a new post with the provided details and photo
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
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.PostDetails{}, err
	}

	user, err := s.userStore.GetUserByID(ctx, post.AuthorID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.PostDetails{}, err
	}

	createPostDetails := ToCorePostDetails(post, animal, user.Username)

	return createPostDetails, err
}

// UpdatePost updates an existing post with the provided details
func (s *postService) UpdatePost(ctx context.Context, postUpdateRequest core.UpdateRequestBodyPost) (core.PostDetails, error) {
	logger.Log().Debug(ctx, fmt.Sprintf("%v", *postUpdateRequest.ID))
	dbPost, err := s.GetPostByID(ctx, *postUpdateRequest.ID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.PostDetails{}, err
	}

	if dbPost.Post.AuthorID != *postUpdateRequest.AuthorID {
		return core.PostDetails{}, core.ErrPostAuthorIDMismatch
	} else if dbPost.Post.IsDeleted {
		return core.PostDetails{}, core.ErrPostIsDeleted
	}

	dbPost = FuncUpdateRequestBodyPost(dbPost, postUpdateRequest)

	post, err := s.postStore.UpdatePost(ctx, dbPost.Post)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.PostDetails{}, err
	}

	animal, err := s.animalStore.UpdateAnimal(ctx, dbPost.Animal)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.PostDetails{}, err
	}

	var updatePostDetails core.PostDetails

	updatePostDetails.Post = post
	updatePostDetails.Animal = animal

	return updatePostDetails, nil
}

// DeletePost deletes a post by its ID
func (s *postService) DeletePost(ctx context.Context, post core.Post) error {
	dbPost, err := s.postStore.GetPostByID(ctx, post.ID)
	if err != nil {
		return err
	}

	if dbPost.AuthorID != post.AuthorID {
		return core.ErrPostAuthorIDMismatch
	} else if dbPost.IsDeleted {
		return core.ErrPostIsDeleted
	}

	err = s.postStore.DeletePost(ctx, post.ID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}

	return nil
}
