package seeker

import "github.com/kotopesp/sos-kotopes/internal/core"

func (seekers *CreateSeeker) ToCoreSeeker() core.Seekers {
	return core.Seekers{
		UserID:      seekers.UserID,
		Location:    seekers.Location,
		Equipment:   seekers.Equipment,
		Description: seekers.Description,
		HaveCar:     seekers.HaveCar,
	}
}

func ToResponseSeeker(seekers *core.Seekers) ResponseSeeker {
	if seekers == nil {
		return ResponseSeeker{}
	}
	return ResponseSeeker{
		ID:          seekers.ID,
		UserID:      seekers.UserID,
		Location:    seekers.Location,
		Equipment:   seekers.Equipment,
		Description: seekers.Description,
		HaveCar:     seekers.HaveCar,
	}
}
