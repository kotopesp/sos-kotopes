package post

import (
	"github.com/kotopesp/sos-kotopes/internal/core"
	modelAnimal "github.com/kotopesp/sos-kotopes/internal/controller/http/model/animal"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/pagination"
)

func ToSplitPostAndAnimal(p CreateRequestBodyPost) (Post, modelAnimal.Animal) {
	post := Post{
		Title:   p.Title,
		Content: p.Content,
		Photo:   p.Photo,
	}

	animal := modelAnimal.Animal{
		AnimalType: p.AnimalType,
		Age:        p.Age,
		Color:      p.Color,
		Gender:     p.Gender,
		Description: p.Description,
		Status:     p.Status,
	}
	return post, animal
}

func ToCorePostDetails(p *CreateRequestBodyPost, authorID int) core.PostDetails {
	if (p == nil) {
		return core.PostDetails{}
	}

	post, animal := ToSplitPostAndAnimal(*p)
	return core.PostDetails{
		Post: ToCorePost(&post, authorID),
		Animal: modelAnimal.ToCoreAnimal(&animal, authorID), 
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
	res := make([]PostResponse, len(posts))

	for i, post := range posts {
		res[i] = ToPostPesponse(post)
	}

	return Response{
		Meta:  meta,
		Posts: res,
	}
}

func ToPostPesponse(post core.PostDetails) PostResponse {
	return PostResponse{
		Title:          post.Post.Title,
		Content:        post.Post.Content,
		AuthorUsername: post.Username,
		CreatedAt:      post.Post.CreatedAt,
		AnimalType:  	post.Animal.AnimalType,
		Age:         	post.Animal.Age,
		Color:       	post.Animal.Color,
		Gender:      	post.Animal.Gender,
		Description: 	post.Animal.Description,
		Status:      	post.Animal.Status,
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

func ToSplitUpdateRequestBodyPost(updateRequestPost UpdateRequestBodyPost) (UpdatePost,  modelAnimal.UpdateAnimal) {

	post := UpdatePost{
		Title:   updateRequestPost.Title,
		Content: updateRequestPost.Content,
		Photo:   updateRequestPost.Photo,
	}

	animal := modelAnimal.UpdateAnimal{
		AnimalType: updateRequestPost.AnimalType,
		Age:        updateRequestPost.Age,
		Color:      updateRequestPost.Color,
		Gender:     updateRequestPost.Gender,
		Description: updateRequestPost.Description,
		Status:     updateRequestPost.Status,
	}

	return post, animal
}

func FuncUpdateRequestBodyPost(postDetails core.PostDetails, updateRequestPost UpdateRequestBodyPost) core.PostDetails {
	updatePost, updateAnimal := ToSplitUpdateRequestBodyPost(updateRequestPost)

	if updatePost.Title != nil {
		postDetails.Post.Title = *updatePost.Title
	}

	if updatePost.Content != nil {
		postDetails.Post.Content = *updatePost.Content
	}

	if updatePost.Photo != nil {
		postDetails.Post.Photo = *updatePost.Photo
	}

	if updateAnimal.AnimalType != nil {
		postDetails.Animal.AnimalType = *updateAnimal.AnimalType
	}

	if updateAnimal.Age != nil {
		postDetails.Animal.Age = *updateAnimal.Age
	}

	if updateAnimal.Color != nil {
		postDetails.Animal.Color = *updateAnimal.Color
	}

	if updateAnimal.Gender != nil {
		postDetails.Animal.Gender = *updateAnimal.Gender
	}

	if updateAnimal.Description != nil {
		postDetails.Animal.Description = *updateAnimal.Description
	}

	if updateAnimal.Status != nil {
		postDetails.Animal.Status = *updateAnimal.Status
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
		Status: 	p.Status,
		AnimalType: p.AnimalType,
		Gender: 	p.Gender,
		Color:  	p.Color,
	}
}
