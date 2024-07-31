package post

import (
	"time"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/pagination"
)

type (
	Post struct {
		Title     string `form:"title" json:"title" validate:"required,max=200"`
		Content   string `form:"content" json:"content" validate:"required,max=2000"`
		Photo     []byte `form:"photo" json:"photo"`
	}

	CreateRequestBodyPost struct {
		Title     string   `form:"title" json:"title" validate:"required,max=200"`
		Content   string   `form:"content" json:"content" validate:"required,max=2000"`
		Photo     []byte   `form:"photo" json:"photo"`
		AnimalType  string `form:"animal_type" json:"animal_type" validate:"required"` 
		Age         int    `form:"age" json:"age" validate:"gte=0"`
		Color       string `form:"color" json:"color" validate:"required"`
		Gender      string `form:"gender" json:"gender" validate:"required"`
		Description string `form:"description" json:"description"`
		Status      string `form:"status" json:"status" validate:"required"`
	}

	PostResponse struct {
		Title          string   `form:"title" json:"title"`
		Content        string    `form:"content" json:"content"`
		AuthorUsername string 	 `form:"author_username" json:"author_username"`
		CreatedAt      time.Time `form:"created_at " json:"created_at"`
		Photo          []byte    `form:"photo" json:"photo"`
		AnimalType     string 	 `form:"animal_type" json:"animal_type"`
		Age            int    	 `form:"age" json:"age"`
		Color          string 	 `form:"color" json:"color"`
		Gender         string 	 `form:"gender" json:"gender"`
		Description    string 	 `form:"description" json:"description"`
		Status         string 	 `form:"status" json:"status"`
		Comments       int 		 `form:"comments" json:"comments"`
	}

	UpdatePost struct {
		Title   *string `form:"title" json:"title" validate:"max=200"`
		Content *string `form:"content" json:"content" validate:"max=2000"`
		Photo   *[]byte `form:"photo" json:"photo"`
	}

	UpdateRequestBodyPost struct {
		Title   	*string `form:"title" json:"title" validate:"omitempty,max=200"`
		Content 	*string `form:"content" json:"content" validate:"omitempty,max=2000"`
		Photo   	*[]byte `form:"photo" json:"photo"`
		AnimalType  *string `form:"animal_type" json:"animal_type"`
		Age         *int    `form:"age" json:"age"`
		Color       *string `form:"color" json:"color"`
		Gender      *string `form:"gender" json:"gender"`
		Description *string `form:"description" json:"description"`
		Status      *string `form:"status" json:"status"`
	}

	Meta struct {
		Total       int `json:"total"`
		TotalPages  int `json:"total_pages"`
		CurrentPage int `json:"current_page"`
		PerPage     int `json:"per_page"`
	}
	
	Response struct {
		Meta  pagination.Pagination `json:"meta"`
		Posts []PostResponse		`json:"posts"`
	}

	GetAllPostsParams struct {
		Limit      int     `query:"limit" validate:"gt=0"`
		Offset     int     `query:"offset" validate:"gte=0"`
		Status     *string `query:"status" validate:"omitempty,oneof=lost found need_home"`
		AnimalType *string `query:"animal_type" validate:"omitempty,oneof=dog cat"`
		Gender     *string `query:"gender" validate:"omitempty,oneof=male female"`
		Color      *string `query:"color" validate:"omitempty"`
		Location   *string `query:"location" validate:"omitempty"`
	}
)
