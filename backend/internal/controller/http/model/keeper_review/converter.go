package keeperreview

import "github.com/kotopesp/sos-kotopes/internal/core"

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

func FromCoreKeeperReview(coreReview core.KeeperReviews) KeeperReviewsResponse {
	return KeeperReviewsResponse{
		ID:        coreReview.ID,
		AuthorID:  coreReview.AuthorID,
		Content:   coreReview.Content,
		Grade:     coreReview.Grade,
		KeeperID:  coreReview.KeeperID,
		DeletedAt: coreReview.DeletedAt,
		CreatedAt: coreReview.CreatedAt,
		UpdatedAt: coreReview.UpdatedAt,
	}
}
