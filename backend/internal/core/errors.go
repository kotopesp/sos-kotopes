package core

import (
	"errors"
	"gorm.io/gorm"
)

var (
    ErrInvalidPostID      		    = errors.New("invalid post ID")
    ErrPostNotFound       		    = errors.New("post not found")
    ErrInternalServerError 		    = errors.New("internal server error")
	ErrInvalidInput        		    = errors.New("invalid input")
	ErrInvalidAuthorID       		= errors.New("invalid author ID")
	ErrInvalidAnimalID     		    = errors.New("invalid animal ID")
	ErrFailedToOpenImage   		    = errors.New("failed to open image")
	ErrFailedToReadImage  		    = errors.New("failed to read image")
	ErrPhotoNotFound       		    = errors.New("photo not found")
	ErrPhotoRequired 	   		    = errors.New("photo is required")
	ErrFailedToGetAuthorIDFromToken = errors.New("failed to get author ID from token")
	ErrPostAlreadyInFavorites 	    = errors.New("post already added to favorites")
	ErrRecordNotFound 				= gorm.ErrRecordNotFound
	ErrAnimalNotFound 				= errors.New("animal not found")
	ErrInvalidUserID 				= errors.New("invalid user ID")
)
