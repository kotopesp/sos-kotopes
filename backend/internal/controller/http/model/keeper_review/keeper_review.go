package keeperreview

import "time"

type KeeperReviews struct {
	ID        int       `json:"id"`
	AuthorID  int       `json:"author_id"`
	Content   string    `json:"content"`
	Grade     int       `json:"grade"`
	KeeperID  int       `json:"keeper_id"`
	IsDeleted bool      `json:"is_deleted"`
	DeletedAt time.Time `json:"deleted_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type KeeperReviewsCreate struct {
	AuthorID int    `json:"author_id"`
	Content  string `json:"content" validate:"required,notblank,max=2000"`
	Grade    int    `json:"grade" validate:"required,numeric_natural,min=1,max=5"`
	KeeperID int    `json:"keeper_id"`
}

type KeeperReviewsUpdate struct {
	Content string `json:"content" validate:"notblank,max=2000"`
	Grade   int    `json:"grade" validate:"numeric_natural,min=1,max=5"`
}

type GetAllKeeperReviewsParams struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}
