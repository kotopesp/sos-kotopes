package post

import (
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/pagination"
	"time"
)

type (
	// CreateRequestBodyPost is the structure used for creating a new post
	CreateRequestBodyPost struct {
		Title       string `form:"title" json:"title" validate:"required,max=200"`
		Content     string `form:"content" json:"content" validate:"required,max=2000"`
		Photo       []byte `form:"photo" json:"photo"`
		AnimalType  string `form:"animal_type" json:"animal_type" validate:"required,oneof=dog cat"`
		Age         int    `form:"age" json:"age" validate:"gte=0"`
		Color       string `form:"color" json:"color" validate:"required"`
		Gender      string `form:"gender" json:"gender" validate:"required,oneof=male female"`
		Description string `form:"description" json:"description"`
		Status      string `form:"status" json:"status" validate:"required,oneof=lost found need_home"`
	}

	// PostResponse represents the structure of a post response with additional details
	PostResponse struct {
		Title          string    `form:"title" json:"title"`
		Content        string    `form:"content" json:"content"`
		AuthorUsername string    `form:"author_username" json:"author_username"`
		CreatedAt      time.Time `form:"created_at " json:"created_at"`
		Photo          []byte    `form:"photo" json:"photo"`
		AnimalType     string    `form:"animal_type" json:"animal_type"`
		Age            int       `form:"age" json:"age"`
		Color          string    `form:"color" json:"color"`
		Gender         string    `form:"gender" json:"gender"`
		Description    string    `form:"description" json:"description"`
		Status         string    `form:"status" json:"status"`
		IsFavourite    bool      `form:"is_favourite" json:"is_favourite"`
		Comments       int       `form:"comments" json:"comments"`
	}

	// UpdatePost, UpdateRequestBodyPost used for updating an existing post
	UpdatePost struct {
		Title   *string `form:"title" json:"title" validate:"max=200"`
		Content *string `form:"content" json:"content" validate:"max=2000"`
		Photo   *[]byte `form:"photo" json:"photo"`
	}

	UpdateRequestBodyPost struct {
		Title       *string `form:"title" json:"title" validate:"omitempty,max=200"`
		Content     *string `form:"content" json:"content" validate:"omitempty,max=2000"`
		Photo       *[]byte `form:"photo" json:"photo"`
		AnimalType  *string `form:"animal_type" json:"animal_type" validate:"omitempty,oneof=dog cat"`
		Age         *int    `form:"age" json:"age" validate:"omitempty,gte=0"`
		Color       *string `form:"color" json:"color"`
		Gender      *string `form:"gender" json:"gender" validate:"omitempty,oneof=male female"`
		Description *string `form:"description" json:"description"`
		Status      *string `form:"status" json:"status" validate:"omitempty,oneof=lost found need_home"`
	}

	// Meta represents metadata about the pagination of posts
	Meta struct {
		Total       int `json:"total"`
		TotalPages  int `json:"total_pages"`
		CurrentPage int `json:"current_page"`
		PerPage     int `json:"per_page"`
	}

	// Response represents the response structure for a list of posts with pagination
	Response struct {
		Meta  pagination.Pagination `json:"meta"`
		Posts []PostResponse        `json:"posts"`
	}

	// GetAllPostsParams represents the parameters for fetching a list of posts with filters
	GetAllPostsParams struct {
		Limit      int     `query:"limit" validate:"gt=0"`                                  // Limit on the number of posts to retrieve
		Offset     int     `query:"offset" validate:"gte=0"`                                // Offset for pagination
		Status     *string `query:"status" validate:"omitempty,oneof=lost found need_home"` // Filter by status of the associated animal
		AnimalType *string `query:"animal_type" validate:"omitempty,oneof=dog cat"`         // Filter by type of the associated animal
		Gender     *string `query:"gender" validate:"omitempty,oneof=male female"`          // Filter by gender of the associated animal
		Color      *string `query:"color" validate:"omitempty"`								              // Filter by color of the associated animal
		Location   *string `query:"location" validate:"omitempty"` 							            // Filter by location of the associated animal
		SearchWord *string `query:"search_word" validate:"omitempty"`                       // Filter by location of the associated animal
	}

	PathParams struct {
		PostID int `params:"id" validate:"omitempty,gt=0"`
	}
)
