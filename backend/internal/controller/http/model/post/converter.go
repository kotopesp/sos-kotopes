package post

import "github.com/kotopesp/sos-kotopes/internal/core"

func ToPost(postDetails *core.PostDetails) Post {
	if postDetails == nil {
		return Post{}
	}
	return Post{
		ID:        postDetails.ID,
		Title:     postDetails.Title,
		Content:   postDetails.Content,
		Username:  postDetails.Username,
		CreatedAt: postDetails.CreatedAt,
		UpdatedAt: postDetails.UpdatedAt,
		AnimalID:  postDetails.AnimalID,
	}
}
