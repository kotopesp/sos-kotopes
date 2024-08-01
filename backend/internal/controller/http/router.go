package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/validator"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

type Router struct {
	app                  *fiber.App
	entityService        core.EntityService
	formValidator        validator.FormValidatorService
	authService          core.AuthService
	keeperService        core.KeeperService
	keeperReviewsService core.KeeperReviewsService
}

func NewRouter(
	app *fiber.App,
	entityService core.EntityService,
	keeperService core.KeeperService,
	keeperReviewsService core.KeeperReviewsService,
	formValidator validator.FormValidatorService,
	authService core.AuthService,
) {
	router := &Router{
		app:           app,
		entityService: entityService,
		keeperService: keeperService,
		formValidator: formValidator,
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

	// e.g. protected resource
	v1.Get("/protected", r.protectedMiddleware(), r.protected)

	// auth
	v1.Post("/auth/login", r.loginBasic)
	v1.Post("/auth/signup", r.signup)
	v1.Post("/auth/token/refresh", r.refreshTokenMiddleware(), r.refresh)

	// auth vk
	v1.Get("/auth/login/vk", r.loginVK)
	v1.Get("/auth/login/vk/callback", r.callback)

	// keepers
	v1.Get("/keepers", r.getKeepers)
	v1.Get("/keepers/:id", r.getKeeperByID)
	v1.Post("/keepers", r.protectedMiddleware(), r.createKeeper)
	v1.Put("/keepers/:id", r.updateKeeperByID)
	v1.Delete("/keepers/:id", r.deleteKeeperByID)

	// keeper reviews
	v1.Get("/keepers/:id/keeper_reviews", r.getKeeperReviews)
	v1.Post("/keepers/:id/keeper_reviews", r.protectedMiddleware(), r.createKeeperReview)
	v1.Put("/keeper_reviews/:id", r.updateKeeperReview)
	v1.Delete("/keeper_reviews/:id", r.deleteKeeperReview)
}

// initRequestMiddlewares initializes all middlewares for http requests
func (r *Router) initRequestMiddlewares() {
	r.app.Use(logger.New())
}

// initResponseMiddlewares initializes all middlewares for http response
func (r *Router) initResponseMiddlewares() {}
