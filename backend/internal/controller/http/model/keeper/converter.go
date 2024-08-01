package keeper

import (
	"strings"

	"github.com/kotopesp/sos-kotopes/internal/core"
)

func (p *GetAllKeepersParams) FromKeeperRequest() core.GetAllKeepersParams {
	sortBy, sortOrder := p.ParseSort()
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

func (p *GetAllKeepersParams) ParseSort() (string, string) {
	parts := strings.Split(p.Sort, ":")
	sortBy := ""
	sortOrder := ""
	if len(parts[0]) > 0 {
		sortBy = parts[0]
	}
	if len(parts[1]) > 0 {
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

func FromCoreKeeperReview(coreKeeper core.Keepers) KeepersResponse {
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
