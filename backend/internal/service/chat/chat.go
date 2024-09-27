package chat

import (
	"context"

	"github.com/kotopesp/sos-kotopes/internal/core"
)

type (
	service struct {
		ChatStore core.ChatStore
		UserStore core.UserStore
	}
)

func New(chatStore core.ChatStore, userStore core.UserStore) core.ChatService {
	return &service{
		ChatStore: chatStore,
		UserStore: userStore,
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

func (s *service) GetChatWithUsersByID(ctx context.Context, id int) (chat core.Chat, err error) {
	chat, err = s.ChatStore.GetChatWithUsersByID(ctx, id)
	return
}

func (s *service) CreateChat(ctx context.Context, data core.Chat, userIds []int) (chat core.Chat, err error) {
	chat, err = s.ChatStore.CreateChat(ctx, data)
	if err != nil {
		return core.Chat{}, err
	}

	for _, userID := range userIds {
		chatMember := core.ChatMember{
			ChatID: chat.ID,
			UserID: userID,
		}
		if _, err := s.ChatStore.AddMemberToChat(ctx, chatMember); err != nil {
			return core.Chat{}, err
		}
		u, err := s.UserStore.GetUser(ctx, userID)
		if err != nil {
			return core.Chat{}, err
		}
		chat.Users = append(chat.Users, u)
	}
	return chat, nil
}

func (s *service) FindChatByUsers(ctx context.Context, userIds []int) (core.Chat, error) {
	chat, err := s.ChatStore.FindChatByUsers(ctx, userIds)
	if err != nil {
		return core.Chat{ID: -1}, err
	}
	return chat, nil
}

func (s *service) DeleteChat(ctx context.Context, id int) (err error) {
	err = s.ChatStore.DeleteChat(ctx, id)
	return
}

func (s *service) AddMemberToChat(ctx context.Context, data core.ChatMember) (member core.ChatMember, err error) {
	return s.ChatStore.AddMemberToChat(ctx, data)
}
