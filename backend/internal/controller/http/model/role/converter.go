package role

import "github.com/kotopesp/sos-kotopes/internal/core"

func ToRole(roleDetails *core.RoleDetails) Role {
	if roleDetails == nil {
		return Role{}
	}
	return Role{
		Name:        roleDetails.Name,
		ID:          roleDetails.ID,
		Username:    roleDetails.Username,
		Description: roleDetails.Description,
		CreatedAt:   roleDetails.CreatedAt,
		UpdatedAt:   roleDetails.UpdatedAt,
	}
}

func (r *GivenRole) ToCoreGivenRole() core.GivenRole {
	if r == nil {
		return core.GivenRole{}
	}
	return core.GivenRole{
		Name:        r.Name,
		Description: r.Description,
	}
}

func (r *UpdateRole) ToCoreUpdateRole() core.UpdateRole {
	if r == nil {
		return core.UpdateRole{}
	}
	return core.UpdateRole{
		Name:        r.Name,
		Description: r.Description,
	}
}
