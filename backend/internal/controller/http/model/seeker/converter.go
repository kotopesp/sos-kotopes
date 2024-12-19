package seeker

import (
	"github.com/kotopesp/sos-kotopes/internal/core"
	"strings"
)

func (seeker *CreateSeeker) GetEquipment() (core.Equipment, error) {
	if len(seeker.Equipment) != 5 {
		return core.Equipment{}, core.ErrInvalidEquipment
	}

	return core.Equipment{
		HaveMetalCage:   seeker.Equipment[0] != "",
		HavePlasticCage: seeker.Equipment[1] != "",
		HaveNet:         seeker.Equipment[2] != "",
		HaveLadder:      seeker.Equipment[3] != "",
		HaveOther:       seeker.Equipment[4],
	}, nil
}

func (seeker *CreateSeeker) ToCoreSeeker() core.Seeker {
	return core.Seeker{
		UserID:           seeker.UserID,
		AnimalType:       seeker.AnimalType,
		Description:      seeker.Description,
		Location:         seeker.Location,
		EquipmentRental:  seeker.EquipmentRental,
		HaveCar:          seeker.HaveCar,
		Price:            seeker.Price,
		WillingnessCarry: seeker.WillingnessCarry,
	}
}

func (seeker *UpdateSeeker) ToCoreUpdateSeeker() core.UpdateSeeker {
	return core.UpdateSeeker{
		UserID:           seeker.UserID,
		AnimalType:       seeker.AnimalType,
		Description:      seeker.Description,
		Location:         seeker.Location,
		EquipmentRental:  seeker.EquipmentRental,
		HaveCar:          seeker.HaveCar,
		Price:            seeker.Price,
		WillingnessCarry: seeker.WillingnessCarry,
	}
}

func parseSort(sort string) (sortBy, sortOrder string) {
	parts := strings.Split(sort, ":")
	if len(parts) != 2 {
		return "", ""
	}
	sortBy = ""
	sortOrder = ""
	if parts[0] != "" {
		sortBy = parts[0]
	}
	if parts[1] != "" {
		sortOrder = parts[1]
	}
	return sortBy, sortOrder
}

func (p *GetAllSeekerParams) ToCoreGetAllSeekersParams() core.GetAllSeekersParams {
	sortBy, sortOrder := "", ""

	if p.Sort != nil {
		sortBy, sortOrder = parseSort(*p.Sort)
	}

	limit, offset := 10, 0

	if p.Limit != nil {
		limit = *p.Limit
	}

	if p.Offset != nil {
		offset = *p.Offset
	}

	return core.GetAllSeekersParams{
		SortBy:     &sortBy,
		SortOrder:  &sortOrder,
		AnimalType: p.AnimalType,
		Location:   p.Location,
		Price:      p.Price,
		HaveCar:    p.HaveCar,
		Limit:      &limit,
		Offset:     &offset,
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
		Equipment:        seeker.EquipmentID,
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
		Equipment:        seeker.EquipmentID,
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
