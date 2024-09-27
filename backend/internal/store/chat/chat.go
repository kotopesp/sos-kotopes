package chat

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/kotopesp/sos-kotopes/internal/core"
	err "github.com/kotopesp/sos-kotopes/internal/store/errors"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
)

type (
	store struct {
		*postgres.Postgres
	}
)

func New(pg *postgres.Postgres) core.ChatStore {
	return &store{pg}
}

func ifChatExists(s *store, ctx context.Context, chatID int) error {
	var counter int64
	if err := s.DB.WithContext(ctx).Model(&core.Chat{}).Where("id", chatID).Where("is_deleted", false).Count(&counter).Error; err != nil {
		return err
	}
	if counter != 1 {
		return err.ErrInvalidChatID
	}
	return nil
}

func (s *store) GetAllChats(ctx context.Context, sortType string) ([]core.Chat, error) {
	var chats []core.Chat

	query := s.DB.WithContext(ctx).Model(&core.Chat{}).Where("is_deleted = false")

	if sortType != "" {
		query = query.Where("chat_type", sortType)
	}
	if err := query.Preload("Users").Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}

func (s *store) GetChatWithUsersByID(ctx context.Context, id int) (core.Chat, error) {
	var chat = core.Chat{ID: id, IsDeleted: false}
	err := s.DB.WithContext(ctx).
		Table("chats"). // Явно указываем, что работаем с таблицей "chats"
		Joins("JOIN chat_member cm ON cm.chat_id = chats.id").
		Joins("JOIN users u ON u.id = cm.user_id").
		Where("chats.id = ? AND cm.is_deleted = false", id).
		Preload("Users").
		First(&chat).Error
	if err != nil {
		return core.Chat{}, err
	}
	return chat, nil
}

func (s *store) CreateChat(ctx context.Context, data core.Chat) (core.Chat, error) {
	if err := s.DB.WithContext(ctx).Create(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}

func (s *store) FindChatByUsers(ctx context.Context, userIds []int) (core.Chat, error) {
	var chat core.Chat

	err := s.DB.WithContext(ctx).
		Table("chats").
		Joins("JOIN chat_member cm ON cm.chat_id = chats.id").
		Where("cm.user_id IN ?", userIds).
		Where("chats.is_deleted = false").
		Group("chats.id, cm.chat_id").
		Having("COUNT(DISTINCT cm.user_id) = ?", len(userIds)).
		Having("COUNT(DISTINCT cm.user_id) = (SELECT COUNT(*) FROM chat_member WHERE chat_member.chat_id = cm.chat_id)").
		First(&chat).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return core.Chat{ID: -1}, nil // Чат с такими пользователями не найден
		}
		return core.Chat{ID: -1}, err
	}

	return chat, nil
}

func (s *store) DeleteChat(ctx context.Context, id int) error {
	if err := ifChatExists(s, ctx, id); err != nil {
		return err
	}
	if err := s.DB.WithContext(ctx).Model(&core.Chat{}).Where("id", id).Where("is_deleted", false).Updates(map[string]interface{}{"is_deleted": true, "deleted_at": time.Now()}).Error; err != nil {
		return err
	}
	return nil
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
