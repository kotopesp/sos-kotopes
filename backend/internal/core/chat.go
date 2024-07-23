package core

import (
	"context"
	"time"
)

type (
	Chat struct {
		ID        int       `gorm:"column:id" json:"id"`
		Type      string    `gorm:"column:type" json:"type"`
		CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	}

	ChatStore interface {
		GetAll(ctx context.Context, sortType string) (chats []Chat, err error)
		GetByID(ctx context.Context, id int) (chat Chat, err error)
		Create(ctx context.Context, data Chat) (chat Chat, err error)
		Delete(ctx context.Context, id int) error
	}

	ChatService interface {
		GetAll(ctx context.Context, sortType string) (chats []Chat, total int, err error)
		GetByID(ctx context.Context, id int) (chat Chat, err error)
		Create(ctx context.Context, data Chat) (chat Chat, err error)
		Delete(ctx context.Context, id int) error
	}
)
