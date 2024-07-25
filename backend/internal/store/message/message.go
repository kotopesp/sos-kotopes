package message

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

func New(pg *postgres.Postgres) core.MessageStore {
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

func ifMessageExists(s *store, ctx context.Context, messageID int) error {
	var counter int64
	if err := s.DB.WithContext(ctx).Model(&core.Message{}).Where("id", messageID).Where("is_deleted", false).Count(&counter).Error; err != nil {
		return err
	}
	if counter != 1 {
		return errors.ErrInvalidMessageID
	}
	return nil
}

func (s *store) GetAllMessages(ctx context.Context, id int, sortType, searchText string) ([]core.Message, error) {
	if err := ifChatExists(s, ctx, id); err != nil {
		return nil, err
	}
	var messages []core.Message
	query := s.DB.WithContext(ctx).Model(&core.Message{}).Where("chat_id", id).Where("is_deleted", false).Where("content LIKE ?", "%"+searchText+"%").Order("created_at " + sortType)
	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *store) CreateMessage(ctx context.Context, data core.Message) (core.Message, error) {
	if err := ifChatExists(s, ctx, data.ChatID); err != nil {
		return core.Message{}, err
	}
	if err := s.DB.WithContext(ctx).Create(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}

func (s *store) UpdateMessage(ctx context.Context, chatID, messageID int, content string) (core.Message, error) {
	if err := ifChatExists(s, ctx, chatID); err != nil {
		return core.Message{}, err
	}
	if err := ifMessageExists(s, ctx, messageID); err != nil {
		return core.Message{}, err
	}
	if err := s.DB.WithContext(ctx).Model(&core.Message{}).Where("id", messageID).Where("chat_id", chatID).Where("is_deleted", false).Updates(map[string]interface{}{"content": content, "updated_at": time.Now()}).Error; err != nil {
		return core.Message{}, err
	}
	var message core.Message
	s.DB.WithContext(ctx).Model(&core.Message{}).Where("id", messageID).First(&message)
	return message, nil
}

func (s *store) DeleteMessage(ctx context.Context, chatID, messageID int) error {
	if err := ifChatExists(s, ctx, chatID); err != nil {
		return err
	}
	if err := ifMessageExists(s, ctx, messageID); err != nil {
		return err
	}
	if err := s.DB.WithContext(ctx).Model(&core.Message{}).Where("id", messageID).Where("chat_id", chatID).Where("is_deleted", false).Updates(map[string]interface{}{"is_deleted": true, "deleted_at": time.Now()}).Error; err != nil {
		return err
	}
	return nil
}
