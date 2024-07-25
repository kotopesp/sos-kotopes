package message

import (
	"context"

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

func (s *service) GetAllMessages(ctx context.Context, chatID int, sortType, searchText string) (messages []core.Message, total int, err error) {
	messages, err = s.MessageStore.GetAllMessages(ctx, chatID, sortType, searchText)
	if err != nil {
		return
	}
	total = len(messages)
	return
}

func (s *service) CreateMessage(ctx context.Context, data core.Message) (core.Message, error) {
	return s.MessageStore.CreateMessage(ctx, data)
}

func (s *service) UpdateMessage(ctx context.Context, chatID, messageID int, data string) (core.Message, error) {
	return s.MessageStore.UpdateMessage(ctx, chatID, messageID, data)
}

func (s *service) DeleteMessage(ctx context.Context, chatID, messageID int) (err error) {
	return s.MessageStore.DeleteMessage(ctx, chatID, messageID)
}
