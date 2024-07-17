package vet

import "gitflic.ru/spbu-se/sos-kotopes/internal/core"

func (v *Vet) ToCoreVet() *core.Vet {
	if v == nil {
		return nil
	}
	return &core.Vet{
		Description: v.Description,
		Location:    v.Location,
	}
}
