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

func (s *service) CreateEquipment(ctx context.Context, equipment core.Equipment) (int, error) {
	return s.seekersStore.CreateEquipment(ctx, equipment)
}

func (s *service) CreateSeeker(ctx context.Context, seeker core.Seekers) (core.Seekers, error) {
	return s.seekersStore.CreateSeeker(ctx, seeker)
}

func (s *service) GetSeeker(ctx context.Context, id int) (core.Seekers, error) {
	return s.seekersStore.GetSeeker(ctx, id)
}

func (s *service) UpdateSeeker(ctx context.Context, seeker core.UpdateSeekers) (core.Seekers, error) {
	return s.seekersStore.UpdateSeeker(ctx, seeker)
}
