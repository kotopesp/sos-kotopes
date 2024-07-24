package message

import (
	"context"
	"time"

	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/postgres"
)

type (
	store struct {
		*postgres.Postgres
	}
)

func New(pg *postgres.Postgres) core.MessageStore {
	return &store{pg}
}

func (s *store) GetAllMessages(ctx context.Context, id int, sortType string, searchText string) ([]core.Message, error) {
	var messages []core.Message
	query := s.DB.WithContext(ctx).Model(&core.Message{}).Where(&core.Message{ChatID: id}).Where("is_deleted", false).Where("content LIKE ?", "%"+searchText+"%").Order("created_at " + sortType)
	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *store) CreateMessage(ctx context.Context, data core.Message) (core.Message, error) {
	if err := s.DB.WithContext(ctx).Create(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}

func (s *store) UpdateMessage(ctx context.Context, chatId int, messageId int, content string) (core.Message, error) {
	if err := s.DB.WithContext(ctx).Model(&core.Message{}).Where(&core.Message{ID: messageId, ChatID: chatId}).Updates(map[string]interface{}{"content": content, "updated_at": time.Now()}).Error; err != nil {
		return core.Message{}, err
	}
	var message core.Message
	s.DB.WithContext(ctx).Model(&core.Message{}).Where(&core.Message{ID: messageId}).First(&message)
	return message, nil
}

func (s *store) DeleteMessage(ctx context.Context, chatId int, messageId int) error {
	if err := s.DB.WithContext(ctx).Model(&core.Message{}).Where(&core.Message{ID: messageId, ChatID: chatId}).Updates(map[string]interface{}{"is_deleted": true, "deleted_at": time.Now()}).Error; err != nil {
		return err
	}
	return nil
}
