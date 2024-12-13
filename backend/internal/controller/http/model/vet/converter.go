package vet

import (
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/user"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

func (p *GetAllVetsParams) FromVetRequestParams() core.GetAllVetParams {
	return core.GetAllVetParams{
		SortBy:    p.SortBy,
		SortOrder: p.SortOrder,
		Location:  p.Location,
		MinRating: p.MinRating,
		MaxRating: p.MaxRating,
		MinPrice:  p.MinPrice,
		MaxPrice:  p.MaxPrice,
		Limit:     p.Limit,
		Offset:    p.Offset,
	}
}

func (v *VetsCreate) ToCoreNewVet() core.Vets {
	return core.Vets{
		UserID:             v.UserID,
		IsOrganization:     v.IsOrganization,
		Username:           v.Username,
		Firstname:          v.Firstname,
		Lastname:           v.Lastname,
		Patronymic:         v.Patronymic,
		Education:          v.Education,
		OrgName:            v.OrgName,
		Location:           v.Location,
		OrgEmail:           v.OrgEmail,
		InnNumber:          v.InnNumber,
		IsRemoteConsulting: v.IsRemoteConsulting,
		IsInpatient:        v.IsInpatient,
		Description:        v.Description,
	}
}

func (v *VetsUpdate) ToCoreUpdateVet() core.UpdateVets {
	return core.UpdateVets{
		ID:                 v.ID,
		UserID:             v.UserID,
		IsOrganization:     v.IsOrganization,
		Username:           v.Username,
		Firstname:          v.Firstname,
		Lastname:           v.Lastname,
		Patronymic:         v.Patronymic,
		Education:          v.Education,
		OrgName:            v.OrgName,
		Location:           v.Location,
		OrgEmail:           v.OrgEmail,
		InnNumber:          v.InnNumber,
		IsRemoteConsulting: v.IsRemoteConsulting,
		IsInpatient:        v.IsInpatient,
		Description:        v.Description,
	}
}

func FromCoreVet(coreVet core.Vets) VetsResponse {
	return VetsResponse{
		ID:                 coreVet.ID,
		UserID:             coreVet.UserID,
		IsOrganization:     coreVet.IsOrganization,
		Username:           coreVet.Username,
		Firstname:          coreVet.Firstname,
		Lastname:           coreVet.Lastname,
		Patronymic:         coreVet.Patronymic,
		Education:          coreVet.Education,
		OrgName:            coreVet.OrgName,
		Location:           coreVet.Location,
		OrgEmail:           coreVet.OrgEmail,
		InnNumber:          coreVet.InnNumber,
		IsRemoteConsulting: coreVet.IsRemoteConsulting,
		IsInpatient:        coreVet.IsInpatient,
		Description:        coreVet.Description,
		CreatedAt:          coreVet.CreatedAt,
		UpdatedAt:          coreVet.UpdatedAt,
	}
}

func FromCoreVetWithUser(coreVet core.VetsDetails) VetsResponseWithUser {
	return VetsResponseWithUser{
		Vet:  FromCoreVet(coreVet.Vet),
		User: user.ToResponseUser(&coreVet.User),
	}
}
