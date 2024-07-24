package role_service

import (
	"context"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
)

type Service struct {
	roleStore core.RoleStore
}

func NewRoleService(store core.RoleStore) core.RoleService {
	return &Service{
		roleStore: store,
	}
}

func (s *Service) GetUserRoles(ctx context.Context, id int) (roles []core.Role, err error) {
	return s.roleStore.GetUserRoles(ctx, id)
}
func (s *Service) GiveRoleToUser(ctx context.Context, id int, role string) (err error) {
	return nil
}
func (s *Service) DeleteUserRole(ctx context.Context, id int, role string) (err error) {
	return nil
}
func (s *Service) UpdateUserRole(ctx context.Context, id int, role string) (err error) {
	return nil
}
