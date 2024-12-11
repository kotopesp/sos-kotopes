package seeker

import "github.com/kotopesp/sos-kotopes/internal/core"

func checker(s string) bool {
	if len(s) == 0 {
		return false
	}
	return true
}

func (seeker *CreateSeeker) GetEquipment() core.Equipment {
	return core.Equipment{
		HaveMetalCage:   checker(seeker.Equipment[0]),
		HavePlasticCage: checker(seeker.Equipment[1]),
		HaveNet:         checker(seeker.Equipment[2]),
		HaveLadder:      checker(seeker.Equipment[3]),
		HaveOther:       seeker.Equipment[4],
	}
}

func (seeker *CreateSeeker) ToCoreSeeker(equipmentId int) core.Seekers {
	return core.Seekers{
		UserID:      seeker.UserID,
		IdEquipment: equipmentId,
		Description: seeker.Description,
		Location:    seeker.Location,
		HaveCar:     seeker.HaveCar,
	}
}

func ToResponseSeeker(seeker *core.Seekers) ResponseSeeker {
	if seeker == nil {
		return ResponseSeeker{}
	}
	return ResponseSeeker{
		ID:          seeker.ID,
		UserID:      seeker.UserID,
		Location:    seeker.Location,
		Equipment:   seeker.IdEquipment,
		Description: seeker.Description,
		HaveCar:     seeker.HaveCar,
	}
}
