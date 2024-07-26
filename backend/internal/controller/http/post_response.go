package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

func (r *Router) createPostResponse(ctx *fiber.Ctx) error {
	var apiPostResponse = core.PostResponse{}

	if err := ctx.BodyParser(&apiPostResponse); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid input"))
	}

	postID, err := ctx.ParamsInt("post_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid post ID"))
	}

	authorID, err := ctx.ParamsInt("author_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid author_id"))
	}

	apiPostResponse.PostID = postID
	apiPostResponse.AuthorID = authorID

	createdResponse, err := r.postResponseService.CreatePostResponse(ctx.UserContext(), apiPostResponse)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(createdResponse))
}

func (r *Router) getResponsesByPostID(ctx *fiber.Ctx) error {
	postID, err := ctx.ParamsInt("post_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid post ID"))
	}

	responses, err := r.postResponseService.GetResponsesByPostID(ctx.UserContext(), postID)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(responses))
}

func (r *Router) updatePostResponse(ctx *fiber.Ctx) error {
	var response core.PostResponse
	if err := ctx.BodyParser(&response); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid input"))
	}

	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid response ID"))
	}
	response.ID = id

	updatedResponse, err := r.postResponseService.UpdatePostResponse(ctx.UserContext(), response)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(updatedResponse))
}

func (r *Router) deletePostResponse(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid response ID"))
	}

	err = r.postResponseService.DeletePostResponse(ctx.UserContext(), id)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse("Response deleted"))
}
