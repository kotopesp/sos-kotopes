package keeper

import (
	"strings"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/pagination"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/user"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

func (p *GetAllKeepersParams) ToCoreGetAllKeepersParams() core.GetAllKeepersParams {
	sortBy, sortOrder := ParseSort(p.Sort)
	return core.GetAllKeepersParams{
		SortBy:    &sortBy,
		SortOrder: &sortOrder,
		Location:  p.Location,
		MinRating: p.MinRating,
		MaxRating: p.MaxRating,
		MinPrice:  p.MinPrice,
		MaxPrice:  p.MaxPrice,
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

func (k *KeepersUpdate) ToCoreUpdatedKeeper() core.UpdateKeepers {
	if k == nil {
		return core.UpdateKeepers{}
	}
	return core.UpdateKeepers{
		ID:          k.ID,
		UserID:      k.UserID,
		Description: k.Description,
		Price:       k.Price,
		Location:    k.Location,
	}
}

func ToKeepersResponse(meta pagination.Pagination, coreKeepersDetails []core.KeepersDetails) KeepersResponseWithMeta {
	offset := (meta.CurrentPage - 1) * meta.PerPage
	paginateCoreKeepersDetails := coreKeepersDetails[offset:min(offset+meta.PerPage, meta.Total)]
	paginateResponseKeepersWithUser := make([]KeepersResponseWithUser, len(paginateCoreKeepersDetails))

	for i, coreKeeperDetails := range paginateCoreKeepersDetails {
		paginateResponseKeepersWithUser[i] = FromCoreKeeperDetails(coreKeeperDetails)
	}

	return KeepersResponseWithMeta{
		Meta: meta,
		Data: paginateResponseKeepersWithUser,
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

func FromCoreKeeperDetails(coreKeeperDetails core.KeepersDetails) KeepersResponseWithUser {
	return KeepersResponseWithUser{
		Keeper: FromCoreKeeper(coreKeeperDetails.Keeper),
		User:   user.ToResponseUser(&coreKeeperDetails.User),
	}
}
