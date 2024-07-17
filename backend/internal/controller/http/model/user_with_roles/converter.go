package userwithroles

import "gitflic.ru/spbu-se/sos-kotopes/internal/core"

func (uwr *UserWithRoles) ToCoreUserWithRoles() *core.UserWithRoles {
	if uwr == nil {
		return nil
	}
	return &core.UserWithRoles{
		User:   uwr.User.ToCoreUser(),
		Seeker: uwr.Seeker.ToCoreSeeker(),
		Keeper: uwr.Keeper.ToCoreKeeper(),
		Vet:    uwr.Vet.ToCoreVet(),
	}
}
