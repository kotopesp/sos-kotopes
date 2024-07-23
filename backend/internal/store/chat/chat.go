package chat

import (
	"context"

	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/postgres"
)

type (
	store struct {
		*postgres.Postgres
	}
)

func New(pg *postgres.Postgres) core.ChatStore {
	return &store{pg}
}

func (s *store) GetAll(ctx context.Context, sortType string) ([]core.Chat, error) {
	var chats []core.Chat
	query := s.DB.WithContext(ctx).Model(&core.Chat{})
	if sortType != "" {
		query = query.Where(&core.Chat{Type: sortType})
	}
	if err := query.Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}

func (s *store) GetByID(ctx context.Context, id int) (core.Chat, error) {
	var chat core.Chat = core.Chat{ID: id}

	if err := s.DB.WithContext(ctx).First(&chat).Error; err != nil {
		return core.Chat{}, err
	}
	return chat, nil
}

func (s *store) Create(ctx context.Context, data core.Chat) (core.Chat, error) {
	if err := s.DB.WithContext(ctx).Create(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}

func (s *store) Delete(ctx context.Context, id int) error {
	if err := s.DB.WithContext(ctx).Delete(&core.Chat{ID: id}).Error; err != nil {
		return err
	}
	return nil
}
