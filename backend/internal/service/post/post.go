package post

import (
	"context"
	"fmt"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

type service struct {
	postStore          core.PostStore
	postFavouriteStore core.PostFavouriteStore
	animalStore        core.AnimalStore
	userStore          core.UserStore
}

// New initializes a new instance of service
func New(postStore core.PostStore, postFavouriteStore core.PostFavouriteStore, animalStore core.AnimalStore, userStore core.UserStore) core.PostService {
	return &service{
		postStore:          postStore,
		postFavouriteStore: postFavouriteStore,
		animalStore:        animalStore,
		userStore:          userStore,
	}
}

// GetAllPosts retrieves all posts with the given parameters
func (s *service) GetAllPosts(ctx context.Context, userID int, params core.GetAllPostsParams) ([]core.PostDetails, int, error) {

	posts, total, err := s.postStore.GetAllPosts(ctx, userID, params)
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

// GetUserPosts retrieves all posts with the given user ID
func (s *service) GetUserPosts(ctx context.Context, id int) (postsDetails []core.PostDetails, count int, err error) {
	posts, total, err := s.postStore.GetUserPosts(ctx, id)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, 0, err
	}
	postsDetails, err = s.BuildPostDetailsList(ctx, posts, total)
	return postsDetails, total, err
}

// GetPostByID retrieves a post by its ID
func (s *service) GetPostByID(ctx context.Context, postID, userID int) (core.PostDetails, error) {
	post, err := s.postStore.GetPostByID(ctx, postID, userID)
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
func (s *service) CreatePost(ctx context.Context, postDetails core.PostDetails) (core.PostDetails, error) {
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
func (s *service) UpdatePost(ctx context.Context, postUpdateRequest core.UpdateRequestBodyPost) (core.PostDetails, error) {
	logger.Log().Debug(ctx, fmt.Sprintf("%v", *postUpdateRequest.ID))
	dbPost, err := s.GetPostByID(ctx, *postUpdateRequest.ID, *postUpdateRequest.AuthorID)
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
func (s *service) DeletePost(ctx context.Context, post core.Post) error {
	dbPost, err := s.postStore.GetPostByID(ctx, post.ID, post.AuthorID)
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
