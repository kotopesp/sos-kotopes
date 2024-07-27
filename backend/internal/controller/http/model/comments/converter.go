package comments

import "github.com/kotopesp/sos-kotopes/internal/core"

func (c *Comments) ToCoreComments() *core.Comments {
	if c == nil {
		//TODO: доделать
	}
	return &core.Comments{
		ID:        c.ID,
		Content:   c.Content,
		AuthorID:  c.AuthorID,
		PostsID:   c.PostsID,
		IsDeleted: c.IsDeleted,
		DeletedAt: c.DeletedAt,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		ParentID:  c.ParentID,
		ReplyID:   c.ReplyID,
	}
}
func (c *GetAllCommentsParams) ToCoreGetAllCommentsParams() *core.GetAllParamsComments {
	if c == nil {
		//TODO: доделать
	}
	return &core.GetAllParamsComments{
		Limit:  &c.Limit,
		Offset: &c.Offset,
	}
}
