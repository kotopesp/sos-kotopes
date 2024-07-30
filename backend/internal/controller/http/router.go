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
	authService          core.AuthService
	userService          core.UserService
	roleService          core.RoleService
	userFavouriteService core.UserFavouriteService
	formValidator        validator.FormValidatorService
}

func NewRouter(
	app *fiber.App,
	entityService core.EntityService,
	authService core.AuthService,
	userService core.UserService,
	roleService core.RoleService,
	userFavouriteService core.UserFavouriteService,
	formValidator validator.FormValidatorService,
) {
	router := &Router{
		app:                  app,
		entityService:        entityService,
		formValidator:        formValidator,
		authService:          authService,
		userService:          userService,
		roleService:          roleService,
		userFavouriteService: userFavouriteService,
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
	v1.Patch("/users", r.protectedMiddleware(), r.UpdateUser)
	v1.Get("/users/:id", r.GetUser)
	v1.Get("/users/:id/posts", r.GetUserPosts)

	// user roles
	v1.Get("/users/:id/roles", r.GetUserRoles)
	v1.Post("/users/:id/roles", r.GiveRoleToUser)
	v1.Delete("/users/:id/roles", r.DeleteUserRole)
	v1.Patch("/users/:id/roles", r.UpdateUserRoles)

	//favourites users todo
	v1.Get("/users/favourites", r.protectedMiddleware(), r.GetFavouriteUsers)
	v1.Post("/users/:id/favourites", r.AddUserToFavourites)
	v1.Delete("/users/:id/favourites", r.DeleteUserFromFavourites)

	// e.g. protected resource
	v1.Get("/protected", r.protectedMiddleware(), r.protected)

	// auth
	v1.Post("/auth/login", r.loginBasic)
	v1.Post("/auth/signup", r.signup)
	v1.Post("/auth/token/refresh", r.refreshTokenMiddleware(), r.refresh)

	// auth vk
	v1.Get("/auth/login/vk", r.loginVK)
	v1.Get("/auth/login/vk/callback", r.callback)
}

// initRequestMiddlewares initializes all middlewares for http requests
func (r *Router) initRequestMiddlewares() {
	r.app.Use(logger.New())
}

// initResponseMiddlewares initializes all middlewares for http response
func (r *Router) initResponseMiddlewares() {}
