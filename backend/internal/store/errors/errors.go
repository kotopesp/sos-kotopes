package errors

import "errors"

var (
	ErrInvalidChatId    = errors.New("chat with this id is not existing")
	ErrInvalidMessageId = errors.New("message with this id is not existing")
)
