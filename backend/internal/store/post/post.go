package poststore

import (
	"context"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"time"
	"errors"
)

type (
	store struct {
		*postgres.Postgres
	}
)

func New(pg *postgres.Postgres) core.PostStore {
	return &store{pg}
}


func (s *store) GetAllPosts(ctx context.Context, params core.GetAllPostsParams) ([]core.Post, int, error) {
    var posts []core.Post

    query := s.DB.WithContext(ctx).Model(&core.Post{}).
			Joins("JOIN animals ON posts.animal_id = animals.id").
			Where("posts.is_deleted = ?", false)

	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}

	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	if params.Status != nil {
		query = query.Where("animals.status = ?", *params.Status)
	}

	if params.AnimalType != nil {
		query = query.Where("animals.animal_type = ?", *params.AnimalType)
	}

	if params.Gender != nil {
		query = query.Where("animals.gender = ?", *params.Gender)
	}

	if params.Color != nil {
		query = query.Where("animals.color = ?", *params.Color)
	}

    var total int64
    if err := query.Count(&total).Error; err != nil {
		logger.Log().Error(ctx, err.Error())
        return nil, 0, err
    }

    if err := query.Select("posts.*").Find(&posts).Error; err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return nil, 0, core.ErrPostNotFound
		}
		logger.Log().Error(ctx, err.Error())
		return nil, 0, err
    }

    return posts, int(total), nil
}

func (s *store) GetPostByID(ctx context.Context, id int) (core.Post, error) {
	var post core.Post

	if err := s.DB.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&post).Error; err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.Post{}, core.ErrPostNotFound
		}

		logger.Log().Error(ctx, err.Error())
		return core.Post{}, err
	}

	return post, nil
}

func (s *store) CreatePost(ctx context.Context, post core.Post) (core.Post, error) {
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	var createdPost core.Post

	if err := s.DB.WithContext(ctx).Create(&post).First(&createdPost, post.ID).Error; err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.Post{}, err
	}
	
	return createdPost, nil
}

func (s *store) UpdatePost(ctx context.Context, post core.Post) (core.Post, error) {
	post.UpdatedAt = time.Now()

	var updatedPost core.Post

	if err := s.DB.WithContext(ctx).Save(&post).First(&updatedPost, post.ID).Error; err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.Post{}, core.ErrPostNotFound
		}

		logger.Log().Error(ctx, err.Error())
		return core.Post{}, err
	}
	
	return updatedPost, nil
}

func (s *store) DeletePost(ctx context.Context, id int) error{
	updates := map[string]interface{}{
        "is_deleted": true,
        "deleted_at": time.Now(), 
    }

	result := s.DB.WithContext(ctx).Model(&core.Post{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		if errors.Is(result.Error, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.ErrPostNotFound
		}
		logger.Log().Error(ctx, result.Error.Error())
        return result.Error
    }

    if result.RowsAffected == 0 {
		logger.Log().Error(ctx, core.ErrPostNotFound.Error())
        return core.ErrPostNotFound
    }

    return nil
}
