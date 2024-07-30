package core

import "errors"

var (
	ErrReviewGradeBounds = errors.New("Grade must be between 1 and 5")
)
