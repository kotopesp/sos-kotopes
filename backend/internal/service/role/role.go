package role

import (
	"context"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
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

func (s *service) GiveRoleToUser(ctx context.Context, id int, givenRole core.GivenRole) (addedRole core.RoleDetails, err error) {
	role, err := s.roleStore.GiveRoleToUser(ctx, id, givenRole)
	if err != nil {
		return core.RoleDetails{}, err
	}

	user, err := s.userStore.GetUserByID(ctx, id)
	username := user.Username
	addedRole = toRoleDetails(&role, givenRole.Name, username)
	return addedRole, err
}

func (s *service) GetUserRoles(ctx context.Context, id int) (rolesDetails []core.RoleDetails, err error) {
	roles, err := s.roleStore.GetUserRoles(ctx, id)
	if err != nil {
		return nil, err
	}

	user, err := s.userStore.GetUserByID(ctx, id)
	username := user.Username
	if seekerRole, exists := roles["seeker"]; exists {
		rolesDetails = append(rolesDetails, toRoleDetails(&seekerRole, "seeker", username))
	}
	if keeperRole, exists := roles["keeper"]; exists {
		rolesDetails = append(rolesDetails, toRoleDetails(&keeperRole, "keeper", username))
	}
	if vetRole, exists := roles["seeker"]; exists {
		rolesDetails = append(rolesDetails, toRoleDetails(&vetRole, "seeker", username))
	}
	return rolesDetails, err
}

func (s *service) UpdateUserRole(ctx context.Context, id int, roleForUpdate core.UpdateRole) (updatedRole core.RoleDetails, err error) {
	role, err := s.roleStore.UpdateUserRole(ctx, id, roleForUpdate)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.RoleDetails{}, err
	}

	user, err := s.userStore.GetUserByID(ctx, id)
	username := user.Username
	updatedRole = toRoleDetails(&role, roleForUpdate.Name, username)
	return updatedRole, err
}

func (s *service) DeleteUserRole(ctx context.Context, id int, roleForDelete string) (err error) {
	return s.roleStore.DeleteUserRole(ctx, id, roleForDelete)
}

func toRoleDetails(role *core.Role, roleName, username string) core.RoleDetails {
	if role == nil {
		return core.RoleDetails{}
	}
	return core.RoleDetails{
		ID:          role.ID,
		Name:        roleName,
		UserID:      role.UserID,
		Username:    username,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}
