package vet_review

import (
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/user"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

func (p *GetAllVetReviewsParams) FromVetReviewRequest() core.GetAllVetReviewsParams {
	return core.GetAllVetReviewsParams{
		Limit:  &p.Limit,
		Offset: &p.Offset,
	}
}

func (v *VetReviewsCreate) ToCoreNewVetReview() core.VetReviews {
	if v == nil {
		return core.VetReviews{}
	}

	return core.VetReviews{
		AuthorID: v.AuthorID,
		Content:  v.Content,
		Grade:    v.Grade,
		VetID:    v.VetID,
	}
}

func (v *VetReviewsUpdate) ToCoreUpdateVetReview() core.UpdateVetReviews {
	return core.UpdateVetReviews{
		ID:       v.ID,
		AuthorID: v.AuthorID,
		Grade:    v.Grade,
		Content:  v.Content,
	}
}

func FromCoreVetReview(coreReview core.VetReviews) VetReviewsResponse {
	return VetReviewsResponse{
		ID:        coreReview.ID,
		AuthorID:  coreReview.AuthorID,
		Content:   coreReview.Content,
		Grade:     coreReview.Grade,
		VetID:     coreReview.VetID,
		CreatedAt: coreReview.CreatedAt,
		UpdatedAt: coreReview.UpdatedAt,
	}
}

func FromCoreVetReviewDetails(coreReview core.VetReviewsDetails) VetReviewsResponseWithUser {
	return VetReviewsResponseWithUser{
		Review: FromCoreVetReview(coreReview.Review),
		User:   user.ToResponseUser(&coreReview.User),
	}
}
