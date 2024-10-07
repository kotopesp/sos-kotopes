package websockets

import (
	"context"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

type WebSocketManager struct {
	connections map[int]map[*websocket.Conn]bool // мапа чатов, где для каждого чата хранятся соединения
	mu          sync.Mutex                       // мьютекс для синхронизации
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		connections: make(map[int]map[*websocket.Conn]bool),
	}
}

func (wsm *WebSocketManager) AddConnection(chatID int, conn *websocket.Conn) {
	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	if wsm.connections[chatID] == nil {
		wsm.connections[chatID] = make(map[*websocket.Conn]bool)
	}
	wsm.connections[chatID][conn] = true
}

func (wsm *WebSocketManager) RemoveConnection(chatID int, conn *websocket.Conn) {
	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	delete(wsm.connections[chatID], conn)
	if len(wsm.connections[chatID]) == 0 {
		delete(wsm.connections, chatID)
	}
}

func (wsm *WebSocketManager) BroadcastMessage(chatID int, message []byte) {
	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	for conn := range wsm.connections[chatID] {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Error sending message: %v", err)
			if err := conn.Close(); err != nil {
				logger.Log().Error(context.Background(), err.Error())
			}
			delete(wsm.connections[chatID], conn)
		}
	}
}

func (wsm *WebSocketManager) HandleWebSocket(c *fiber.Ctx) error {
	chatID, err := c.ParamsInt("chatID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid chat ID")
	}

	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("chatID", chatID)
		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}

func (wsm *WebSocketManager) WebSocketEndpoint(c *websocket.Conn) {
	chatID, ok := c.Locals("chatID").(int)
	if !ok {
		logger.Log().Error(context.Background(), model.ErrInvalidRequestParameter.Error())
	}

	// Добавляем соединение
	wsm.AddConnection(chatID, c)

	defer func() {
		// Удаляем соединение по завершению работы
		wsm.RemoveConnection(chatID, c)
		err := c.Close()
		if err != nil {
			logger.Log().Error(context.Background(), err.Error())
		}
	}()

	// Читаем и отправляем сообщение всем подключенным к чату
	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			logger.Log().Error(context.Background(), err.Error())
			break
		}

		if messageType == websocket.TextMessage {
			wsm.BroadcastMessage(chatID, message)
		}
	}
}
