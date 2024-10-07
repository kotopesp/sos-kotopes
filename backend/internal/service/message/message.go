package message

import (
	"context"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/chat"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

type (
	service struct {
		MessageStore core.MessageStore
	}
)

func New(store core.MessageStore) core.MessageService {
	return &service{
		MessageStore: store,
	}
}

func (s *service) GetAllMessages(ctx context.Context, chatID, userID int, sortType, searchText string) (messages []chat.Message, total int, err error) {
	messages, err = s.MessageStore.GetAllMessages(ctx, chatID, userID, sortType, searchText)
	if err != nil {
		return
	}
	total = len(messages)
	return
}

func (s *service) MarkMessagesAsRead(ctx context.Context, chatID, userID int) error {
	return s.MessageStore.MarkMessagesAsRead(ctx, chatID, userID)
}

func (s *service) GetUnreadMessageCount(ctx context.Context, chatID, userID int) (int64, error) {
	return s.MessageStore.GetUnreadMessageCount(ctx, chatID, userID)
}

func (s *service) CreateMessage(ctx context.Context, data chat.Message) (chat.Message, error) {
	return s.MessageStore.CreateMessage(ctx, data)
}

func (s *service) UpdateMessage(ctx context.Context, chatID, userID, messageID int, data string) (core.Message, error) {
	return s.MessageStore.UpdateMessage(ctx, chatID, userID, messageID, data)
}

func (s *service) DeleteMessage(ctx context.Context, chatID, userID, messageID int) (err error) {
	return s.MessageStore.DeleteMessage(ctx, chatID, userID, messageID)
}
