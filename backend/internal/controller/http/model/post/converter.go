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
		Photo:       p.Photo,
		AuthorID:    authorID,
		Description: p.Description,
		Location:   p.Location,
	}

	animal := core.Animal{
		UserID:     authorID,
		AnimalType: p.AnimalType,
		Color:      p.Color,
		Gender:     p.Gender,
		Status:     p.Status,
	}

	return core.PostDetails{
		Post:   post,
		Animal: animal,
	}
}

func (p *UpdateRequestBodyPost) ToCorePostDetails() core.UpdateRequestBodyPost {
	if p == nil {
		return core.UpdateRequestBodyPost{}
	}

	return core.UpdateRequestBodyPost{
		Title:       p.Title,
		Content:     p.Content,
		Photo:       p.Photo,
		AnimalType:  p.AnimalType,
		Age:         p.Age,
		Color:       p.Color,
		Gender:      p.Gender,
		Description: p.Description,
		Status:      p.Status,
	}
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
	return PostResponse{
		AuthorUsername: post.Username,
		CreatedAt:      post.Post.CreatedAt,
		AnimalType:     post.Animal.AnimalType,
		Color:          post.Animal.Color,
		Gender:         post.Animal.Gender,
		Description:    post.Post.Description,
		Status:         post.Animal.Status,
		Photo:          post.Post.Photo,
		Location:       post.Post.Location,
		IsFavourite:    false,
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
		Gender:     p.Gender,
		Color:      p.Color,
	}
}
