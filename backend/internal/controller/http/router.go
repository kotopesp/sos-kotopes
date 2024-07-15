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
	userService   core.UserService
}

func NewRouter(
	app *fiber.App,
	entityService core.EntityService,
	authService interface{},
	userService core.UserService,
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
	//v1.Patch("/users", r.UpdatePhoto)
	//v1.Delete("/users", r.DeletePhoto)
	v1.Patch("/users", r.ChangeName)
	v1.Patch("/users", r.ChangeDescription)
	//v1.Post("/users", r.SendMessage)
	//v1.Post("/users", r.AddUserToFavourites)
	//v1.Post("/users", r.DeleteUserFromFavourites)
}

// initRequestMiddlewares initializes all middlewares for http requests
func (r *Router) initRequestMiddlewares() {
	r.app.Use(logger.New())
}

// initResponseMiddlewares initializes all middlewares for http response
func (r *Router) initResponseMiddlewares() {}
