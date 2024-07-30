package comments

import "github.com/kotopesp/sos-kotopes/internal/core"

func (c *Comments) ToCoreComments() core.Comments {
	return core.Comments{
		Content:  c.Content,
		ParentID: c.ParentID,
		ReplyID:  c.ReplyID,
	}
}

func (c *CommentUpdate) ToCoreComment() core.Comments {
	return core.Comments{
		Content: c.Content,
	}
}

func ToModelComment(c core.Comments) Comments {
	return Comments{
		ID:        c.ID,
		Content:   c.Content,
		ParentID:  c.ParentID,
		ReplyID:   c.ReplyID,
		UpdatedAt: c.UpdatedAt,
		IsDeleted: c.IsDeleted,
	}
}

func ToModelCommentsSlice(c []core.Comments) []Comments {
	modelCommentsSlice := make([]Comments, len(c))
	for i, comment := range c {
		modelCommentsSlice[i] = ToModelComment(comment)
	}
	return modelCommentsSlice
}

func (c *GetAllCommentsParams) ToCoreGetAllCommentsParams() core.GetAllParamsComments {
	return core.GetAllParamsComments{
		Limit:  &c.Limit,
		Offset: &c.Offset,
	}
}
