package seeker

import (
	"github.com/kotopesp/sos-kotopes/internal/core"
)

func (seeker *CreateSeeker) ToCoreSeeker() core.Seeker {
	return core.Seeker{
		AnimalType:       seeker.AnimalType,
		Description:      seeker.Description,
		Location:         seeker.Location,
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
		Location:         seeker.Location,
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
	sortBy, sortOrder := "", ""

	if p.SortBy != nil {
		sortBy = *p.SortBy
	}

	if p.SortOrder != nil {
		sortBy = *p.SortOrder
	}

	limit, offset := 10, 0

	if p.Limit != nil {
		limit = *p.Limit
	}

	if p.Offset != nil {
		offset = *p.Offset
	}

	return core.GetAllSeekersParams{
		SortBy:             &sortBy,
		SortOrder:          &sortOrder,
		AnimalType:         p.AnimalType,
		Location:           p.Location,
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
		Limit:              &limit,
		Offset:             &offset,
	}
}

func ToResponseSeeker(seeker *core.Seeker) ResponseSeeker {
	if seeker == nil {
		return ResponseSeeker{}
	}
	return ResponseSeeker{
		ID:               seeker.ID,
		UserID:           seeker.UserID,
		AnimalType:       seeker.AnimalType,
		Location:         seeker.Location,
		EquipmentRental:  seeker.EquipmentRental,
		HaveMetalCage:    seeker.HaveMetalCage,
		HavePlasticCage:  seeker.HavePlasticCage,
		HaveNet:          seeker.HaveNet,
		HaveLadder:       seeker.HaveLadder,
		HaveOther:        seeker.HaveOther,
		Description:      seeker.Description,
		HaveCar:          seeker.HaveCar,
		Price:            seeker.Price,
		WillingnessCarry: seeker.WillingnessCarry,
	}
}

func toResponseSeeker(seeker core.Seeker) ResponseSeeker {
	return ResponseSeeker{
		ID:               seeker.ID,
		UserID:           seeker.UserID,
		AnimalType:       seeker.AnimalType,
		Location:         seeker.Location,
		EquipmentRental:  seeker.EquipmentRental,
		HaveMetalCage:    seeker.HaveMetalCage,
		HavePlasticCage:  seeker.HavePlasticCage,
		HaveNet:          seeker.HaveNet,
		HaveLadder:       seeker.HaveLadder,
		HaveOther:        seeker.HaveOther,
		Description:      seeker.Description,
		HaveCar:          seeker.HaveCar,
		Price:            seeker.Price,
		WillingnessCarry: seeker.WillingnessCarry,
	}
}

func ToResponseSeekers(coreSeekers []core.Seeker) ResponseSeekers {
	responseSeekers := make([]ResponseSeeker, len(coreSeekers))

	for i, coreSeeker := range coreSeekers {
		responseSeekers[i] = toResponseSeeker(coreSeeker)
	}

	return ResponseSeekers{
		ResponseSeekers: responseSeekers,
	}
}
