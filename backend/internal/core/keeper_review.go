package core

import (
	"context"
	"time"
)

type KeeperReviews struct {
	ID        int       `gorm:"primaryKey;autoIncrement;column:id"`
	AuthorID  int       `gorm:"column:author_id"`
	Content   string    `gorm:"column:content"`
	Grade     int       `gorm:"column:grade"`
	KeeperID  int       `gorm:"column:keeper_id"`
	IsDeleted bool      `gorm:"column:is_deleted"`
	DeletedAt time.Time `gorm:"column:deleted_at"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP();column:created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

type GetAllKeeperReviewsParams struct {
	Limit  *int
	Offset *int
}

type KeeperReviewsStore interface {
	GetAll(ctx *context.Context, params GetAllKeeperReviewsParams) ([]KeeperReviews, error)
	Create(ctx *context.Context, keeperReview KeeperReviews) error
	DeleteByID(ctx *context.Context, id int) error
	SoftDeleteByID(ctx *context.Context, id int) error
	UpdateByID(ctx *context.Context, keeperReview KeeperReviews) error
}

type KeeperReviewsService interface {
	GetAll(ctx *context.Context, params GetAllKeeperReviewsParams) ([]KeeperReviews, error)
	Create(ctx *context.Context, keeperReview KeeperReviews) error
	DeleteByID(ctx *context.Context, id int) error
	SoftDeleteByID(ctx *context.Context, id int) error
	UpdateByID(ctx *context.Context, keeperReview KeeperReviews) error
}
