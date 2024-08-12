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

type UpdateKeeperReviews struct {
	ID        int       `gorm:"primaryKey;autoIncrement;column:id"`
	Content   string    `gorm:"column:content"`
	Grade     int       `gorm:"column:grade"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

type KeeperReviewsDetails struct {
	Review KeeperReviews
	User   User
}

type GetAllKeeperReviewsParams struct {
	Limit  *int
	Offset *int
}

type KeeperReviewsStore interface {
	GetAllReviews(ctx context.Context, params GetAllKeeperReviewsParams) ([]KeeperReviews, error)
	CreateReview(ctx context.Context, keeperReview KeeperReviews) error
	DeleteReviewByID(ctx context.Context, id int) error
	SoftDeleteReviewByID(ctx context.Context, id int) error
	UpdateReviewByID(ctx context.Context, keeperReview UpdateKeeperReviews) (KeeperReviews, error)
}

type KeeperReviewsService interface {
	GetAllReviews(ctx context.Context, params GetAllKeeperReviewsParams) ([]KeeperReviewsDetails, error)
	CreateReview(ctx context.Context, keeperReview KeeperReviews) error
	DeleteReviewByID(ctx context.Context, id int) error
	SoftDeleteReviewByID(ctx context.Context, id int) error
	UpdateReviewByID(ctx context.Context, keeperReview UpdateKeeperReviews) (KeeperReviewsDetails, error)
}

// TableName table name in db for gorm
func (KeeperReviews) TableName() string {
	return "keeper_reviews"
}
