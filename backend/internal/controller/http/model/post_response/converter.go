package response

import (
	"github.com/kotopesp/sos-kotopes/internal/core"
)

func (r *Response) ToCoreResponse() *core.PostResponse {
	if r == nil {
		return &core.PostResponse{}
	}
	return &core.PostResponse{
		ID:        r.ID,
		PostID:    r.PostID,
		AuthorID:  r.AuthorID,
		Content:   r.Content,
		IsDeleted: r.IsDeleted,
		DeletedAt: r.DeletedAt,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
