package http

import (
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"github.com/gofiber/fiber/v2"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/logger"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
)

func (r *Router) getPosts(ctx *fiber.Ctx) error {
	sortBy := ctx.Query("sortBy")
	sortOrder := ctx.Query("sortOrder")
	searchTerm := ctx.Query("searchTerm")
	limit := ctx.QueryInt("limit")
	offset := ctx.QueryInt("offset")

	params := core.GetAllPostsParams{
		SortBy:     &sortBy,
		SortOrder:  &sortOrder,
		SearchTerm: &searchTerm,
		Limit:      &limit,
		Offset:     &offset,
	}

	posts, total, err := r.postService.GetAll(ctx.UserContext(), params)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	response := struct {
		Total int         `json:"total"`
		Posts []core.Post `json:"posts"`
	}{
		Total: total,
		Posts: posts,
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(response))
}

func (r *Router) getPostByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid post ID"))
	}

	post, err := r.postService.GetByID(ctx.UserContext(), id)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(post))
}

func (r *Router) createPost(ctx *fiber.Ctx) error {
	var post core.Post
	if err := ctx.BodyParser(&post); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid input"))
	}

	createdPost, err := r.postService.Create(ctx.UserContext(), post)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(createdPost))
}

func (r *Router) updatePost(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid post ID"))
	}

	var post core.Post
	if err := ctx.BodyParser(&post); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid input"))
	}
	post.ID = id

	updatedPost, err := r.postService.Update(ctx.UserContext(), post)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(updatedPost))
}

func (r *Router) deletePost(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid post ID"))
	}

	err = r.postService.Delete(ctx.UserContext(), id)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse("Post deleted"))
}
