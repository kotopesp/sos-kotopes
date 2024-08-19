package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/validator"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

type Router struct {
	app                  *fiber.App
	formValidator        validator.FormValidatorService
	authService          core.AuthService
	commentService       core.CommentService
	postService          core.PostService
	userService          core.UserService
	roleService          core.RoleService
	userFavouriteService core.UserFavouriteService
	keeperService        core.KeeperService
}

func NewRouter(
	app *fiber.App,
	authService core.AuthService,
	commentService core.CommentService,
	postService core.PostService,
	userService core.UserService,
	roleService core.RoleService,
	formValidator validator.FormValidatorService,
	keeperService core.KeeperService,
) {
	router := &Router{
		app:            app,
		formValidator:  formValidator,
		authService:    authService,
		postService:    postService,
		userService:    userService,
		roleService:    roleService,
		commentService: commentService,
		keeperService:  keeperService,
	}

	router.initRequestMiddlewares()

	router.initRoutes()

	router.initResponseMiddlewares()
}

func (r *Router) initRoutes() {
	r.app.Get("/ping", r.ping)

	v1 := r.app.Group("/api/v1")
	// comment_service
	v1.Get("/posts/:post_id/comments", r.getComments)
	v1.Post("/posts/:post_id/comments", r.protectedMiddleware(), r.createComment)
	v1.Patch("/posts/:post_id/comments/:comment_id", r.protectedMiddleware(), r.updateComment)
	v1.Delete("/posts/:post_id/comments/:comment_id", r.protectedMiddleware(), r.deleteComment)

	// favourites users todo
	v1.Get("/users/favourites", r.protectedMiddleware(), r.GetFavouriteUsers)
	v1.Post("/users/:id/favourites", r.AddUserToFavourites)
	v1.Delete("/users/:id/favourites", r.DeleteUserFromFavourites)

	// users
	v1.Get("/users/:id", r.getUser)
	v1.Patch("/users", r.protectedMiddleware(), r.updateUser)

	// user roles
	v1.Post("/users/roles", r.protectedMiddleware(), r.giveRoleToUser)
	v1.Get("/users/:id/roles", r.getUserRoles)
	v1.Patch("/users/roles", r.protectedMiddleware(), r.updateUserRoles)
	v1.Delete("/users/roles", r.protectedMiddleware(), r.deleteUserRole)

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
	v1.Put("/keepers/:id", r.protectedMiddleware(), r.updateKeeperByID)
	v1.Delete("/keepers/:id", r.protectedMiddleware(), r.deleteKeeperByID)

	// keeper reviews
	v1.Get("/keepers/:id/keeper_reviews", r.getKeeperReviews)
	v1.Post("/keepers/:id/keeper_reviews", r.protectedMiddleware(), r.createKeeperReview)
	v1.Put("/keeper_reviews/:id", r.protectedMiddleware(), r.updateKeeperReview)
	v1.Delete("/keeper_reviews/:id", r.protectedMiddleware(), r.deleteKeeperReview)

	// posts
	v1.Get("/posts", r.getPosts)
	v1.Get("/users/:id/posts", r.getUserPosts)
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
