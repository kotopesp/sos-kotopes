package chatmember

import (
	"context"
	"time"

	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/internal/store/errors"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/postgres"
)

type (
	store struct {
		*postgres.Postgres
	}
)

func New(pg *postgres.Postgres) core.ChatMemberStore {
	return &store{pg}
}

func ifChatExists(s *store, ctx context.Context, chatID int) error {
	var counter int64
	if err := s.DB.WithContext(ctx).Model(&core.Chat{}).Where("id", chatID).Where("is_deleted", false).Count(&counter).Error; err != nil {
		return err
	}
	if counter != 1 {
		return errors.ErrInvalidChatID
	}
	return nil
}

func (s *store) GetAllMembers(ctx context.Context, chatID int) (members []core.ChatMember, err error) {
	if err := ifChatExists(s, ctx, chatID); err != nil {
		return nil, err
	}
	query := s.DB.WithContext(ctx).Model(&core.ChatMember{}).Where("chat_id", chatID).Where("is_deleted", false)
	if err := query.Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (s *store) AddMemberToChat(ctx context.Context, data core.ChatMember) (core.ChatMember, error) {
	if err := ifChatExists(s, ctx, data.ChatID); err != nil {
		return core.ChatMember{}, err
	}
	if err := s.DB.WithContext(ctx).Create(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}

func (s *store) UpdateMemberInfo(ctx context.Context, chatID, userID int) (core.ChatMember, error) {
	if err := ifChatExists(s, ctx, chatID); err != nil {
		return core.ChatMember{}, err
	}
	if err := s.DB.WithContext(ctx).Model(&core.ChatMember{}).Where("user_id", userID).Where("chat_id", chatID).Updates(map[string]interface{}{"updated_at": time.Now()}).Error; err != nil {
		return core.ChatMember{}, err
	}
	var member core.ChatMember
	s.DB.WithContext(ctx).Model(&core.Message{}).Where("id", userID).First(&member)
	return member, nil
}

func (s *store) DeleteMemberFromChat(ctx context.Context, chatID, userID int) (err error) {
	if err := ifChatExists(s, ctx, chatID); err != nil {
		return err
	}
	if err := s.DB.WithContext(ctx).Model(&core.ChatMember{}).Where("chat_id", chatID).Where("user_id", userID).Where("is_deleted", false).Updates(map[string]interface{}{"is_deleted": true, "deleted_at": time.Now()}).Error; err != nil {
		return err
	}
	return nil
}
