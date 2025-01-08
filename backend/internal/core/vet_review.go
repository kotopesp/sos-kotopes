package core

import (
	"context"
	"time"
)

type VetReviews struct {
	ID        int       `gorm:"primaryKey;autoIncrement;column:id"`
	AuthorID  int       `gorm:"column:author_id"`
	Content   string    `gorm:"column:content"`
	Grade     int       `gorm:"column:grade"`
	VetID     int       `gorm:"column:vet_id"`
	IsDeleted bool      `gorm:"column:is_deleted"`
	DeletedAt time.Time `gorm:"column:deleted_at"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP();column:created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

type UpdateVetReviews struct {
	ID        int       `gorm:"primaryKey;autoIncrement;column:id"`
	AuthorID  int       `gorm:"column:author_id"`
	Content   string    `gorm:"column:content"`
	Grade     int       `gorm:"column:grade"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

type VetReviewsDetails struct {
	Review VetReviews
	User   User
}

type GetAllVetReviewsParams struct {
	Limit  *int
	Offset *int
}

type VetReviewsStore interface {
	GetAllReviews(ctx context.Context, params GetAllVetReviewsParams) ([]VetReviews, error)
	GetByIDReview(ctx context.Context, id int) (VetReviews, error)
	CreateReview(ctx context.Context, vetReview VetReviews) error
	SoftDeleteReviewByID(ctx context.Context, id int) error
	UpdateReviewByID(ctx context.Context, vetReview UpdateVetReviews) (VetReviews, error)
}

type VetReviewsService interface {
	GetAllReviews(ctx context.Context, params GetAllVetReviewsParams) ([]VetReviewsDetails, error)
	CreateReview(ctx context.Context, vetReview VetReviews) error
	SoftDeleteReviewByID(ctx context.Context, id int, userID int) error
	UpdateReviewByID(ctx context.Context, vetReview UpdateVetReviews) (VetReviewsDetails, error)
}

// TableName table name in db for gorm
func (VetReviews) TableName() string {
	return "vet_reviews"
}
