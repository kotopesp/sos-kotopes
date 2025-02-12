package poststore

import (
	"context"
	"errors"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
	"gorm.io/gorm"
	"time"
)

type store struct {
	posts         *postgres.Postgres
	reportedPosts *postgres.Postgres
}

func New(pg *postgres.Postgres) core.PostStore {
	return &store{pg, pg}
}

// GetAllPosts retrieves all posts from the database based on the GetAllPostsParams
func (s *store) GetAllPosts(ctx context.Context, params core.GetAllPostsParams) ([]core.Post, int, error) {
	var posts []core.Post

	query := s.posts.DB.WithContext(ctx).Model(&core.Post{}).
		Joins("JOIN animals ON posts.animal_id = animals.id").
		Where("posts.is_deleted = ?", false)

	// Apply filtering based on the GetAllPostsParams
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
		logger.Log().Error(ctx, err.Error())
		return nil, 0, err
	}

	return posts, int(total), nil
}

// GetUserPosts retrieves all posts from the database based on the given user ID
func (s *store) GetUserPosts(ctx context.Context, id int) (posts []core.Post, count int, err error) {
	err = s.posts.DB.WithContext(ctx).
		Where("author_id = ?", id).
		Order("created_at DESC").
		Find(&posts).Error
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, core.ErrNoSuchUser
		}
		return nil, 0, err
	}

	count = len(posts)
	return posts, count, nil
}

// GetPostByID retrieves a post from the database by its ID
func (s *store) GetPostByID(ctx context.Context, id int) (core.Post, error) {
	var post core.Post

	if err := s.posts.DB.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&post).Error; err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.Post{}, core.ErrPostNotFound
		}

		logger.Log().Error(ctx, err.Error())
		return core.Post{}, err
	}

	return post, nil
}

// CreatePost inserts a new post record into the database
func (s *store) CreatePost(ctx context.Context, post core.Post) (core.Post, error) {
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	var createdPost core.Post

	if err := s.posts.DB.WithContext(ctx).Create(&post).First(&createdPost, post.ID).Error; err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.Post{}, err
	}

	return createdPost, nil
}

// UpdatePost updates an existing post record in the database
func (s *store) UpdatePost(ctx context.Context, post core.Post) (core.Post, error) {
	post.UpdatedAt = time.Now()

	var updatedPost core.Post

	if err := s.posts.DB.WithContext(ctx).Save(&post).First(&updatedPost, post.ID).Error; err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.Post{}, core.ErrPostNotFound
		}

		logger.Log().Error(ctx, err.Error())
		return core.Post{}, err
	}

	return updatedPost, nil
}

// DeletePost marks a post as deleted in the database by updating the is_deleted flag and setting the deleted_at timestamp
func (s *store) DeletePost(ctx context.Context, id int) error {
	updates := map[string]interface{}{
		"is_deleted": true,
		"deleted_at": time.Now(),
	}

	result := s.posts.DB.WithContext(ctx).Model(&core.Post{}).Where("id = ?", id).Updates(updates)
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

// SendToModeration - метод, для сохранения постов с определенным числом репортов, в отдельную таблицу.
func (s *store) SendToModeration(ctx context.Context, post core.Post) (err error) {
	reportedPost := core.ReportedPost{
		Post:             post,
		SentToModeration: time.Now(),
	}

	if err = s.reportedPosts.DB.WithContext(ctx).Create(reportedPost).First(post.ID).Error; err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}

	return nil
}

func (s *store) ReportPost(ctx context.Context, post core.Post, reason string) (reportedPost core.Post, err error) {
	err = s.posts.DB.WithContext(ctx).Model(&reportedPost).Where("id = ?", post.ID).Error
	if err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.Post{}, core.ErrPostNotFound
		}
		logger.Log().Error(ctx, err.Error())
		return core.Post{}, err
	}

	reportedPost.Reports.Number++
	reportedPost.Reports.Reasons = append(reportedPost.Reports.Reasons, reason)
	reportedPost.LastReportedAt = time.Now()
	// some threshold meaning we use, to decide about moderation
	if reportedPost.Reports.Number > 15 {
		err = s.SendToModeration(ctx, reportedPost)
		if err != nil {
			return core.Post{}, err
		}
	}

	if err = s.posts.DB.WithContext(ctx).Save(&reportedPost).Error; err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.Post{}, err
	}

	return reportedPost, err
}
