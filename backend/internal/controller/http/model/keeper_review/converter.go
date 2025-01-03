package keeperreview

import (
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/user"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

func (p *GetAllKeeperReviewsParams) ToCoreGetAllKeeperReviewsParams() core.GetAllKeeperReviewsParams {
	limit, offset := 10, 0
	if p.Limit != nil {
		limit = *p.Limit
	}
	if p.Offset != nil {
		offset = *p.Offset
	}

	return core.GetAllKeeperReviewsParams{
		Limit:  &limit,
		Offset: &offset,
	}
}

func (k *CreateKeeperReview) ToCoreKeeperReview() core.KeeperReview {
	if k == nil {
		return core.KeeperReview{}
	}

	return core.KeeperReview{
		Content: k.Content,
		Grade:   k.Grade,
	}
}

func (k *UpdateKeeperReview) ToCoreUpdateKeeperReview() core.KeeperReview {
	ukr := core.KeeperReview{}

	if k.Content != nil {
		ukr.Content = k.Content
	}
	if k.Grade != nil {
		ukr.Grade = *k.Grade
	}

	return ukr
}

func ToModelResponseKeeperReview(k core.KeeperReview) ResponseKeeperReview {
	return ResponseKeeperReview{
		ID:        k.ID,
		AuthorID:  k.AuthorID,
		Author:    user.ToResponseUser(&k.Author),
		Content:   k.Content,
		Grade:     k.Grade,
		KeeperID:  k.KeeperID,
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}
