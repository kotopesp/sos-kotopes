package postservice

import (
	"context"	
	"github.com/kotopesp/sos-kotopes/internal/core"
)

func (s *postService) GetAnimalByID(ctx context.Context, id int) (core.Animal, error) {
	return s.AnimalStore.GetAnimalByID(ctx, id)
}
