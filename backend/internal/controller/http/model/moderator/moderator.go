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

type GetCommentsForModerationRequest struct {
	Filter string `query:"filter" validate:"required,oneof=ASC DESC"`
}

type ModeratedCommentRequest struct {
	CommentID int `params:"id" validate:"required,min=1"`
}

type CommentsForModerationResponse struct {
	CommentID int      `json:"comment_id"`
	Content   string   `json:"content"`
	PostID    int      `json:"post_id"`
	AuthorID  int      `json:"author_id"`
	CreatedAt string   `json:"created_at"`
	Reasons   []string `json:"reasons"`
}
