package userwithroles

import (
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/keeper"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/seeker"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/user"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/vet"
)

type (
	UserWithRoles struct {
		User   *user.User     `json:"user"`
		Keeper *keeper.Keeper `json:"keeper"`
		Seeker *seeker.Seeker `json:"seeker"`
		Vet    *vet.Vet       `json:"vet"`
	}
)
