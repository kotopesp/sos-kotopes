package seeker

import "github.com/kotopesp/sos-kotopes/internal/core"

func ToResponseSeeker(seekers *core.Seekers) ResponseSeekers {
	if seekers == nil {
		return ResponseSeekers{}
	}
	return ResponseSeekers{
		ID:          seekers.ID,
		UserID:      seekers.UserID,
		Location:    seekers.Location,
		Equipment:   seekers.Equipment,
		Description: seekers.Description,
		Car:         seekers.Car,
	}
}

func (seekers *ResponseSeekers) ToCoreSeeker() core.Seekers {
	return core.Seekers{
		ID:          seekers.ID,
		UserID:      seekers.UserID,
		Location:    seekers.Location,
		Equipment:   seekers.Equipment,
		Description: seekers.Description,
		Car:         seekers.Car,
	}
}
