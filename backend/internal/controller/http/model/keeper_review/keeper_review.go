package keeperreview

import (
	"time"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/user"
)

// KeeperReviewsCreate represents the data required to create a new keeper review.
type CreateKeeperReview struct {
	AuthorID int     `form:"author_id" validate:"required,min=0"`
	Content  *string `form:"content" validate:"max=2000"`
	Grade    int     `form:"grade" validate:"required,numeric,min=1,max=5"`
	KeeperID int     `form:"keeper_id" validate:"required,min=0"`
}

// KeeperReviewsUpdate represents the data to update an existing keeper review.
type UpdateKeeperReview struct {
	Content *string `form:"content" validate:"notblank,max=2000"`
	Grade   *int    `form:"grade" validate:"numeric,min=1,max=5"`
}

// KeeperReviewsResponse represents the data to send keeper review back to client.
type ResponseKeeperReview struct {
	ID        int               `json:"id"`
	AuthorID  int               `json:"author_id"`
	Author    user.ResponseUser `json:"author"`
	Content   *string           `json:"content"`
	Grade     int               `json:"grade"`
	KeeperID  int               `json:"keeper_id"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// GetAllKeeperReviewsParams represents the query parameters for fetching multiple keeper reviews.
type GetAllKeeperReviewsParams struct {
	Limit  *int `query:"limit" validate:"omitempty,gt=0"`
	Offset *int `query:"offset" validate:"omitempty,gte=0"`
}
