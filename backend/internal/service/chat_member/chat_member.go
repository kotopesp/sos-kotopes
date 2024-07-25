package chatmember

import (
	"context"

	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
)

type (
	service struct {
		ChatMemberStore core.ChatMemberStore
	}
)

func New(store core.ChatMemberStore) core.ChatMemberService {
	return &service{
		ChatMemberStore: store,
	}
}

func (s *service) GetAllMembers(ctx context.Context, chatID int) (members []core.ChatMember, total int, err error) {
	members, err = s.ChatMemberStore.GetAllMembers(ctx, chatID)
	if err != nil {
		return
	}
	total = len(members)
	return
}

func (s *service) AddMemberToChat(ctx context.Context, data core.ChatMember) (member core.ChatMember, err error) {
	return s.ChatMemberStore.AddMemberToChat(ctx, data)
}

func (s *service) UpdateMemberInfo(ctx context.Context, chatID, userID int) (member core.ChatMember, err error) {
	return s.ChatMemberStore.UpdateMemberInfo(ctx, chatID, userID)
}

func (s *service) DeleteMemberFromChat(ctx context.Context, chatID, userID int) (err error) {
	return s.ChatMemberStore.DeleteMemberFromChat(ctx, chatID, userID)
}
