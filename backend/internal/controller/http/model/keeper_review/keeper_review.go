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
	Content  string `json:"content"`
	Grade    int    `json:"grade"`
	KeeperID int    `json:"keeper_id"`
}

type KeeperReviewsUpdate struct {
	Content string `json:"content"`
	Grade   int    `json:"grade"`
}

type GetAllKeeperReviewsParams struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}
