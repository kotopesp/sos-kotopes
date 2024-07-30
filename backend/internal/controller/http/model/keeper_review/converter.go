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
