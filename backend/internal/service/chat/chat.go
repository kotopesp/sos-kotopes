package chat

import (
	"context"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/chat"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

type (
	service struct {
		ChatStore core.ChatStore
		UserStore core.UserStore
	}
)

func convertToModelUser(user core.User) (userResponse chat.User) {
	userResponse.ID = user.ID
	userResponse.Username = user.Username
	return
}

func New(chatStore core.ChatStore, userStore core.UserStore) core.ChatService {
	return &service{
		ChatStore: chatStore,
		UserStore: userStore,
	}
}

func (s *service) GetAllChats(ctx context.Context, sortType string, userID int) (chats []chat.Chat, total int, err error) {
	chats, err = s.ChatStore.GetAllChats(ctx, sortType, userID)
	if err != nil {
		return
	}
	total = len(chats)
	return
}

func (s *service) GetChatWithUsersByID(ctx context.Context, id int) (data chat.Chat, err error) {
	data, err = s.ChatStore.GetChatWithUsersByID(ctx, id)
	return
}

func (s *service) CreateChat(ctx context.Context, data chat.Chat, userIds []int) (chat.Chat, error) {
	createdChat, err := s.ChatStore.CreateChat(ctx, data)
	if err != nil {
		return chat.Chat{}, err
	}

	for _, userID := range userIds {
		chatMember := core.ChatMember{
			ChatID: createdChat.ID,
			UserID: userID,
		}
		if _, err := s.ChatStore.AddMemberToChat(ctx, chatMember); err != nil {
			return chat.Chat{}, err
		}
		u, err := s.UserStore.GetUser(ctx, userID)
		if err != nil {
			return chat.Chat{}, err
		}
		createdChat.Users = append(createdChat.Users, convertToModelUser(u))
	}
	return createdChat, nil
}

func (s *service) FindChatByUsers(ctx context.Context, userIds []int) (chat.Chat, error) {
	foundChat, err := s.ChatStore.FindChatByUsers(ctx, userIds)
	if err != nil {
		return chat.Chat{ID: -1}, err
	}
	return foundChat, nil
}

func (s *service) DeleteChat(ctx context.Context, id int) (err error) {
	err = s.ChatStore.DeleteChat(ctx, id)
	return
}

func (s *service) AddMemberToChat(ctx context.Context, data core.ChatMember) (member core.ChatMember, err error) {
	return s.ChatStore.AddMemberToChat(ctx, data)
}
