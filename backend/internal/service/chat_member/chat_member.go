package chatmember

import (
	"context"

	"github.com/kotopesp/sos-kotopes/internal/core"
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

func (s *service) GetAllMembers(ctx context.Context, chatId int) (members []core.ChatMember, total int, err error) {
	members, err = s.ChatMemberStore.GetAllMembers(ctx, chatId)
	if err != nil {
		return
	}
	total = len(members)
	return
}

func (s *service) AddMemberToChat(ctx context.Context, data core.ChatMember) (member core.ChatMember, err error) {
	return s.ChatMemberStore.AddMemberToChat(ctx, data)
}

func (s *service) UpdateMemberInfo(ctx context.Context, chatId int, userId int) (member core.ChatMember, err error) {
	return s.ChatMemberStore.UpdateMemberInfo(ctx, chatId, userId)
}

func (s *service) DeleteMemberFromChat(ctx context.Context, chatId int, userId int) (err error) {
	return s.ChatMemberStore.DeleteMemberFromChat(ctx, chatId, userId)
}
