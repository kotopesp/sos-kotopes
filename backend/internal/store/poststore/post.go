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

func NewPostStore(pg *postgres.Postgres) core.PostStore {
	return &store{pg}
}

func (s *store) GetAuthorUsernameByID(ctx context.Context, authorID int) (string, error) {
    var username string

    err := s.DB.WithContext(ctx).Model(&core.User{}).Select("username").Where("id = ?", authorID).Scan(&username).Error
    if err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return "", core.ErrUsernameNotFound
		}
		logger.Log().Error(ctx, err.Error())
		return "", err
	}

	return username, nil
}

func (s *store) GetAllPosts(ctx context.Context, limit, offset int) ([]core.Post, int, error) {
    var posts []core.Post

    query := s.DB.WithContext(ctx).Model(&core.Post{}).Where("is_deleted = ?", false)

    var total int64
    if err := query.Count(&total).Error; err != nil {
		logger.Log().Error(ctx, err.Error())
        return nil, 0, err
    }

    if err := query.Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
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

func (s *store) CreatePost(ctx context.Context, post core.Post) error {
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	if err := s.DB.WithContext(ctx).Create(&post).Error; err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}
	
	return nil
}

func (s *store) UpdatePost(ctx context.Context, post core.Post) error {
	post.UpdatedAt = time.Now()

	if err := s.DB.WithContext(ctx).Save(&post).Error; err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.ErrPostNotFound
		}

		logger.Log().Error(ctx, err.Error())
		return err
	}
	
	return nil
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
