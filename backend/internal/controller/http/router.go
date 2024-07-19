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
}

func NewRouter(
	app *fiber.App,
	entityService core.EntityService,
	authService interface{},
) {
	router := &Router{
		app:           app,
		entityService: entityService,
		authService:   authService,
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

	// keepers
	v1.Get("/keepers", r.getKeepers)
	v1.Get("/keepers/:id", r.getKeeperByID)
	v1.Post("/keepers", r.createKeeper)
	v1.Put("/keepers/:id", r.updateKeeper)
	v1.Delete("/keepers/:id", r.deleteKeeper)
}

// initRequestMiddlewares initializes all middlewares for http requests
func (r *Router) initRequestMiddlewares() {
	r.app.Use(logger.New())
}

// initResponseMiddlewares initializes all middlewares for http response
func (r *Router) initResponseMiddlewares() {}
