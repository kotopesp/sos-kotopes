package post

import (
	"github.com/kotopesp/sos-kotopes/internal/core"
	modelAnimal "github.com/kotopesp/sos-kotopes/internal/controller/http/model/animal"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/pagination"
)

func ToCorePostDetails(p *CreateRequestBodyPost, authorID int) core.PostDetails {
	if (p == nil) {
		return core.PostDetails{}
	}
	return core.PostDetails{
		Post: ToCorePost(&p.Post, authorID),
		Animal: modelAnimal.ToCoreAnimal(&p.Animal, authorID), 
	}
}

func ToCorePost(post *Post, authorID int) core.Post {
	if (post == nil) {
		return core.Post{}
	}
	return core.Post{
		Title:     post.Title,
		Content:   post.Content,
		Photo:     post.Photo,
		AuthorID:  authorID,
	}
}

func ToResponse(meta pagination.Pagination, posts []core.PostDetails) Response {
	return Response{
		Meta:  meta,
		Posts: posts,
	}
}

func ToPostPesponse(post core.PostDetails) PostResponse {
	return PostResponse{
		Title:          post.Post.Title,
		Content:        post.Post.Content,
		AuthorUsername: post.Username,
		CreatedAt:      post.Post.CreatedAt,
		Animal:         modelAnimal.ToAnimalResponse(&post.Animal),
		Photo:          post.Post.Photo,
		Comments:       0,
	}
}

func ToCorePostFavourite(userID, postID int) core.PostFavourite {
	return core.PostFavourite{
		UserID: userID,
		PostID: postID,
	}
}

func FuncUpdateRequestBodyPost(postDetails core.PostDetails, updateRequestPost UpdateRequestBodyPost) core.PostDetails {
	if updateRequestPost.UpdatePost.Title != nil {
		postDetails.Post.Title = *updateRequestPost.UpdatePost.Title
	}

	if updateRequestPost.UpdatePost.Content != nil {
		postDetails.Post.Content = *updateRequestPost.UpdatePost.Content
	}

	if updateRequestPost.UpdatePost.Photo != nil {
		postDetails.Post.Photo = *updateRequestPost.UpdatePost.Photo
	}

	if updateRequestPost.UpdateAnimal.AnimalType != nil {
		postDetails.Animal.AnimalType = *updateRequestPost.UpdateAnimal.AnimalType
	}

	if updateRequestPost.UpdateAnimal.Age != nil {
		postDetails.Animal.Age = *updateRequestPost.UpdateAnimal.Age
	}

	if updateRequestPost.UpdateAnimal.Color != nil {
		postDetails.Animal.Color = *updateRequestPost.UpdateAnimal.Color
	}

	if updateRequestPost.UpdateAnimal.Gender != nil {
		postDetails.Animal.Gender = *updateRequestPost.UpdateAnimal.Gender
	}

	if updateRequestPost.UpdateAnimal.Description != nil {
		postDetails.Animal.Description = *updateRequestPost.UpdateAnimal.Description
	}
	
	if updateRequestPost.UpdateAnimal.Status != nil {
		postDetails.Animal.Status = *updateRequestPost.UpdateAnimal.Status
	}

	return postDetails
}

func (p *GetAllPostsParams) ToCoreGetAllPostsParams() core.GetAllPostsParams{
	if p == nil {
		return core.GetAllPostsParams{}
	}
	return core.GetAllPostsParams{
		Limit:  	&p.Limit,
		Offset: 	&p.Offset,
		Status: 	&p.Status,
		AnimalType: &p.AnimalType,
		Gender: 	&p.Gender,
		Color:  	&p.Color,
	}
}
