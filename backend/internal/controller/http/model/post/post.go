package post

import (
	"time"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/pagination"
)

type (
	// CreateRequestBodyPost is the structure used for creating a new post
	CreateRequestBodyPost struct {
		Title       string   `form:"title" json:"title" validate:"required,max=200" example:"Title of the post"`
		Content     string   `form:"content" json:"content" validate:"required,max=2000" example:"Content of the post"`
		Photos      [][]byte `form:"photos" json:"photos"`
		AnimalType  string   `form:"animal_type" json:"animal_type" validate:"required,oneof=dog cat" example:"dog"`
		Age         int      `form:"age" json:"age" validate:"gte=0" example:"2"`
		Color       string   `form:"color" json:"color" validate:"required" example:"brown"`
		Gender      string   `form:"gender" json:"gender" validate:"required,oneof=male female" example:"male"`
		Description string   `form:"description" json:"description" example:"Description of the post"`
		Status      string   `form:"status" json:"status" validate:"required,oneof=lost found need_home" example:"lost"`
	}

	// PostResponse represents the structure of a post response with additional details
	PostResponse struct {
		Title          string    `form:"title" json:"title" example:"Title of the post"`
		Content        string    `form:"content" json:"content" example:"Content of the post"`
		AuthorUsername string    `form:"author_username" json:"author_username" example:"username"`
		CreatedAt      time.Time `form:"created_at " json:"created_at" example:"2006-01-02T15:04:05Z07:00"`
		URLsPhotos     []string  `form:"urls_photos" json:"urls_photos"`
		AnimalType     string    `form:"animal_type" json:"animal_type" example:"dog"`
		Age            int       `form:"age" json:"age" example:"2"`
		Color          string    `form:"color" json:"color" example:"brown"`
		Gender         string    `form:"gender" json:"gender" example:"male"`
		Description    string    `form:"description" json:"description" example:"Description of the post"`
		Status         string    `form:"status" json:"status" example:"lost"`
		IsFavourite    bool      `form:"is_favourite" json:"is_favourite" example:"true"`
		Comments       int       `form:"comments" json:"comments" example:"2"`
	}

	// UpdatePost, UpdateRequestBodyPost used for updating an existing post
	UpdatePost struct {
		Title   *string `form:"title" json:"title" validate:"max=200"`
		Content *string `form:"content" json:"content" validate:"max=2000"`
		Photos  *[]byte `form:"photos" json:"photos"`
	}

	UpdateRequestBodyPost struct {
		Title       *string   `form:"title" json:"title" validate:"omitempty,max=200"`
		Content     *string   `form:"content" json:"content" validate:"omitempty,max=2000"`
		Photos      *[][]byte `form:"photos" json:"photos"`
		AnimalType  *string   `form:"animal_type" json:"animal_type" validate:"omitempty,oneof=dog cat"`
		Age         *int      `form:"age" json:"age" validate:"omitempty,gte=0"`
		Color       *string   `form:"color" json:"color"`
		Gender      *string   `form:"gender" json:"gender" validate:"omitempty,oneof=male female"`
		Description *string   `form:"description" json:"description"`
		Status      *string   `form:"status" json:"status" validate:"omitempty,oneof=lost found need_home"`
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
		Color      *string `query:"color" validate:"omitempty"`                             // Filter by color of the associated animal
		Location   *string `query:"location" validate:"omitempty"`                          // Filter by location of the associated animal
		SearchWord *string `query:"search_word" validate:"omitempty"`                       // Filter by location of the associated animal
	}

	PathParams struct {
		PostID  int `params:"id" validate:"omitempty,gt=0"`
		PhotoID int `params:"photo_id" validate:"omitempty,gt=0"`
	}
)
