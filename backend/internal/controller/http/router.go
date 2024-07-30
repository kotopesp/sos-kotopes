package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/validator"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

type Router struct {
	app           *fiber.App
	entityService core.EntityService
	formValidator validator.FormValidatorService
	authService   core.AuthService
	postService   core.PostService
}

func NewRouter(
	app *fiber.App,
	entityService core.EntityService,
	formValidator validator.FormValidatorService,
	authService core.AuthService,
	postService core.PostService,
) {
	router := &Router{
		app:           app,
		entityService: entityService,
		formValidator: formValidator,
		authService:   authService,
		postService:   postService,
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

	// posts
	v1.Get("/posts", r.getPosts)
	v1.Get("/posts/favourites", r.protectedMiddleware(), r.getFavoritePostsUserByID) // gets all favourite posts from the user (there may be collisions with "/posts/:id")
	v1.Get("/posts/:id", r.getPostByID)
	v1.Post("/posts", r.protectedMiddleware(), r.createPost)
	v1.Patch("/posts/:id", r.updatePost)
	v1.Delete("/posts/:id", r.deletePost)

	// favorites posts
	v1.Get("/posts/favourites/:id", r.protectedMiddleware(), r.getFavoritePostUserAndPostByID)
	v1.Post("/posts/:id/favourites", r.protectedMiddleware(), r.addFavoritePost)
	v1.Delete("/posts/favourites/:id", r.protectedMiddleware(), r.deleteFavoritePostByID)
}

// initRequestMiddlewares initializes all middlewares for http requests
func (r *Router) initRequestMiddlewares() {
	r.app.Use(logger.New())
}

// initResponseMiddlewares initializes all middlewares for http response
func (r *Router) initResponseMiddlewares() {}
