package http

import (
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Router struct {
	app           *fiber.App
	entityService core.EntityService
	authService   interface{}
	chatService   core.ChatService
}

func NewRouter(
	app *fiber.App,
	entityService core.EntityService,
	authService interface{},
	chatService core.ChatService,
) {
	router := &Router{
		app:           app,
		entityService: entityService,
		authService:   authService,
		chatService:   chatService,
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
	v1.Get("/chats", r.getChats)
	v1.Get("/chats/:id", r.getChatByID)
	v1.Post("/chats", r.createChat)
	v1.Delete("/chats/:id", r.deleteChat)
}

// initRequestMiddlewares initializes all middlewares for http requests
func (r *Router) initRequestMiddlewares() {
	r.app.Use(logger.New())
}

// initResponseMiddlewares initializes all middlewares for http response
func (r *Router) initResponseMiddlewares() {}
