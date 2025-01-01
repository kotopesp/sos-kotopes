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
		AuthorID: k.AuthorID,
		Content:  k.Content,
		Grade:    k.Grade,
		KeeperID: k.KeeperID,
	}
}

func (k *UpdateKeeperReview) ToCoreUpdateKeeperReview() core.KeeperReview {
	return core.KeeperReview{
		Grade:   *k.Grade,
		Content: k.Content,
	}
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
