package core

import (
	"errors"
	"gorm.io/gorm"
)

var (
	// post errors
    ErrInvalidPostID      		    = errors.New("invalid post ID")
    ErrPostNotFound       		    = errors.New("post not found")
	ErrRecordNotFound 				= gorm.ErrRecordNotFound

	// file errors
	ErrFailedToOpenImage   		    = errors.New("failed to open image")
	ErrFailedToReadImage  		    = errors.New("failed to read image")

	// user errors
	ErrFailedToGetAuthorIDFromToken = errors.New("failed to get author ID from token")

	// animal errors
	ErrAnimalNotFound 				= errors.New("animal not found")

	// favourite errors
	ErrPostAlreadyInFavourites 	    = errors.New("post already added to favourites")
)
