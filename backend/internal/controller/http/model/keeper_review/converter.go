package keeperreview

import (
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/user"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

func (p *GetAllKeeperReviewsParams) FromKeeperReviewRequest() core.GetAllKeeperReviewsParams {
	return core.GetAllKeeperReviewsParams{
		Limit:  &p.Limit,
		Offset: &p.Offset,
	}
}

func (k *KeeperReviewsCreate) ToCoreNewKeeperReview() core.KeeperReviews {
	if k == nil {
		return core.KeeperReviews{}
	}

	return core.KeeperReviews{
		AuthorID: k.AuthorID,
		Content:  k.Content,
		Grade:    k.Grade,
		KeeperID: k.KeeperID,
	}
}

func (k *KeeperReviewsUpdate) ToCoreUpdateKeeperReview() core.UpdateKeeperReviews {
	return core.UpdateKeeperReviews{
		ID:       k.ID,
		AuthorID: k.AuthorID,
		Grade:    k.Grade,
		Content:  k.Content,
	}
}

func FromCoreKeeperReview(coreReview core.KeeperReviews) KeeperReviewsResponse {
	return KeeperReviewsResponse{
		ID:        coreReview.ID,
		AuthorID:  coreReview.AuthorID,
		Content:   coreReview.Content,
		Grade:     coreReview.Grade,
		KeeperID:  coreReview.KeeperID,
		CreatedAt: coreReview.CreatedAt,
		UpdatedAt: coreReview.UpdatedAt,
	}
}

func FromCoreKeeperReviewDetails(coreReview core.KeeperReviewsDetails) KeeperReviewsResponseWithUser {
	return KeeperReviewsResponseWithUser{
		Review: FromCoreKeeperReview(coreReview.Review),
		User:   user.ToResponseUser(&coreReview.User),
	}
}
