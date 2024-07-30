package role

import (
	"context"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/role"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

type service struct {
	roleStore core.RoleStore
}

func New(store core.RoleStore) core.RoleService {
	return &service{
		roleStore: store,
	}
}

func (s *service) GetUserRoles(ctx context.Context, id int) (roles []core.Role, err error) {
	return s.roleStore.GetUserRoles(ctx, id)
}
func (s *service) GiveRoleToUser(ctx context.Context, id int, role role.GiveRole) (err error) {
	return s.roleStore.GiveRoleToUser(ctx, id, role)
}
func (s *service) DeleteUserRole(ctx context.Context, id int, role string) (err error) {
	return s.roleStore.DeleteUserRole(ctx, id, role)
}
func (s *service) UpdateUserRole(ctx context.Context, id int, role role.UpdateRole) (err error) {
	return s.roleStore.UpdateUserRole(ctx, id, role)
}
