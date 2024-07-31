package animal

import (
	"github.com/kotopesp/sos-kotopes/internal/core"
)

func ToCoreAnimal(a *Animal, keeperID int) core.Animal {
	if a == nil {
		return core.Animal{}
	}
	return core.Animal{
		KeeperID:    keeperID,
		AnimalType:  a.AnimalType,
		Age:         a.Age,
		Color:       a.Color,
		Gender:      a.Gender,
		Description: a.Description,
		Status:      a.Status,
	}
}

func ToAnimalResponse(a *core.Animal) AnimalResponse {
	if a == nil {
		return AnimalResponse{}
	}
	return AnimalResponse{
		AnimalType:  a.AnimalType,
		Age:         a.Age,
		Color:       a.Color,
		Gender:      a.Gender,
		Description: a.Description,
		Status:      a.Status,
	}
}
