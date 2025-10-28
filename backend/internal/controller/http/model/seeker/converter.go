package seeker

import (
	"github.com/kotopesp/sos-kotopes/internal/core"
)

func (seeker *CreateSeeker) ToCoreSeeker() core.Seeker {
	return core.Seeker{
		AnimalType:       seeker.AnimalType,
		Description:      seeker.Description,
		LocationID:       seeker.LocationId,
		EquipmentRental:  seeker.EquipmentRental,
		HaveMetalCage:    seeker.HaveMetalCage,
		HavePlasticCage:  seeker.HavePlasticCage,
		HaveNet:          seeker.HaveNet,
		HaveLadder:       seeker.HaveLadder,
		HaveOther:        seeker.HaveOther,
		HaveCar:          seeker.HaveCar,
		Price:            seeker.Price,
		WillingnessCarry: seeker.WillingnessCarry,
	}
}

func (seeker *UpdateSeeker) ToCoreUpdateSeeker() core.UpdateSeeker {
	return core.UpdateSeeker{
		AnimalType:       seeker.AnimalType,
		Description:      seeker.Description,
		LocationID:       seeker.LocationId,
		EquipmentRental:  seeker.EquipmentRental,
		HaveMetalCage:    seeker.HaveMetalCage,
		HavePlasticCage:  seeker.HavePlasticCage,
		HaveNet:          seeker.HaveNet,
		HaveLadder:       seeker.HaveLadder,
		HaveOther:        seeker.HaveOther,
		HaveCar:          seeker.HaveCar,
		Price:            seeker.Price,
		WillingnessCarry: seeker.WillingnessCarry,
	}
}

func (p *GetAllSeekerParams) ToCoreGetAllSeekersParams() core.GetAllSeekersParams {
	return core.GetAllSeekersParams{
		SortBy:             p.SortBy,
		SortOrder:          p.SortOrder,
		AnimalType:         p.AnimalType,
		LocationID:         p.LocationId,
		MinPrice:           p.MinPrice,
		MaxPrice:           p.MaxPrice,
		MinEquipmentRental: p.MinEquipmentRental,
		MaxEquipmentRental: p.MaxEquipmentRental,
		HaveMetalCage:      p.HaveMetalCage,
		HavePlasticCage:    p.HavePlasticCage,
		HaveNet:            p.HaveNet,
		HaveLadder:         p.HaveLadder,
		HaveOther:          p.HaveOther,
		HaveCar:            p.HaveCar,
		Limit:              p.Limit,
		Offset:             p.Offset,
	}
}

func ToResponseSeeker(seeker *core.Seeker) ResponseSeeker {
	responseSeeker := ResponseSeeker{}

	if seeker.User.Firstname != nil {
		responseSeeker.Firstname = *seeker.User.Firstname
	}

	if seeker.User.Lastname != nil {
		responseSeeker.Lastname = *seeker.User.Lastname
	}

	if seeker.User.Photo != nil {
		responseSeeker.Photo = *seeker.User.Photo
	}

	responseSeeker.ID = seeker.ID
	responseSeeker.UserID = seeker.UserID
	responseSeeker.AnimalType = seeker.AnimalType
	responseSeeker.LocationId = seeker.LocationID
	responseSeeker.EquipmentRental = seeker.EquipmentRental
	responseSeeker.HaveMetalCage = seeker.HaveMetalCage
	responseSeeker.HavePlasticCage = seeker.HavePlasticCage
	responseSeeker.HaveNet = seeker.HaveNet
	responseSeeker.HaveLadder = seeker.HaveLadder
	responseSeeker.HaveOther = seeker.HaveOther
	responseSeeker.Description = seeker.Description
	responseSeeker.HaveCar = seeker.HaveCar
	responseSeeker.Price = seeker.Price
	responseSeeker.WillingnessCarry = seeker.WillingnessCarry

	return responseSeeker
}

func ToResponseSeekers(coreSeekers []core.Seeker) ResponseSeekers {
	responseSeekers := make([]ResponseSeeker, len(coreSeekers))

	for i := range coreSeekers {
		responseSeekers[i] = ToResponseSeeker(&coreSeekers[i])
	}

	return ResponseSeekers{
		ResponseSeekers: responseSeekers,
	}
}
