package http

import (
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Router struct {
	app               *fiber.App
	entityService     core.EntityService
	authService       interface{}
	chatService       core.ChatService
	messageService    core.MessageService
	chatMemberService core.ChatMemberService
}

func NewRouter(
	app *fiber.App,
	entityService core.EntityService,
	authService interface{},
	chatService core.ChatService,
	messageService core.MessageService,
	chatMemberService core.ChatMemberService,
) {
	router := &Router{
		app:               app,
		entityService:     entityService,
		authService:       authService,
		chatService:       chatService,
		messageService:    messageService,
		chatMemberService: chatMemberService,
	}

	router.initRequestMiddlewares()

	router.initRoutes()

	router.initResponseMiddlewares()
}

func (r *Router) initRoutes() {
	r.app.Get("/ping", r.ping)

	v1 := r.app.Group("/api/v1")

	// entities
	v1.Get("/entities", r.getEntities)
	v1.Get("/entities/:id", r.getEntityByID)

	// chats
	v1.Get("/chats", r.getAllChats)
	v1.Get("/chats/:chat_id", r.getChatByID)
	v1.Post("/chats", r.createChat)
	v1.Delete("/chats/:chat_id", r.deleteChat)

	// messages
	v1.Get("/chats/:chat_id/messages", r.getAllMessages)
	v1.Post("/chats/:chat_id/messages", r.createMessage)
	v1.Patch("/chats/:chat_id/messages/:message_id", r.updateMessage)
	v1.Delete("/chats/:chat_id/messages/:message_id", r.deleteMessage)

	// chat members
	v1.Get("/chats/:chat_id/members", r.getAllMembers)
	v1.Post("/chats/:chat_id/members", r.addMemberToChat)
	v1.Patch("/chats/:chat_id/members/:user_id", r.updateMemberInfo)
	v1.Delete("/chats/:chat_id/members/:user_id", r.deleteMemberFromChat)
}

// initRequestMiddlewares initializes all middlewares for http requests
func (r *Router) initRequestMiddlewares() {
	r.app.Use(logger.New())
}

// initResponseMiddlewares initializes all middlewares for http response
func (r *Router) initResponseMiddlewares() {}
