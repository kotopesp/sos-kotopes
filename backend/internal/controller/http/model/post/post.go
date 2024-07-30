package post

import (
	"time"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/animal"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/pagination"
)

type (
	Post struct {
		Title     string    `form:"title" json:"title" validate:"required"`
		Content   string    `form:"content" json:"content" validate:"required"`
		Photo     []byte    `form:"photo" json:"photo"`
	}

	PostPesponse struct {
		Title          string 	 		 	 `form:"title" json:"title"`
		Content        string    		 	 `form:"content" json:"content"`
		AuthorUsername string 	 		 	 `form:"author_username" json:"author_username"`
		CreatedAt      time.Time 			 `form:"created_at " json:"created_at"`
		Photo          []byte   			 `form:"photo" json:"photo"`
		Animal         animal.AnimalResponse `form:"animal" json:"animal"`
	}

	UpdateRequestBodyPost struct {
		Title   *string `form:"title" json:"title"`
		Content *string `form:"content" json:"content"`
		Photo   *[]byte `form:"photo" json:"photo"`
	}

	Meta struct {
		Total       int `json:"total"`
		TotalPages  int `json:"total_pages"`
		CurrentPage int `json:"current_page"`
		PerPage     int `json:"per_page"`
	}
	
	Response struct {
		Meta  pagination.Pagination `json:"meta"`
		Posts []PostPesponse 		`json:"posts"`
	}

	GetAllPostsParams struct {
		Limit      int    `query:"limit" validate:"gt=0"`
		Offset     int    `query:"offset" validate:"gte=0"`
	}
)
