package core

import "errors"

var (
	ErrReviewGradeBounds = errors.New("grade must be between 1 and 5")
	ErrRecordNotFound    = errors.New("record not found")
)
