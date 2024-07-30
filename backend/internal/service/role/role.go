package role

import (
	"context"
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

func (s *service) GetUserRoles(ctx context.Context, id int) (rolesDetails []core.RoleDetails, err error) {
	roles, err := s.roleStore.GetUserRoles(ctx, id)
	if seekerRole, exists := roles["seeker"]; exists {
		rolesDetails = append(rolesDetails, seekerRole.ToRoleDetails("seeker"))
	}
	if keeperRole, exists := roles["keeper"]; exists {
		rolesDetails = append(rolesDetails, keeperRole.ToRoleDetails("keeper"))
	}
	if vetRole, exists := roles["seeker"]; exists {
		rolesDetails = append(rolesDetails, vetRole.ToRoleDetails("seeker"))
	}
	return rolesDetails, err
}
func (s *service) GiveRoleToUser(ctx context.Context, id int, givenRole core.GiveRole) (err error) {
	return s.roleStore.GiveRoleToUser(ctx, id, givenRole)
}
func (s *service) DeleteUserRole(ctx context.Context, id int, roleForDelete string) (err error) {
	return s.roleStore.DeleteUserRole(ctx, id, roleForDelete)
}
func (s *service) UpdateUserRole(ctx context.Context, id int, roleForUpdate core.UpdateRole) (err error) {
	return s.roleStore.UpdateUserRole(ctx, id, roleForUpdate)
}
