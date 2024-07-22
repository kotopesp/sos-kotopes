package errors

import "errors"

var (
    ErrPostAlreadyInFavorites = errors.New("post already added to favorites")
)