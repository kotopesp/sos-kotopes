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
	v1.Patch("/users/:id", r.UpdateUser)
	v1.Get("/users/:id", r.GetUser)
	v1.Get("/users/:id/posts", r.GetUserPosts)

	//favourites posts	todo
	//	v1.Delete("/posts/:id/favourites", r.DeletePostFromFavourites)
	//favourites users todo
	//	v1.Get("/users/favourites", r.GetUserFavourites)
	//	v1.Post("/users/:id/favourites", r.AddUserToFavourites)
	//	v1.Delete("/users/:id/favourites", r.DeleteUserFromFavourites)
	//favourites comments todo
	//	v1.Get("/posts/comments/favourites", r.GetFavouriteComments)
	//	v1.Delete("/posts/:id/comments/:id/favourites", r.DeleteCommentFromFavourites)
	// user roles todo
	//	v1.Get("/users/:id/roles", r.GetUserRoles)
	//	v1.Post("/users/:id/roles", r.GiveRoleToUser)
	//	v1.Patch("/users/:id/roles", r.UpdateUserRoles)
	//	v1.Delete("/users/:id/roles", r.DeleteUserRole)
	// reviews todo
	// v1.Get("/users/:id/reviews", r.GetUserReviews)
	//	v1.Post("/users/:id/reviews", r.ReviewUser)
	//	v1.Patch("/users/:id/reviews", r.UpdateReviewOnUser)
	//	v1.Delete("/users/:id/reviews", r.DeleteReviewOnUser)

}

// initRequestMiddlewares initializes all middlewares for http requests
func (r *Router) initRequestMiddlewares() {
	r.app.Use(logger.New())
}

// initResponseMiddlewares initializes all middlewares for http response
func (r *Router) initResponseMiddlewares() {}
