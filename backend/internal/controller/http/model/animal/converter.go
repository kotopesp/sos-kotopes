package animal

import (
	"github.com/kotopesp/sos-kotopes/internal/core"
)

func(a *Animal) ToCoreAnimal(keeperID int) core.Animal {
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

func FuncUpdateRequestBodyAnimal(animal *core.Animal, updateRequestAnimal *UpdateRequestBodyAnimal) core.Animal {
	if updateRequestAnimal.AnimalType != nil {
		animal.AnimalType = *updateRequestAnimal.AnimalType
	}
	if updateRequestAnimal.Age != nil {
		animal.Age = *updateRequestAnimal.Age
	}
	if updateRequestAnimal.Color != nil {
		animal.Color = *updateRequestAnimal.Color
	}
	if updateRequestAnimal.Gender != nil {
		animal.Gender = *updateRequestAnimal.Gender
	}
	if updateRequestAnimal.Description != nil {
		animal.Description = *updateRequestAnimal.Description
	}
	if updateRequestAnimal.Status != nil {
		animal.Status = *updateRequestAnimal.Status
	}

	return *animal
}
