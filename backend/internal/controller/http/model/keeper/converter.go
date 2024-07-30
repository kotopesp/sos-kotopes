package keeper

import "github.com/kotopesp/sos-kotopes/internal/core"

func (p *GetAllKeepersParams) FromKeeperRequest() core.GetAllKeepersParams {
	return core.GetAllKeepersParams{
		SortBy:    &p.SortBy,
		SortOrder: &p.SortOrder,
		Location:  &p.Location,
		MinRating: &p.MinRating,
		MaxRating: &p.MaxRating,
		MinPrice:  &p.MinPrice,
		MaxPrice:  &p.MaxPrice,
		Limit:     &p.Limit,
		Offset:    &p.Offset,
	}
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
