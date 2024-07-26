package comments

import "gitflic.ru/spbu-se/sos-kotopes/internal/core"

func (c *Comments) ToCoreComments() *core.Comments {
	if c == nil {
		//TODO: доделать
	}
	return &core.Comments{
		Id:         c.ID,
		Content:    c.Content,
		Author_id:  c.Author_id,
		Posts_id:   c.Posts_id,
		Is_deleted: c.Is_deleted,
		Deleted_at: c.Deleted_at,
		Created_at: c.Created_at,
		Updated_at: c.Updated_at,
		Parent_id:  c.Parent_id,
		Reply_id:   c.Reply_id,
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
