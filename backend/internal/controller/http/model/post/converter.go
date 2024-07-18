package post

import (
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
)

func (p *GetAllPostsParams) FromUserRequest() *core.GetAllPostsParams {
	if p == nil {
		return nil
	}

	return &core.GetAllPostsParams{
		SortBy:     &p.SortBy,
		SortOrder:  &p.SortOrder,
		SearchTerm: &p.SearchTerm,
		Limit:      &p.Limit,
		Offset:     &p.Offset,
	}
}