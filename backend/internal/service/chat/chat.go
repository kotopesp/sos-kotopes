package chat

import (
	"context"

	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
)

type (
	service struct {
		ChatStore core.ChatStore
	}
)

func New(store core.ChatStore) core.ChatService {
	return &service{
		ChatStore: store,
	}
}

func (s *service) GetAll(ctx context.Context, sortType string) (chats []core.Chat, total int, err error) {
	chats, err = s.ChatStore.GetAll(ctx, sortType)
	if err != nil {
		return
	}
	total = len(chats)
	return
}

func (s *service) GetByID(ctx context.Context, id int) (chat core.Chat, err error) {
	chat, err = s.ChatStore.GetByID(ctx, id)
	return
}

func (s *service) Create(ctx context.Context, data core.Chat) (chat core.Chat, err error) {
	chat, err = s.ChatStore.Create(ctx, data)
	return
}

func (s *service) Delete(ctx context.Context, id int) (err error) {
	err = s.ChatStore.Delete(ctx, id)
	return
}
