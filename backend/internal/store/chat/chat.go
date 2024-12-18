package chat

import (
	"context"
	"errors"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/chat"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
	"gorm.io/gorm"
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

func convertToModelUser(users []core.User) (userResponse []chat.User) {
	for _, u := range users {
		userResponse = append(userResponse, chat.User{
			ID:       u.ID,
			Username: u.Username,
		})
	}
	return
}

func convertToModelChat(coreChat core.Chat, users []chat.User, lastMessage chat.Message, unreadCount int) chat.Chat {
	return chat.Chat{
		ID:          coreChat.ID,
		ChatType:    coreChat.ChatType,
		IsDeleted:   coreChat.IsDeleted,
		DeletedAt:   coreChat.DeletedAt,
		CreatedAt:   coreChat.CreatedAt,
		UpdatedAt:   coreChat.UpdatedAt,
		Users:       users,
		LastMessage: lastMessage,
		UnreadCount: unreadCount,
	}
}

func convertToModelLastMessage(message core.Message) chat.Message {
	return chat.Message{
		UserID:      message.UserID,
		TextContent: message.Content,
		CreatedAt:   message.CreatedAt,
		SenderName:  message.SenderName,
	}
}

func (s *store) GetUnreadMessageCount(ctx context.Context, chatID, userID int) (int64, error) {
	var count int64

	err := s.DB.WithContext(ctx).
		Model(&core.Message{}).
		Joins("LEFT JOIN message_read mr ON messages.id = mr.message_id AND mr.user_id = ?", userID).
		Where("messages.chat_id = ? AND mr.id IS NULL", chatID).
		Where("EXISTS (SELECT 1 FROM chat_members WHERE chat_members.chat_id = ? AND chat_members.user_id = ?)", chatID, userID).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *store) GetAllChats(ctx context.Context, sortType string, userID int) ([]chat.Chat, error) {
	query := s.DB.WithContext(ctx).
		Model(&core.Chat{IsDeleted: false}).
		Joins("JOIN chat_members ON chat_members.chat_id = chats.id").
		Where("chat_members.user_id = ?", userID).
		Where("chat_members.is_deleted = false")

	var chats []core.Chat
	if sortType != "" {
		query = query.Where("chat_type", sortType)
	}
	if err := query.
		Preload("Users").
		Find(&chats).Error; err != nil {
		return nil, err
	}

	var chatResponses []chat.Chat

	for _, chatEl := range chats {
		var message core.Message
		err := s.DB.WithContext(ctx).
			Model(&core.Message{}).
			Where("is_deleted = false").
			Where("chat_id = ?", chatEl.ID).
			Order("created_at desc").
			Limit(1).
			Find(&message).Error
		if err != nil {
			message = core.Message{}
		}
		unreadCount, err := s.GetUnreadMessageCount(ctx, chatEl.ID, userID)
		if err != nil {
			return nil, err
		}

		usersResponse := convertToModelUser(chatEl.Users)
		lastMessage := convertToModelLastMessage(message)

		chatResponses = append(chatResponses, convertToModelChat(chatEl, usersResponse, lastMessage, int(unreadCount)))
	}

	return chatResponses, nil
}

func (s *store) GetChatWithUsersByID(ctx context.Context, chatID, userID int) (chat.Chat, error) {
	var foundChat = core.Chat{ID: chatID, IsDeleted: false}
	err := s.DB.WithContext(ctx).
		Table("chats").
		Joins("JOIN chat_members cm ON cm.chat_id = chats.id").
		Joins("JOIN users u ON u.id = cm.user_id").
		Where("chats.id = ? AND cm.is_deleted = false", chatID).
		Where("EXISTS (SELECT 1 FROM chat_members WHERE chat_members.chat_id = ? AND chat_members.user_id = ?)", chatID, userID).
		Preload("Users").
		First(&foundChat).Error
	if err != nil {
		return chat.Chat{}, err
	}

	usersResponse := convertToModelUser(foundChat.Users)

	foundChatResponse := convertToModelChat(foundChat, usersResponse, chat.Message{}, 0)

	return foundChatResponse, nil
}

func (s *store) CreateChat(ctx context.Context, data chat.Chat) (chat.Chat, error) {
	dataToInsert := core.Chat{
		ChatType:  data.ChatType,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err := s.DB.WithContext(ctx).
		Create(&dataToInsert).Error; err != nil {
		return data, err
	}
	data.ID = dataToInsert.ID
	return data, nil
}

func (s *store) FindChatByUsers(ctx context.Context, userIds []int, chatType string) (chat.Chat, error) {
	var foundChat core.Chat

	err := s.DB.WithContext(ctx).
		Table("chats").
		Joins("JOIN chat_members cm ON cm.chat_id = chats.id").
		Where("cm.user_id IN ?", userIds).
		Where("chats.is_deleted = false").
		Where("chats.chat_type = ?", chatType).
		Group("chats.id, cm.chat_id").
		Having("COUNT(DISTINCT cm.user_id) = ?", len(userIds)).
		Having("COUNT(DISTINCT cm.user_id) = (SELECT COUNT(*) FROM chat_members WHERE chat_members.chat_id = cm.chat_id)").
		First(&foundChat).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return chat.Chat{ID: -1}, nil
		}
		return chat.Chat{ID: -1}, err
	}

	usersResponse := convertToModelUser(foundChat.Users)
	chatResponse := convertToModelChat(foundChat, usersResponse, chat.Message{}, 0)

	return chatResponse, nil
}

func (s *store) DeleteChat(ctx context.Context, chatID, userID int) error {
	if err := ifChatExists(s, ctx, chatID); err != nil {
		return err
	}
	if err := s.DB.WithContext(ctx).
		Model(&core.Chat{}).
		Where("id", chatID).
		Where("is_deleted", false).
		Where("EXISTS (SELECT 1 FROM chat_members WHERE chat_members.chat_id = ? AND chat_members.user_id = ?)", chatID, userID).
		Updates(map[string]interface{}{"is_deleted": true, "deleted_at": time.Now()}).Error; err != nil {
		return err
	}
	return nil
}

func (s *store) AddMemberToChat(ctx context.Context, data core.ChatMember) (core.ChatMember, error) {
	if err := ifChatExists(s, ctx, data.ChatID); err != nil {
		return core.ChatMember{}, err
	}
	if err := s.DB.WithContext(ctx).
		Create(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}
