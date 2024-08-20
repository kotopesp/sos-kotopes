package post

import (
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/pagination"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

// ToCorePostDetails converts CreateRequestBodyPost from model to core.PostDetails
func (p *CreateRequestBodyPost) ToCorePostDetails(authorID int) core.PostDetails {
	if p == nil {
		return core.PostDetails{}
	}

	post := core.Post{
		Title:    p.Title,
		Content:  p.Content,
		AuthorID: authorID,
	}

	animal := core.Animal{
		KeeperID:    authorID,
		AnimalType:  p.AnimalType,
		Age:         p.Age,
		Color:       p.Color,
		Gender:      p.Gender,
		Description: p.Description,
		Status:      p.Status,
	}

	photos := []core.Photo{}
	for _, photo := range p.Photos {
		photos = append(photos, core.Photo{
			Photo: photo,
		})
	}	

	return core.PostDetails{
		Post:   post,
		Animal: animal,
		Photos: photos,
	}
}

func (p *UpdateRequestBodyPost) ToCorePostDetails() (core.UpdateRequestPost, core.UpdateRequestPhotos) {
	if p == nil {
		return core.UpdateRequestPost{}, core.UpdateRequestPhotos{}
	}

	UpdateRequestPost := core.UpdateRequestPost{
		Title:       p.Title,
		Content:     p.Content,
		AnimalType:  p.AnimalType,
		Age:         p.Age,
		Color:       p.Color,
		Gender:      p.Gender,
		Description: p.Description,
		Status:      p.Status,
	}

	UpdateRequestPhotos := core.UpdateRequestPhotos{
		Photo: p.Photo,
	}

	return UpdateRequestPost, UpdateRequestPhotos
}

// ToResponse converts a list of core.PostDetails to Response with pagination meta
func ToResponse(meta pagination.Pagination, posts []core.PostDetails) Response {
	res := make([]PostResponse, len(posts))

	for i, post := range posts {
		res[i] = ToPostResponse(post)
	}

	return Response{
		Meta:  meta,
		Posts: res,
	}
}

// ToPostResponse converts core.PostDetails to PostResponse
func ToPostResponse(post core.PostDetails) PostResponse {
	var photos [][]byte

	for _, photo := range post.Photos {
		photos = append(photos, photo.Photo)
	}

	return PostResponse{
		Title:          post.Post.Title,
		Content:        post.Post.Content,
		AuthorUsername: post.Username,
		CreatedAt:      post.Post.CreatedAt,
		AnimalType:     post.Animal.AnimalType,
		Age:            post.Animal.Age,
		Color:          post.Animal.Color,
		Gender:         post.Animal.Gender,
		Description:    post.Animal.Description,
		Status:         post.Animal.Status,
		Photo:          photos,
		IsFavourite:    post.Post.IsFavourite,
		Comments:       0,
	}
}

// ToCorePostFavourite converts user ID and post ID to core.PostFavourite
func ToCorePostFavourite(userID, postID int) core.PostFavourite {
	return core.PostFavourite{
		UserID: userID,
		PostID: postID,
	}
}

// ToCoreGetAllPostsParams converts GetAllPostsParams from model to core.GetAllPostsParams
func (p *GetAllPostsParams) ToCoreGetAllPostsParams() core.GetAllPostsParams {
	if p == nil {
		return core.GetAllPostsParams{}
	}
	return core.GetAllPostsParams{
		Limit:      &p.Limit,
		Offset:     &p.Offset,
		Status:     p.Status,
		AnimalType: p.AnimalType,
		Gender: 	p.Gender,
		Color:  	p.Color,
		SearchWord: p.SearchWord,
	}
}
