package seeker

import (
	"context"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

type service struct {
	seekersStore core.SeekersStore
}

func New(seekersStore core.SeekersStore) core.SeekersService {
	return &service{seekersStore: seekersStore}
}

func (s *service) CreateSeeker(ctx context.Context, seeker core.Seeker, equipment core.Equipment) (core.Seeker, error) {
	return s.seekersStore.CreateSeeker(ctx, seeker, equipment)
}

func (s *service) GetSeeker(ctx context.Context, id int) (core.Seeker, error) {
	return s.seekersStore.GetSeeker(ctx, id)
}

func (s *service) UpdateSeeker(ctx context.Context, seeker core.UpdateSeeker) (core.Seeker, error) {
	return s.seekersStore.UpdateSeeker(ctx, seeker)
}
