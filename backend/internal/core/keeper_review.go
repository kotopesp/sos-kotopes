package core

import (
	"context"
	"time"
)

type KeeperReview struct {
	ID        int        `gorm:"primaryKey;autoIncrement;column:id"`
	AuthorID  int        `gorm:"column:author_id"`
	Author    User       `gorm:"foreignKey:AuthorID;references:ID"`
	Content   *string    `gorm:"column:content"`
	Grade     int        `gorm:"column:grade"`
	KeeperID  int        `gorm:"column:keeper_id"`
	Keeper    Keeper     `gorm:"foreignKey:KeeperID;references:ID"`
	IsDeleted bool       `gorm:"column:is_deleted"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP();column:created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime;column:updated_at"`
}

// type UpdateKeeperReview struct {
// 	Content *string `gorm:"column:content"`
// 	Grade   *int    `gorm:"column:grade"`
// }

type GetAllKeeperReviewsParams struct {
	Limit  *int
	Offset *int
}

type KeeperReviewStore interface {
	GetAllReviews(ctx context.Context, keeperID int, params GetAllKeeperReviewsParams) (data []KeeperReview, err error)
	GetReviewByID(ctx context.Context, id int) (KeeperReview, error)
	CreateReview(ctx context.Context, review KeeperReview) error
	DeleteReview(ctx context.Context, id int) error
	UpdateReview(ctx context.Context, id int, review KeeperReview) (data KeeperReview, err error)
}

type KeeperReviewService interface {
	GetAllReviews(ctx context.Context, keeperID int, params GetAllKeeperReviewsParams) (data []KeeperReview, err error)
	CreateReview(ctx context.Context, review KeeperReview) (data KeeperReview, err error)
	DeleteReview(ctx context.Context, id int, userID int) error
	UpdateReview(ctx context.Context, id int, userID int, review KeeperReview) (data KeeperReview, err error)
}

// TableName table name in db for gorm
func (KeeperReview) TableName() string {
	return "keeper_reviews"
}
