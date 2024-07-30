package role

import "github.com/kotopesp/sos-kotopes/internal/core"

func ToRole(roleDetails *core.RoleDetails) Role {
	if roleDetails == nil {
		return Role{}
	}
	return Role{
		Name:        roleDetails.Name,
		ID:          roleDetails.ID,
		UserID:      roleDetails.UserID,
		Description: roleDetails.Description,
		CreatedAt:   roleDetails.CreatedAt,
		UpdatedAt:   roleDetails.UpdatedAt,
	}
}
