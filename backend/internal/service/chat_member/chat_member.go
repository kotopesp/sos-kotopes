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

func (s *service) GetAllMembers(ctx context.Context, chatID, userID int) (members []core.ChatMember, total int, err error) {
	members, err = s.ChatMemberStore.GetAllMembers(ctx, chatID, userID)
	if err != nil {
		return
	}
	total = len(members)
	return
}

func (s *service) AddMemberToChat(ctx context.Context, data core.ChatMember, userID int) (member core.ChatMember, err error) {
	return s.ChatMemberStore.AddMemberToChat(ctx, data, userID)
}

func (s *service) UpdateMemberInfo(ctx context.Context, chatID, userID, currentUserID int) (member core.ChatMember, err error) {
	return s.ChatMemberStore.UpdateMemberInfo(ctx, chatID, userID, currentUserID)
}

func (s *service) DeleteMemberFromChat(ctx context.Context, chatID, userID, currentUserID int) (err error) {
	return s.ChatMemberStore.DeleteMemberFromChat(ctx, chatID, userID, currentUserID)
}
