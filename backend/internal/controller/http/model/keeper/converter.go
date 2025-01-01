package keeper

import (
	"strings"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/pagination"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/user"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

func (p *GetAllKeepersParams) ToCoreGetAllKeepersParams() core.GetAllKeepersParams {
	sortBy, sortOrder := "", ""
	if p.Sort != nil {
		sortBy, sortOrder = ParseSort(*p.Sort)
	}
	limit, offset := 10, 0
	if p.Limit != nil {
		limit = *p.Limit
	}
	if p.Offset != nil {
		offset = *p.Offset
	}

	return core.GetAllKeepersParams{
		SortBy:               &sortBy,
		SortOrder:            &sortOrder,
		LocationID:           p.LocationID,
		MinRating:            p.MinRating,
		MaxRating:            p.MaxRating,
		MinPrice:             p.MinPrice,
		MaxPrice:             p.MaxPrice,
		HasCage:              p.HasCage,
		BoardingDuration:     p.BoardingDuration,
		BoardingCompensation: p.BoardingCompensation,
		AnimalAcceptance:     p.AnimalAcceptance,
		AnimalCategory:       p.AnimalCategory,
		Limit:                &limit,
		Offset:               &offset,
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

func (k *CreateKeeper) ToCoreKeeper() core.Keeper {
	if k == nil {
		return core.Keeper{}
	}
	return core.Keeper{
		UserID:               k.UserID,
		Description:          k.Description,
		Price:                k.Price,
		LocationID:           k.LocationID,
		HasCage:              k.HasCage,
		BoardingDuration:     k.BoardingDuration,
		BoardingCompensation: k.BoardingCompensation,
		AnimalAcceptance:     k.AnimalAcceptance,
		AnimalCategory:       k.AnimalCategory,
	}
}

func (k *UpdateKeeper) ToCoreUpdateKeeper() core.Keeper {
	if k == nil {
		return core.Keeper{}
	}
	return core.Keeper{
		Description:          k.Description,
		Price:                k.Price,
		LocationID:           k.LocationID,
		HasCage:              *k.HasCage,
		BoardingDuration:     *k.BoardingDuration,
		BoardingCompensation: *k.BoardingCompensation,
		AnimalAcceptance:     *k.AnimalAcceptance,
		AnimalCategory:       *k.AnimalCategory,
	}
}

func ToModelResponseKeepers(meta pagination.Pagination, coreKeepers []core.Keeper) ResponseKeepers {
	offset := (meta.CurrentPage - 1) * meta.PerPage
	paginateCoreKeepers := coreKeepers[offset:max(0, min(offset+meta.PerPage, meta.Total))]
	paginateKeepersResponse := make([]ResponseKeeper, len(paginateCoreKeepers))

	for i, coreKeeper := range paginateCoreKeepers {
		paginateKeepersResponse[i] = ToModelResponseKeeper(coreKeeper)
	}

	return ResponseKeepers{
		Meta: meta,
		Data: paginateKeepersResponse,
	}
}

func ToModelResponseKeeper(coreKeeper core.Keeper) ResponseKeeper {
	return ResponseKeeper{
		ID:                   coreKeeper.ID,
		UserID:               coreKeeper.UserID,
		User:                 user.ToResponseUser(&coreKeeper.User),
		Description:          coreKeeper.Description,
		Price:                coreKeeper.Price,
		LocationID:           coreKeeper.LocationID,
		HasCage:              coreKeeper.HasCage,
		BoardingDuration:     coreKeeper.BoardingDuration,
		BoardingCompensation: coreKeeper.BoardingCompensation,
		AnimalAcceptance:     coreKeeper.AnimalAcceptance,
		AnimalCategory:       coreKeeper.AnimalCategory,
		CreatedAt:            coreKeeper.CreatedAt,
		UpdatedAt:            coreKeeper.UpdatedAt,
	}
}
