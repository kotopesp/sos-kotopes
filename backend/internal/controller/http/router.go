package http

import (
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Router struct {
	app           *fiber.App
	entityService core.EntityService
	authService   core.AuthService
}

func NewRouter(
	app *fiber.App,
	entityService core.EntityService,
	authService core.AuthService,
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

	// e.g protected resource
	v1.Get("/protected", r.protected)

	// auth
	v1.Post("/auth/login", r.login)
	v1.Post("/auth/signup", r.signup)
	v1.Post("/auth/token/refresh", r.refresh)
}

// initRequestMiddlewares initializes all middlewares for http requests
func (r *Router) initRequestMiddlewares() {
	r.app.Use(logger.New())

	v1 := r.app.Group("/api/v1")
	// protected paths (need to have access token)
	protectedPaths := []string{
		"/protected", // e.g, need access token to access /api/v1/protected
	}
	// or in initRoutes: `v1.Get("/protected", r.protectedMiddleware(), r.protected)`
	v1.Use(protectedPaths, r.protectedMiddleware())

	// refresh token middleware
	refreshTokenPath := "/auth/token/refresh"
	v1.Use(refreshTokenPath, r.refreshTokenMiddleware())
}

// initResponseMiddlewares initializes all middlewares for http response
func (r *Router) initResponseMiddlewares() {}
