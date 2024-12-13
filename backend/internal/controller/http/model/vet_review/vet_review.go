package vet_review

import (
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/user"
	"time"
)

type VetReviewsCreate struct {
	AuthorID int    `form:"author_id" validate:"required,min=0"`
	Content  string `form:"content" validate:"required,notblank,max=2000"`
	Grade    int    `form:"grade" validate:"required,numeric,min=1,max=5"`
	VetID    int    `form:"vet_id" validate:"required,min=0"`
}

type VetReviewsUpdate struct {
	ID       int    `form:"id"`
	AuthorID int    `form:"author_id" validate:"required,min=0"`
	Content  string `form:"content" validate:"notblank,max=2000"`
	Grade    int    `form:"grade" validate:"numeric,min=1,max=5"`
}

type VetReviewsResponse struct {
	ID        int       `json:"id"`
	AuthorID  int       `json:"author_id" validate:"required,min=0"`
	Content   string    `json:"content" validate:"required,notblank,max=2000"`
	Grade     int       `json:"grade" validate:"required,numeric,min=1,max=5"`
	VetID     int       `json:"vet_id" validate:"required,min=0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VetReviewsResponseWithUser struct {
	Review VetReviewsResponse
	User   user.ResponseUser
}

type GetAllVetReviewsParams struct {
	Limit  int `query:"limit" validate:"omitempty,gt=0"`
	Offset int `query:"offset" validate:"omitempty,gte=0"`
}
