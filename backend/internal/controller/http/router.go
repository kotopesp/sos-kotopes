package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/validator"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

type Router struct {
	app           *fiber.App
	formValidator validator.FormValidatorService
	authService   core.AuthService
	postService   core.PostService
	keeperService core.KeeperService
}

func NewRouter(
	app *fiber.App,
	formValidator validator.FormValidatorService,
	authService core.AuthService,
	postService core.PostService,
	keeperService core.KeeperService,
) {
	router := &Router{
		app:           app,
		keeperService: keeperService,
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
	v1.Put("/keepers/:id", r.updateKeeperByID)
	v1.Delete("/keepers/:id", r.deleteKeeperByID)

	// keeper reviews
	v1.Get("/keepers/:id/keeper_reviews", r.getKeeperReviews)
	v1.Post("/keepers/:id/keeper_reviews", r.protectedMiddleware(), r.createKeeperReview)
	v1.Put("/keeper_reviews/:id", r.updateKeeperReview)
	v1.Delete("/keeper_reviews/:id", r.deleteKeeperReview)

	// posts
	v1.Get("/posts", r.getPosts)
	v1.Get("/posts/favourites", r.protectedMiddleware(), r.getFavouritePostsUserByID) // gets all favourite posts from the user (there may be collisions with "/posts/:id")
	v1.Get("/posts/:id", r.getPostByID)
	v1.Post("/posts", r.protectedMiddleware(), r.createPost)
	v1.Patch("/posts/:id", r.protectedMiddleware(), r.updatePost)
	v1.Delete("/posts/:id", r.protectedMiddleware(), r.deletePost)

	// favourites posts
	v1.Post("/posts/:id/favourites", r.protectedMiddleware(), r.addFavouritePost)
	v1.Delete("/posts/favourites/:id", r.protectedMiddleware(), r.deleteFavouritePostByID)
}

// initRequestMiddlewares initializes all middlewares for http requests
func (r *Router) initRequestMiddlewares() {
	r.app.Use(logger.New())
}

// initResponseMiddlewares initializes all middlewares for http response
func (r *Router) initResponseMiddlewares() {}
