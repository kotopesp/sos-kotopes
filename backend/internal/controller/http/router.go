package http

import (
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core/user_core"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Router struct {
	app           *fiber.App
	entityService core.EntityService
	authService   interface{}
	userService   user_core.UserService
}

func NewRouter(
	app *fiber.App,
	entityService core.EntityService,
	authService interface{},
	userService user_core.UserService,
) {
	router := &Router{
		app:           app,
		entityService: entityService,
		authService:   authService,
		userService:   userService,
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

	// users
	v1.Patch("/users/:id", r.UpdateUser)
	v1.Get("/users/:id", r.GetUser)
	v1.Get("/users/:id/posts", r.GetUserPosts)

}

// initRequestMiddlewares initializes all middlewares for http requests
func (r *Router) initRequestMiddlewares() {
	r.app.Use(logger.New())
}

// initResponseMiddlewares initializes all middlewares for http response
func (r *Router) initResponseMiddlewares() {}
