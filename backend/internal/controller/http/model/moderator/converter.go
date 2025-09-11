package moderator

import (
	"time"

	post "github.com/kotopesp/sos-kotopes/internal/controller/http/model/post"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

func (c CreateModeratorRequest) ToCoreModerator() core.Moderator {
	return core.Moderator{
		UserID: c.ID,
	}
}

func ToPostList(posts []core.PostForModeration) []core.Post {
	postList := make([]core.Post, len(posts))

	for i, p := range posts {
		postList[i] = p.Post
	}

	return postList
}

func ToPostsForModerationResponse(postsAndReasons []core.PostForModeration, details []core.PostDetails) (response []PostsForModerationResponse) {
	for i, postWithReason := range postsAndReasons {
		response = append(response, PostsForModerationResponse{
			Post:    post.ToPostResponse(details[i]),
			Reasons: postWithReason.Reasons,
		})
	}

	return response
}

func ToCommentsForModerationResponse(comments []core.CommentForModeration) []CommentsForModerationResponse {
	var response []CommentsForModerationResponse
	for _, c := range comments {
		response = append(response, CommentsForModerationResponse{
			CommentID: c.Comment.ID,
			Content:   c.Comment.Content,
			PostID:    c.Comment.PostID,
			AuthorID:  c.Comment.AuthorID,
			CreatedAt: c.Comment.CreatedAt.Format(time.RFC3339),
			Reasons:   c.Reasons,
		})
	}
	return response
}
