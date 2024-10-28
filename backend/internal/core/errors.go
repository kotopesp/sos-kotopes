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

	// photo errors
	ErrPhotoNotFound = errors.New("photo not found")

	// user errors
	ErrFailedToGetAuthorIDFromToken = errors.New("failed to get author ID from token")

	// animal errors
	ErrAnimalNotFound = errors.New("animal not found")

	// favourite errors
	ErrPostAlreadyInFavourites = errors.New("post already added to favourites")

	// auth errors
	ErrUnauthorized                      = errors.New("unauthorized")
	ErrInvalidUpdateRefreshSessionParams = errors.New("invalid UpdateRefreshSession params")

	// comment errors
	ErrCommentAuthorIDMismatch     = errors.New("your user_id and db author_id mismatch")
	ErrCommentPostIDMismatch       = errors.New("your posts_id and db posts_id mismatch")
	ErrNoSuchComment               = errors.New("no such comment")
	ErrCommentIsDeleted            = errors.New("comment is deleted")
	ErrInvalidCommentParentID      = errors.New("invalid comment parent_id")
	ErrReplyToCommentOfAnotherPost = errors.New("reply to comment of another post")
	ErrParentCommentNotFound       = errors.New("parent comment not found")
	ErrReplyCommentNotFound        = errors.New("reply comment not found")
	ErrInvalidCommentReplyID       = errors.New("invalid comment reply_id")
	ErrNullCommentParentID         = errors.New("comment parent id is null")
)
