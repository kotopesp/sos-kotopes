package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

type Router struct {
	app                  *fiber.App
	entityService        core.EntityService
	authService          interface{}
	postService          core.PostService
	postFavouriteService core.PostFavoriteService
	postResponseService  core.PostResponseService
}

func NewRouter(
	app *fiber.App,
	entityService core.EntityService,
	authService interface{},
	postService core.PostService,
	postFavouriteService core.PostFavoriteService,
	postResponseService core.PostResponseService,
) {
	router := &Router{
		app:                  app,
		entityService:        entityService,
		authService:          authService,
		postService:          postService,
		postFavouriteService: postFavouriteService,
		postResponseService:  postResponseService,
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
	v1.Get("/posts/favorites", r.getFavoritePosts) // получает все посты у user (могут быть коллизии с "/posts/:id")
	v1.Get("/posts/:id", r.getPostByID)
	v1.Get("/posts/:id/photo", r.getPostPhoto)
	v1.Post("/posts", r.createPost)
	v1.Put("/posts/:id", r.updatePost)
	v1.Delete("/posts/:id", r.deletePost)

	// favorites posts

	v1.Get("/posts/favorites/:id", r.getFavoritePostUserByID)
	v1.Post("/posts/:id/favorites", r.addFavoritePost)
	// v1.Delete("/posts/favorites", r.deleteFavoriteAllPostsFromUser) // удалить все посты у user (не знаю нужна ли)
	v1.Delete("/posts/favorites/:id", r.deleteFavoritePostByID)

	// responses post
	v1.Post("/posts/:post_id/responses", r.createPostResponse)
	v1.Get("/posts/:post_id/responses", r.getResponsesByPostID)
	v1.Put("/posts/:post_id/responses", r.updatePostResponse)
	v1.Delete("/posts/:post_id/responses", r.deletePostResponse)
}

// initRequestMiddlewares initializes all middlewares for http requests
func (r *Router) initRequestMiddlewares() {
	r.app.Use(logger.New())
}

// initResponseMiddlewares initializes all middlewares for http response
func (r *Router) initResponseMiddlewares() {}
