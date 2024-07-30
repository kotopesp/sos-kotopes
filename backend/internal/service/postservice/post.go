package postservice

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/internal/controller/http"
)

type (
	postService struct {
		PostStore core.PostStore
		PostFavouriteStore core.PostFavouriteStore
		AnimalStore core.AnimalStore
	}
)

func NewPostService(postStore core.PostStore, postFavouriteStore core.PostFavouriteStore, animalStore core.AnimalStore) core.PostService {
	return &postService{
		PostStore: 			postStore,
		PostFavouriteStore: postFavouriteStore,
		AnimalStore: 		animalStore,
	}
}

func (s *postService) GetAllPosts(ctx context.Context, limit, offset int) ([]core.Post, int, error) {

	posts, total, err := s.PostStore.GetAllPosts(ctx, limit, offset)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, 0, core.ErrPostNotFound
	}

	return posts, total, nil
}

func (s *postService) GetPostByID(ctx context.Context, id int) (core.Post, core.Animal, error) {
	post, err := s.PostStore.GetPostByID(ctx, id)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.Post{}, core.Animal{}, core.ErrPostNotFound
	}

	animal, err := s.AnimalStore.GetAnimalByID(ctx, post.AnimalID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.Post{}, core.Animal{}, core.ErrAnimalNotFound
	}

	return post, animal, nil
}

func (s *postService) CreatePost(ctx context.Context, post core.Post, fileHeader *multipart.FileHeader, animal core.Animal) (error) {
	photoBytes, err := http.GetPhotoBytes(fileHeader)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}
	post.Photo = *photoBytes
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	animal.CreatedAt = time.Now()
	animal.UpdatedAt = time.Now()

	animal, err = s.AnimalStore.CreateAnimal(ctx, animal)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}

	post.AnimalID = animal.ID
	err = s.PostStore.CreatePost(ctx, post)

	return err
}

func (s *postService) UpdatePost(ctx context.Context, post core.Post, animal core.Animal) error {
	err := s.PostStore.UpdatePost(ctx, post)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}
	err = s.AnimalStore.UpdateAnimal(ctx, animal)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}
	return nil
}

func (s *postService) DeletePost(ctx context.Context, id int) error {
	err := s.PostStore.DeletePost(ctx, id)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}
	return nil
}

func (s *postService) GetAuthorUsernameByID(ctx context.Context, authorID int) (string, error) {
	authorUsernames, err := s.PostStore.GetAuthorUsernameByID(ctx, authorID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return "", err
	}
	return authorUsernames, nil
}
