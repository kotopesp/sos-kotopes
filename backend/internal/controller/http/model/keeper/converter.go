package keeper

import (
	"strings"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/pagination"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

func (p *GetAllKeepersParams) FromKeeperRequest() core.GetAllKeepersParams {
	sortBy, sortOrder := ParseSort(p.Sort)
	return core.GetAllKeepersParams{
		SortBy:    &sortBy,
		SortOrder: &sortOrder,
		Location:  &p.Location,
		MinRating: &p.MinRating,
		MaxRating: &p.MaxRating,
		MinPrice:  &p.MinPrice,
		MaxPrice:  &p.MaxPrice,
		Limit:     &p.Limit,
		Offset:    &p.Offset,
	}
}

func ParseSort(sort string) (sortBy, sortOrder string) {
	parts := strings.Split(sort, ":")
	if len(parts) != 2 {
		return "", ""
	}
	sortBy = ""
	sortOrder = ""
	if parts[0] != "" {
		sortBy = parts[0]
	}
	if parts[1] != "" {
		sortOrder = parts[1]
	}
	return sortBy, sortOrder
}

func (k *KeepersCreate) ToCoreNewKeeper() core.Keepers {
	if k == nil {
		return core.Keepers{}
	}
	return core.Keepers{
		UserID:      k.UserID,
		Description: k.Description,
		Price:       k.Price,
		Location:    k.Location,
	}
}

func (k *KeepersUpdate) ToCoreUpdatedKeeper() core.Keepers {
	if k == nil {
		return core.Keepers{}
	}
	return core.Keepers{
		ID:          k.ID,
		Description: k.Description,
		Price:       k.Price,
		Location:    k.Location,
	}
}

func ToKeepersResponse(meta pagination.Pagination, coreKeepers []core.Keepers) KeepersResponseWithMeta {
	offset := (meta.CurrentPage - 1) * meta.PerPage
	paginateCoreKeepers := coreKeepers[offset:min(offset+meta.PerPage, meta.Total)]
	paginateResponseKeepers := make([]KeepersResponse, meta.PerPage)

	for i, coreKeeper := range paginateCoreKeepers {
		paginateResponseKeepers[i] = FromCoreKeeper(coreKeeper)
	}

	return KeepersResponseWithMeta{
		Meta: meta,
		Data: paginateResponseKeepers,
	}
}

func FromCoreKeeper(coreKeeper core.Keepers) KeepersResponse {
	return KeepersResponse{
		ID:          coreKeeper.ID,
		UserID:      coreKeeper.UserID,
		Description: coreKeeper.Description,
		Price:       coreKeeper.Price,
		Location:    coreKeeper.Location,
		CreatedAt:   coreKeeper.CreatedAt,
		UpdatedAt:   coreKeeper.UpdatedAt,
	}
}
