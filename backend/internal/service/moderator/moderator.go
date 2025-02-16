package moderator

import (
	"context"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

type service struct {
	moderatorStore core.ModeratorStore
}

func New(store core.ModeratorStore) core.ModeratorService {
	return &service{}
}

// GetModerator - returns moderator struct by its id.
func (s *service) GetModerator(ctx context.Context, id int) (moderator core.Moderator, err error) {
	moderator, err = s.moderatorStore.GetModeratorByID(ctx, id)
	if err != nil {
		logger.Log().Error(ctx, err.Error())

		return core.Moderator{}, err
	}

	return moderator, nil
}
