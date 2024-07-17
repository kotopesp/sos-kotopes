package seeker

import "gitflic.ru/spbu-se/sos-kotopes/internal/core"

func (s *Seeker) ToCoreSeeker() *core.Seeker {
	if s == nil {
		return nil
	}
	return &core.Seeker{
		Description: s.Description,
		Location:    s.Location,
	}
}
