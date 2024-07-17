package keeper

import "gitflic.ru/spbu-se/sos-kotopes/internal/core"

func (k *Keeper) ToCoreKeeper() *core.Keeper {
	if k == nil {
		return nil
	}
	return &core.Keeper{
		Description: k.Description,
		Location:    k.Location,
	}
}
