package post

import (
	"github.com/kotopesp/sos-kotopes/internal/core"
	modelAnimal "github.com/kotopesp/sos-kotopes/internal/controller/http/model/animal"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/pagination"
)
func (p *Post) ToCorePost(authorID int) core.Post {
	if (p == nil) {
		return core.Post{}
	}
	return core.Post{
		Title:     p.Title,
		Content:   p.Content,
		AuthorID:  authorID,
	}
}

func ToResponse(meta pagination.Pagination, posts []PostPesponse) Response {
	return Response{
		Meta:  meta,
		Posts: posts,
	}
}

func ToPostPesponse(authorUsername string, post core.Post, animal core.Animal) PostPesponse {
	return PostPesponse{
		Title:          post.Title,
		Content:        post.Content,
		AuthorUsername: authorUsername,
		CreatedAt:      post.CreatedAt,
		Animal:         modelAnimal.ToAnimalResponse(&animal),
		Photo:          post.Photo,
	}
}

func ToCorePostFavourite(userID, postID int) core.PostFavourite {
	return core.PostFavourite{
		UserID: userID,
		PostID: postID,
	}
}

func FuncUpdateRequestBodyPost(post *core.Post, updateRequestPost *UpdateRequestBodyPost) core.Post {
	if updateRequestPost.Title != nil {
		post.Title = *updateRequestPost.Title
	}
	if updateRequestPost.Content != nil {
		post.Content = *updateRequestPost.Content
	}
	if updateRequestPost.Photo != nil {
		post.Photo = *updateRequestPost.Photo
	}

	return *post
}
