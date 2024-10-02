package chat

import (
	"time"
)

type (
	Chat struct {
		ID          int       `json:"id"`
		ChatType    string    `json:"chat_type"`
		IsDeleted   bool      `json:"is_deleted"`
		DeletedAt   time.Time `json:"deleted_at"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Users       []User    `json:"users"`
		LastMessage Message   `json:"last_message"`
		UnreadCount int       `json:"unread_count"`
	}

	CreateChatRequest struct {
		ChatType string `json:"chat_type"`
		UserIDs  []int  `json:"users"`
	}

	User struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		// Firstname    *string    `json:"firstname"`
		// Lastname     *string    `json:"lastname"`
		// Photo        *[]byte    `json:"photo"`
		// PasswordHash string     `json:"password_hash"`
		// Description  *string    `json:"description"`
		// IsDeleted    bool       `json:"is_deleted"`
		// CreatedAt    time.Time  `json:"created_at"`
		// UpdatedAt    time.Time  `json:"updated_at"`
		// DeletedAt    *time.Time `json:"deleted_at"`
	}
)
