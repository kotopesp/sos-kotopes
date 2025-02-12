package moderator

import (
	"context"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

type service struct {
	moderatorStore core.ModeratorStore
}

func New(store core.ModeratorStore) core.ModeratorService {
	return &service{}
}

func (s *service) GetModerator(ctx context.Context, id int) (core.Moderator, error) {
	return core.Moderator{}, nil
}

func (s *service) UpdateModerator(ctx context.Context, id int, update core.UpdateModerator) (core.Moderator, error) {
	return core.Moderator{}, nil
}
