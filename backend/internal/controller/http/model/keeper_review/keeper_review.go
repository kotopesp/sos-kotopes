package keeperreview

import (
	"time"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/user"
)

type CreateKeeperReview struct {
	Content *string `form:"content" validate:"omitempty,max=2000"`
	Grade   int     `form:"grade" validate:"required,numeric,min=1,max=5"`
}

type UpdateKeeperReview struct {
	Content *string `form:"content" validate:"omitempty,notblank,max=2000"`
	Grade   *int    `form:"grade" validate:"omitempty,numeric,min=1,max=5"`
}

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

type GetAllKeeperReviewsParams struct {
	Limit  *int `query:"limit" validate:"omitempty,gt=0"`
	Offset *int `query:"offset" validate:"omitempty,gte=0"`
}
