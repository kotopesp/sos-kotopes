package errors

import "errors"

var (
	ErrInvalidChatID    = errors.New("chat with this id is not existing")
	ErrInvalidMessageID = errors.New("message with this id is not existing")
)
