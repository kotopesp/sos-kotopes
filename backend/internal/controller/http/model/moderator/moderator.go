package moderator

import "github.com/kotopesp/sos-kotopes/internal/controller/http/model/post"

type CreateModeratorRequest struct {
	ID int `json:"id" validate:"required,gte=0"`
}

type GetPostsForModerationRequest struct {
	Filter string `query:"filter" validate:"required,oneof=ASC DESC"`
}

type PostsForModerationResponse struct {
	Post    post.PostResponse
	Reasons []string
}

type ModeratedPostRequest struct {
	PostID int `params:"id" validate:"omitempty,gt=0"`
}
