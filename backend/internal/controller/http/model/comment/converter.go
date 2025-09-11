package comment

import (
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/pagination"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

func (c *Comment) ToCoreComment() core.Comment {
	return core.Comment{
		Content:  c.Content,
		ParentID: c.ParentID,
		ReplyID:  c.ReplyID,
	}
}

func (c *Update) ToCoreComment() core.Comment {
	return core.Comment{
		Content: c.Content,
	}
}

func (c *Create) ToCoreComment() core.Comment {
	return core.Comment{
		Content:  c.Content,
		ParentID: c.ParentID,
		ReplyID:  c.ReplyID,
	}
}

func ToModelComment(c core.Comment) Comment {
	if c.Status == core.Deleted {
		c.Content = ""
	}

	return Comment{
		ID:       c.ID,
		ParentID: c.ParentID,
		ReplyID:  c.ReplyID,
		User: User{
			ID:       c.Author.ID,
			Username: c.Author.Username,
		},
		Content:   c.Content,
		Status:    string(c.Status),
		CreatedAt: c.CreatedAt,
	}
}

func ToModelCommentsSlice(c []core.Comment) []Comment {
	modelCommentsSlice := make([]Comment, len(c))
	for i, comment := range c {
		modelCommentsSlice[i] = ToModelComment(comment)
	}
	return modelCommentsSlice
}

func (params *GetAllCommentsParams) ToCoreGetAllCommentsParams(postID int) core.GetAllCommentsParams {
	return core.GetAllCommentsParams{
		PostID: postID,
		Limit:  &params.Limit,
		Offset: &params.Offset,
	}
}

func ToGetAllCommentsResponse(data []Comment, meta pagination.Pagination) GetAllCommentsResponse {
	return GetAllCommentsResponse{
		Data: data,
		Meta: meta,
	}
}
