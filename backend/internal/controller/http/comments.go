package http

import (
	"strconv"

	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/comments"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"github.com/gofiber/fiber/v2"
)

func (r *Router) getCommentsByPostID(ctx *fiber.Ctx) error {

	var params = comments.GetAllCommentsParams{}

	if err := ctx.QueryParser(&params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	posts_id, err := strconv.Atoi(ctx.FormValue("postsID"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ctx.FormValue("postsID"))
	}

	comments, err := r.commentsService.GetCommentsByPostID(ctx.UserContext(),
		*params.ToCoreGetAllCommentsParams(), posts_id)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(
		model.Response{
			Data: comments,
		}))
}
func (r *Router) createComment(ctx *fiber.Ctx) error {

	comment := comments.Comments{}

	if err := ctx.BodyParser(&comment); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	coreComment := comment.ToCoreComments()

	post_id, err := strconv.Atoi(ctx.FormValue("postsID")) //ctx.ParamsInt("postsID")

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("invalid post id for comments"))
	}

	coreComment.Posts_id = post_id

	createdComment, err := r.commentsService.CreateComment(ctx.UserContext(), *coreComment, post_id)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(createdComment))
}

func (r *Router) updateComment(ctx *fiber.Ctx) error {

	id, err := strconv.Atoi(ctx.FormValue("id"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("invalid comment id"))
	}

	newComment := core.Comments{}

	if err := ctx.BodyParser(&newComment); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	newComment.Id = id

	updatedComment, err := r.commentsService.UpdateComments(ctx.UserContext(), newComment)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(
		model.Response{Data: updatedComment}))
}

func (r *Router) deleteComment(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.FormValue("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid comment id")
	}

	err = r.commentsService.DeleteComments(ctx.UserContext(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return nil
}
