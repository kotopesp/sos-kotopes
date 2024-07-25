package chat

import (
	"context"

	"github.com/kotopesp/sos-kotopes/internal/core"
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

func (s *service) GetAllChats(ctx context.Context, sortType string) (chats []core.Chat, total int, err error) {
	chats, err = s.ChatStore.GetAllChats(ctx, sortType)
	if err != nil {
		return
	}
	total = len(chats)
	return
}

func (s *service) GetChatByID(ctx context.Context, id int) (chat core.Chat, err error) {
	chat, err = s.ChatStore.GetChatByID(ctx, id)
	return
}

func (s *service) CreateChat(ctx context.Context, data core.Chat) (chat core.Chat, err error) {
	chat, err = s.ChatStore.CreateChat(ctx, data)
	return
}

func (s *service) DeleteChat(ctx context.Context, id int) (err error) {
	err = s.ChatStore.DeleteChat(ctx, id)
	return
}
