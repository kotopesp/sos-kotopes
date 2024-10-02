package chat

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/chat"
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
		UnreadCount: int(unreadCount),
	}
}

func convertToModelLastMessage(message core.Message) chat.Message {
	return chat.Message{
		UserID:     message.UserID,
		Content:    message.Content,
		CreatedAt:  message.CreatedAt,
		SenderName: message.SenderName,
	}
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

func (s *store) GetAllChats(ctx context.Context, sortType string, userID int) ([]chat.Chat, error) {
	// userID := 2

	// lastMessageSubQuery := s.DB.WithContext(ctx).Table("message").
	// 	Select("content, created_at").
	// 	Where("chat_id = chats.id").
	// 	Order("created_at DESC").
	// 	Limit(1)
	// unreadCountSubQuery := s.DB.WithContext(ctx).Table("message").
	// 	Select("COUNT(*)").
	// 	Where("chat_id = chats.id").
	// 	Where("is_read = false").
	// 	Where("user_id != ?", userID)

	// query := s.DB.WithContext(ctx).Table("chats").Where("is_deleted = false")

	query := s.DB.WithContext(ctx).
		Model(&core.Chat{IsDeleted: false}).
		Joins("JOIN chat_member ON chat_member.chat_id = chats.id").
		Where("chat_member.user_id = ?", userID).
		Where("chat_member.is_deleted = false")

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
		err := s.DB.WithContext(ctx).Model(&core.Message{}).Where("is_deleted = false").Where("chat_id = ?", chatEl.ID).Order("created_at desc").Limit(1).Find(&message).Error
		if err != nil {
			message = core.Message{}
		}

		unreadCount, err := s.GetUnreadMessageCount(ctx, chatEl.ID, userID)
		if err != nil {
			return nil, err
		}
		// err = s.DB.WithContext(ctx).Model(&core.Message{}).
		// 	Where("chat_id = ? AND is_read = false AND user_id != ?", chatEl.ID, userID).
		// 	Count(&unreadCount).Error
		// if err != nil {
		// 	return nil, err
		// }

		usersResponse := convertToModelUser(chatEl.Users)
		lastMessage := convertToModelLastMessage(message)

		chatResponses = append(chatResponses, convertToModelChat(chatEl, usersResponse, lastMessage, int(unreadCount)))
	}

	return chatResponses, nil
}

func (s *store) GetChatWithUsersByID(ctx context.Context, id int) (chat.Chat, error) {
	var foundChat = core.Chat{ID: id, IsDeleted: false}
	err := s.DB.WithContext(ctx).
		Table("chats").
		Joins("JOIN chat_member cm ON cm.chat_id = chats.id").
		Joins("JOIN users u ON u.id = cm.user_id").
		Where("chats.id = ? AND cm.is_deleted = false", id).
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
		ChatType: data.ChatType,
	}
	if err := s.DB.WithContext(ctx).Create(&dataToInsert).Error; err != nil {
		return data, err
	}
	data.ID = dataToInsert.ID
	return data, nil
}

func (s *store) FindChatByUsers(ctx context.Context, userIds []int) (chat.Chat, error) {
	var foundChat core.Chat

	err := s.DB.WithContext(ctx).
		Table("chats").
		Joins("JOIN chat_member cm ON cm.chat_id = chats.id").
		Where("cm.user_id IN ?", userIds).
		Where("chats.is_deleted = false").
		Group("chats.id, cm.chat_id").
		Having("COUNT(DISTINCT cm.user_id) = ?", len(userIds)).
		Having("COUNT(DISTINCT cm.user_id) = (SELECT COUNT(*) FROM chat_member WHERE chat_member.chat_id = cm.chat_id)").
		First(&foundChat).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return chat.Chat{ID: -1}, nil // Чат с такими пользователями не найден
		}
		return chat.Chat{ID: -1}, err
	}

	usersResponse := convertToModelUser(foundChat.Users)
	chatResponse := convertToModelChat(foundChat, usersResponse, chat.Message{}, 0)

	return chatResponse, nil
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
