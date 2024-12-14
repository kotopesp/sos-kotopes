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
		Location:             p.Location,
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

func (k *KeepersCreate) ToCoreNewKeeper() core.Keepers {
	if k == nil {
		return core.Keepers{}
	}
	return core.Keepers{
		UserID:               k.UserID,
		Description:          k.Description,
		Price:                k.Price,
		Location:             k.Location,
		HasCage:              k.HasCage,
		BoardingDuration:     k.BoardingDuration,
		BoardingCompensation: k.BoardingCompensation,
		AnimalAcceptance:     k.AnimalAcceptance,
		AnimalCategory:       k.AnimalCategory,
	}
}

func (k *KeepersUpdate) ToCoreUpdatedKeeper() core.UpdateKeepers {
	if k == nil {
		return core.UpdateKeepers{}
	}
	return core.UpdateKeepers{
		ID:                   k.ID,
		UserID:               k.UserID,
		Description:          k.Description,
		Price:                k.Price,
		Location:             k.Location,
		HasCage:              k.HasCage,
		BoardingDuration:     k.BoardingDuration,
		BoardingCompensation: k.BoardingCompensation,
		AnimalAcceptance:     k.AnimalAcceptance,
		AnimalCategory:       k.AnimalCategory,
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
		ID:                   coreKeeper.ID,
		UserID:               coreKeeper.UserID,
		Description:          coreKeeper.Description,
		Price:                coreKeeper.Price,
		Location:             coreKeeper.Location,
		HasCage:              coreKeeper.HasCage,
		BoardingDuration:     coreKeeper.BoardingDuration,
		BoardingCompensation: coreKeeper.BoardingCompensation,
		AnimalAcceptance:     coreKeeper.AnimalAcceptance,
		AnimalCategory:       coreKeeper.AnimalCategory,
		CreatedAt:            coreKeeper.CreatedAt,
		UpdatedAt:            coreKeeper.UpdatedAt,
	}
}

func FromCoreKeeperDetails(coreKeeperDetails core.KeepersDetails) KeepersResponseWithUser {
	return KeepersResponseWithUser{
		Keeper: FromCoreKeeper(coreKeeperDetails.Keeper),
		User:   user.ToResponseUser(&coreKeeperDetails.User),
	}
}
