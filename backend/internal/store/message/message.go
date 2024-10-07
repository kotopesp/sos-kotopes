package message

import (
	"context"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/chat"
	"github.com/kotopesp/sos-kotopes/internal/core"
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
	if err := s.DB.WithContext(ctx).
		Model(&core.Chat{}).
		Where("id", chatID).
		Where("is_deleted", false).
		Count(&counter).Error; err != nil {
		return err
	}
	if counter != 1 {
		return model.ErrInvalidChatID
	}
	return nil
}

func ifMessageExists(s *store, ctx context.Context, messageID int) error {
	var counter int64
	if err := s.DB.WithContext(ctx).
		Model(&core.Message{}).
		Where("id", messageID).
		Where("is_deleted", false).
		Count(&counter).Error; err != nil {
		return err
	}
	if counter != 1 {
		return model.ErrInvalidMessageID
	}
	return nil
}

func (s *store) MarkMessagesAsRead(ctx context.Context, chatID, userID int) error {
	var unreadMessages []core.Message

	err := s.DB.WithContext(ctx).
		Model(core.Message{}).
		Joins("LEFT JOIN message_read mr ON messages.id = mr.message_id AND mr.user_id = ?", userID).
		Where("messages.chat_id = ? AND mr.id IS NULL", chatID).
		Select("messages.id").
		Find(&unreadMessages).Error

	if err != nil {
		return err
	}

	if len(unreadMessages) == 0 {
		return nil
	}

	var reads []core.MessageRead
	for _, message := range unreadMessages {
		reads = append(reads, core.MessageRead{
			MessageID: message.ID,
			UserID:    userID,
			ReadAt:    time.Now().UTC(),
		})
	}

	err = s.DB.Create(&reads).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *store) GetUnreadMessageCount(ctx context.Context, chatID, userID int) (int64, error) {
	var count int64

	err := s.DB.WithContext(ctx).
		Model(&core.Message{}).
		Joins("LEFT JOIN message_read mr ON messages.id = mr.message_id AND mr.user_id = ?", userID).
		Where("messages.chat_id = ? AND mr.id IS NULL", chatID).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *store) GetAllMessages(ctx context.Context, chatID, userID int, sortType, searchText string) ([]chat.Message, error) {
	if err := ifChatExists(s, ctx, chatID); err != nil {
		return nil, err
	}
	var messages []core.Message
	query := s.DB.WithContext(ctx).
		Model(&core.Message{}).
		Where("chat_id", chatID).
		Where("is_deleted", false).
		Where("content LIKE ?", "%"+searchText+"%").
		Where("EXISTS (SELECT 1 FROM chat_members WHERE chat_members.chat_id = ? AND chat_members.user_id = ?)", chatID, userID).
		Order("created_at " + sortType)
	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}
	var messagesResponse []chat.Message
	for _, mes := range messages {
		messagesResponse = append(messagesResponse, chat.Message{
			UserID:     mes.UserID,
			ChatID:     mes.ChatID,
			Content:    mes.Content,
			CreatedAt:  mes.CreatedAt,
			SenderName: mes.SenderName,
		})
	}
	return messagesResponse, nil
}

func (s *store) CreateMessage(ctx context.Context, data chat.Message) (chat.Message, error) {
	var user core.User
	if err := s.DB.WithContext(ctx).
		Table("users").
		Where("id = ?", data.UserID).
		Where("EXISTS (SELECT 1 FROM chat_members WHERE chat_members.chat_id = ? AND chat_members.user_id = ?)", data.ChatID, data.UserID).
		First(&user).Error; err != nil {
		return chat.Message{}, err
	}

	dataToInsert := core.Message{
		UserID:     data.UserID,
		ChatID:     data.ChatID,
		Content:    data.Content,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
		SenderName: user.Username,
	}
	if err := ifChatExists(s, ctx, data.ChatID); err != nil {
		return chat.Message{}, err
	}
	if err := s.DB.WithContext(ctx).Create(&dataToInsert).Error; err != nil {
		return data, err
	}
	readMessage := core.MessageRead{
		MessageID: dataToInsert.ID,
		UserID:    dataToInsert.UserID,
		ReadAt:    time.Now().UTC(),
	}
	if err := s.DB.WithContext(ctx).Create(&readMessage).Error; err != nil {
		return chat.Message{}, err
	}
	data.SenderName = user.Username
	data.CreatedAt = dataToInsert.CreatedAt
	return data, nil
}

func (s *store) UpdateMessage(ctx context.Context, chatID, userID, messageID int, content string) (core.Message, error) {
	if err := ifChatExists(s, ctx, chatID); err != nil {
		return core.Message{}, err
	}
	if err := ifMessageExists(s, ctx, messageID); err != nil {
		return core.Message{}, err
	}
	if err := s.DB.WithContext(ctx).
		Model(&core.Message{}).
		Where("id", messageID).
		Where("chat_id", chatID).
		Where("is_deleted", false).
		Where("EXISTS (SELECT 1 FROM chat_members WHERE chat_members.chat_id = ? AND chat_members.user_id = ?)", chatID, userID).
		Updates(map[string]interface{}{"content": content, "updated_at": time.Now()}).Error; err != nil {
		return core.Message{}, err
	}
	var message core.Message
	s.DB.WithContext(ctx).Model(&core.Message{}).Where("id", messageID).First(&message)
	return message, nil
}

func (s *store) DeleteMessage(ctx context.Context, chatID, userID, messageID int) error {
	if err := ifChatExists(s, ctx, chatID); err != nil {
		return err
	}
	if err := ifMessageExists(s, ctx, messageID); err != nil {
		return err
	}
	if err := s.DB.WithContext(ctx).
		Model(&core.Message{}).
		Where("id", messageID).
		Where("chat_id", chatID).
		Where("is_deleted", false).
		Where("EXISTS (SELECT 1 FROM chat_members WHERE chat_members.chat_id = ? AND chat_members.user_id = ?)", chatID, userID).
		Updates(map[string]interface{}{"is_deleted": true, "deleted_at": time.Now()}).Error; err != nil {
		return err
	}
	return nil
}
