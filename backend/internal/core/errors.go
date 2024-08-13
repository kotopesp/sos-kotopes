package core

import (
	"errors"
	"gorm.io/gorm"
)

var (
	// post errors
	ErrInvalidPostID        = errors.New("invalid post ID")
	ErrPostNotFound         = errors.New("post not found")
	ErrRecordNotFound       = gorm.ErrRecordNotFound
	ErrPostIsDeleted        = errors.New("post is deleted")
	ErrPostAuthorIDMismatch = errors.New("your user_id and db author_id mismatch")

	// user errors
	ErrFailedToGetAuthorIDFromToken = errors.New("failed to get author ID from token")

	// animal errors
	ErrAnimalNotFound = errors.New("animal not found")

	// favourite errors
	ErrPostAlreadyInFavourites       = errors.New("post already added to favourites")
	ErrUserAlreadyInFavourites       = errors.New("user already added to favourites")
	ErrCantAddYourselfIntoFavourites = errors.New("user can't add himself into his favourites list")
)
