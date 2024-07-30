package role

import (
	"context"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

type service struct {
	roleStore core.RoleStore
	userStore core.UserStore
}

func New(roleStore core.RoleStore, userStore core.UserStore) core.RoleService {
	return &service{
		roleStore: roleStore,
		userStore: userStore,
	}
}

func (s *service) GetUserRoles(ctx context.Context, id int) (rolesDetails []core.RoleDetails, err error) {
	roles, err := s.roleStore.GetUserRoles(ctx, id)
	if err != nil {
		return nil, err
	}
	user, err := s.userStore.GetUserByID(ctx, id)
	username := user.Username
	if seekerRole, exists := roles["seeker"]; exists {
		rolesDetails = append(rolesDetails, seekerRole.ToRoleDetails("seeker", username))
	}
	if keeperRole, exists := roles["keeper"]; exists {
		rolesDetails = append(rolesDetails, keeperRole.ToRoleDetails("keeper", username))
	}
	if vetRole, exists := roles["seeker"]; exists {
		rolesDetails = append(rolesDetails, vetRole.ToRoleDetails("seeker", username))
	}
	return rolesDetails, err
}

func (s *service) GiveRoleToUser(ctx context.Context, id int, givenRole core.GivenRole) (addedRole core.RoleDetails, err error) {
	role, roleName, err := s.roleStore.GiveRoleToUser(ctx, id, givenRole)
	if err != nil {
		return core.RoleDetails{}, err
	}
	user, err := s.userStore.GetUserByID(ctx, id)
	username := user.Username
	addedRole = role.ToRoleDetails(roleName, username)
	return addedRole, err
}
func (s *service) DeleteUserRole(ctx context.Context, id int, roleForDelete string) (err error) {
	return s.roleStore.DeleteUserRole(ctx, id, roleForDelete)
}
func (s *service) UpdateUserRole(ctx context.Context, id int, roleForUpdate core.UpdateRole) (err error) {
	return s.roleStore.UpdateUserRole(ctx, id, roleForUpdate)
}
