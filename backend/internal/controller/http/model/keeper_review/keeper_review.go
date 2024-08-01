package keeperreview

import "time"

// KeeperReviews represents the review entity for keepers.
type KeeperReviews struct {
	ID        int       `json:"id"`
	AuthorID  int       `json:"author_id" validate:"required,min=0"`
	Content   string    `json:"content" validate:"required,notblank,max=2000"`
	Grade     int       `json:"grade" validate:"required,numeric,min=1,max=5"`
	KeeperID  int       `json:"keeper_id" validate:"required,min=0"`
	IsDeleted bool      `json:"is_deleted"`
	DeletedAt time.Time `json:"deleted_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// KeeperReviewsCreate represents the data required to create a new keeper review.
type KeeperReviewsCreate struct {
	AuthorID int    `json:"author_id" validate:"required,min=0"`
	Content  string `json:"content" validate:"required,notblank,max=2000"`
	Grade    int    `json:"grade" validate:"required,numeric,min=1,max=5"`
	KeeperID int    `json:"keeper_id" validate:"required,min=0"`
}

// KeeperReviewsUpdate represents the data to update an existing keeper review.
type KeeperReviewsUpdate struct {
	Content string `json:"content" validate:"notblank,max=2000"`
	Grade   int    `json:"grade" validate:"numeric,min=1,max=5"`
}

// GetAllKeeperReviewsParams represents the query parameters for fetching multiple keeper reviews.
type GetAllKeeperReviewsParams struct {
	Limit  int `query:"limit" validate:"omitempty,gt=0"`
	Offset int `query:"offset" validate:"omitempty,gte=0"`
}
