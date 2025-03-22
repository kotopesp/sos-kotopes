package moderator

import (
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
