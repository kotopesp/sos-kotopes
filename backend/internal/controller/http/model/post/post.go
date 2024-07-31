package post

import (
	"time"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/animal"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/pagination"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

type (
	Post struct {
		Title     string    `form:"title" json:"title" validate:"required,max=200"`
		Content   string    `form:"content" json:"content" validate:"required,max=2000"`
		Photo     []byte    `form:"photo" json:"photo"`
	}

	CreateRequestBodyPost struct {
		Post Post
		Animal animal.Animal
	}

	PostResponse struct {
		Title          string 	 		 	 `form:"title" json:"title"`
		Content        string    		 	 `form:"content" json:"content"`
		AuthorUsername string 	 		 	 `form:"author_username" json:"author_username"`
		CreatedAt      time.Time 			 `form:"created_at " json:"created_at"`
		Photo          []byte   			 `form:"photo" json:"photo"`
		Animal         animal.AnimalResponse `form:"animal" json:"animal"`
		Comments       int 				 	 `form:"comments" json:"comments"`
	}

	UpdatePost struct {
		Title   *string `form:"title" json:"title" validate:"max=200"`
		Content *string `form:"content" json:"content" validate:"max=2000"`
		Photo   *[]byte `form:"photo" json:"photo"`
	}

	UpdateRequestBodyPost struct {
		UpdatePost UpdatePost
		UpdateAnimal animal.UpdateAnimal
	}

	Meta struct {
		Total       int `json:"total"`
		TotalPages  int `json:"total_pages"`
		CurrentPage int `json:"current_page"`
		PerPage     int `json:"per_page"`
	}
	
	Response struct {
		Meta  pagination.Pagination `json:"meta"`
		Posts []core.PostDetails 		`json:"posts"`
	}

	GetAllPostsParams struct {
		Limit      int    `query:"limit" validate:"gt=0"`
		Offset     int    `query:"offset" validate:"gte=0"`
		Status     string `query:"status" validate:"one_of=lost,found,need_home"`
		AnimalType string `query:"animal_type" validate:"one_of=dog,cat"`
		Gender     string `query:"gender" validate:"one_of=male,female"`
		Color      string `query:"color"`
		Location   string `query:"location"`
	}
)
