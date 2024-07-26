package core

import (
	"context"
	"time"
)

type KeeperReviews struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	AuthorID  int       `json:"author_id"`
	Content   string    `json:"content"`
	Grade     int       `json:"grade"`
	KeeperID  int       `json:"keeper_id"`
	IsDeleted bool      `json:"is_deleted"`
	DeletedAt time.Time `json:"deleted_at"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP()" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
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
