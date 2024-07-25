package chat

import (
	"context"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/internal/store/errors"
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
		return errors.ErrInvalidChatID
	}
	return nil
}

func (s *store) GetAllChats(ctx context.Context, sortType string) ([]core.Chat, error) {
	var chats []core.Chat
	query := s.DB.WithContext(ctx).Model(&core.Chat{}).Where("is_deleted", false)
	if sortType != "" {
		query = query.Where("chat_type", sortType)
	}
	if err := query.Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}

func (s *store) GetChatByID(ctx context.Context, id int) (core.Chat, error) {
	var chat = core.Chat{ID: id, IsDeleted: false}

	if err := s.DB.WithContext(ctx).First(&chat).Error; err != nil {
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

func (s *store) DeleteChat(ctx context.Context, id int) error {
	if err := ifChatExists(s, ctx, id); err != nil {
		return err
	}
	if err := s.DB.WithContext(ctx).Model(&core.Chat{}).Where("id", id).Where("is_deleted", false).Updates(map[string]interface{}{"is_deleted": true, "deleted_at": time.Now()}).Error; err != nil {
		return err
	}
	return nil
}
