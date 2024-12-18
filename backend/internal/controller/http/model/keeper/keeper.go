package keeper

import (
	"time"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/pagination"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/user"
)

// KeepersCreate represents the data required to create a new keeper.
type CreateKeeper struct {
	UserID               int      `form:"user_id" validate:"required,min=0"`
	Description          *string  `form:"description" validate:"notblank,max=600"`
	Price                *float64 `form:"price" validate:"min=0"`
	LocationID           *int     `form:"location_id"`
	HasCage              bool     `form:"has_cage" validate:"required,boolean"`
	BoardingDuration     string   `form:"boarding_duration" validate:"required,oneof=hours days weeks months depends"`
	BoardingCompensation string   `form:"boarding_compensation" validate:"required,oneof=paid free depends"`
	AnimalAcceptance     string   `form:"animal_acceptance" validate:"required,oneof=homeless home depends"`
	AnimalCategory       string   `form:"animal_category" validate:"required,oneof=dog cat both"`
}

// KeepersUpdate represents the data to update an existing keeper.
type UpdateKeeper struct {
	Description          *string  `form:"description" validate:"notblank,max=600"`
	Price                *float64 `form:"price" validate:"min=0"`
	LocationID           *int     `form:"location_id"`
	HasCage              *bool    `form:"has_cage" validate:"boolean"`
	BoardingDuration     *string  `form:"boarding_duration" validate:"oneof=hours days weeks months depends"`
	BoardingCompensation *string  `form:"boarding_compensation" validate:"oneof=paid free depends"`
	AnimalAcceptance     *string  `form:"animal_acceptance" validate:"oneof=homeless home depends"`
	AnimalCategory       *string  `form:"animal_category" validate:"oneof=dog cat both"`
}

// KeepersResponse represents the response keeper entity.
type ResponseKeeper struct {
	ID                   int               `json:"id"`
	UserID               int               `json:"user_id"`
	User                 user.ResponseUser `json:"user"`
	Description          *string           `json:"description"`
	Price                *float64          `json:"price"`
	LocationID           *int              `json:"location_id"`
	HasCage              bool              `json:"has_cage"`
	BoardingDuration     string            `json:"boarding_duration"`
	BoardingCompensation string            `json:"boarding_compensation"`
	AnimalAcceptance     string            `json:"animal_acceptance"`
	AnimalCategory       string            `json:"animal_category"`
	CreatedAt            time.Time         `json:"created_at"`
	UpdatedAt            time.Time         `json:"updated_at"`
}

// KeepersResponseWithMeta represents the respose keeper entity with meta information.
type ResponseKeepers struct {
	Meta pagination.Pagination `json:"meta"`
	Data []ResponseKeeper      `json:"payload"`
}

// GetAllKeepersParams represents the query parameters for filtering and sorting keepers.
type GetAllKeepersParams struct {
	Sort                 *string  `query:"sort" validate:"omitempty,sort_keeper"`
	LocationID           *int     `query:"location_id" validate:"omitempty"`
	MinRating            *float64 `query:"min_rating" validate:"omitempty,gte=1,lte=5"`
	MaxRating            *float64 `query:"max_rating" validate:"omitempty,gte=1,lte=5"`
	MinPrice             *float64 `query:"min_price" validate:"omitempty,gte=0"`
	MaxPrice             *float64 `query:"max_price" validate:"omitempty,gte=0"`
	HasCage              *bool    `query:"has_cage" validate:"omitempty,boolean"`
	BoardingDuration     *string  `query:"boarding_duration" validate:"omitempty,oneof=hours days weeks months depends"`
	BoardingCompensation *string  `query:"boarding_compensation" validate:"omitempty,oneof=paid free depends"`
	AnimalAcceptance     *string  `query:"animal_acceptance" validate:"omitempty,oneof=homeless home depends"`
	AnimalCategory       *string  `query:"animal_category" validate:"omitempty,oneof=dog cat both"`
	Limit                *int     `query:"limit" validate:"omitempty,gt=0"`
	Offset               *int     `query:"offset" validate:"omitempty,gte=0"`
}
