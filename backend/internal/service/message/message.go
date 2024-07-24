package message

import (
	"context"

	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
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

func (s *service) GetAllMessages(ctx context.Context, chatId int, sortType string, searchText string) (messages []core.Message, total int, err error) {
	messages, err = s.MessageStore.GetAllMessages(ctx, chatId, sortType, searchText)
	if err != nil {
		return
	}
	total = len(messages)
	return
}

func (s *service) CreateMessage(ctx context.Context, data core.Message) (core.Message, error) {
	return s.MessageStore.CreateMessage(ctx, data)
}

func (s *service) UpdateMessage(ctx context.Context, chatId int, messageId int, data string) (core.Message, error) {
	return s.MessageStore.UpdateMessage(ctx, chatId, messageId, data)
}

func (s *service) DeleteMessage(ctx context.Context, chatId int, messageId int) (err error) {
	return s.MessageStore.DeleteMessage(ctx, chatId, messageId)
}
