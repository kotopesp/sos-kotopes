package vet

import (
	"context"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

func (s *service) GetAllReviews(ctx context.Context, params core.GetAllVetReviewsParams) ([]core.VetReviewsDetails, error) {
	vetReviews, err := s.vetReviewsStore.GetAllReviews(ctx, params)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, err
	}

	vetReviewsDetails := make([]core.VetReviewsDetails, len(vetReviews))

	for i, review := range vetReviews {
		veterinaryReviewUser, err := s.userStore.GetUserByID(ctx, review.AuthorID)
		if err != nil {
			logger.Log().Error(ctx, err.Error())
			return nil, err
		}

		vetReviewsDetails[i] = core.VetReviewsDetails{
			Review: review,
			User:   veterinaryReviewUser,
		}
	}

	return vetReviewsDetails, nil
}

func (s *service) CreateReview(ctx context.Context, review core.VetReviews) error {
	if review.CreatedAt.IsZero() {
		review.CreatedAt = time.Now()
	}

	if review.Grade < 1 || review.Grade > 5 {
		return core.ErrReviewGradeBounds
	}

	return s.vetReviewsStore.CreateReview(ctx, review)
}

func (s *service) SoftDeleteReviewByID(ctx context.Context, id, userID int) error {
	storedReview, err := s.vetReviewsStore.GetByIDReview(ctx, id)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}

	if storedReview.AuthorID != userID {
		logger.Log().Error(ctx, core.ErrVetReviewUserIDMismatch.Error())
		return core.ErrVetReviewUserIDMismatch
	} else if storedReview.IsDeleted {
		logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
		return core.ErrRecordNotFound
	}

	return s.vetReviewsStore.SoftDeleteReviewByID(ctx, id)
}

func (s *service) UpdateReviewByID(ctx context.Context, review core.UpdateVetReviews) (core.VetReviewsDetails, error) {
	storedReview, err := s.vetReviewsStore.GetByIDReview(ctx, review.ID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.VetReviewsDetails{}, err
	}

	if storedReview.AuthorID != review.AuthorID {
		logger.Log().Error(ctx, core.ErrVetReviewUserIDMismatch.Error())
		return core.VetReviewsDetails{}, core.ErrVetReviewUserIDMismatch
	} else if storedReview.IsDeleted {
		logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
		return core.VetReviewsDetails{}, core.ErrRecordNotFound
	}

	updatedVetReview, err := s.vetReviewsStore.UpdateReviewByID(ctx, review)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.VetReviewsDetails{}, err
	}

	veterinaryReviewUser, err := s.userStore.GetUserByID(ctx, updatedVetReview.AuthorID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.VetReviewsDetails{}, err
	}

	return core.VetReviewsDetails{
		Review: updatedVetReview,
		User:   veterinaryReviewUser,
	}, nil
}
