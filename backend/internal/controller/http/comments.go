package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/comments"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

func (r *Router) getCommentsByPostID(ctx *fiber.Ctx) error {

	params := comments.GetAllCommentsParams{}
	if err := ctx.QueryParser(&params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	postsID, err := ctx.ParamsInt("postsID")
	if err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON("")
	}

	commentsByPost, err := r.commentsService.GetCommentsByPostID(ctx.UserContext(),
		*params.ToCoreGetAllCommentsParams(), postsID)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON("")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(commentsByPost))
}
func (r *Router) createComment(ctx *fiber.Ctx) error {

	comment := comments.Comments{}
	if err := ctx.BodyParser(&comment); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	coreComment := comment.ToCoreComments()

	postID, err := ctx.ParamsInt("postsID")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("invalid post id for comments"))
	}

	coreComment.PostsID = postID

	createdComment, err := r.commentsService.CreateComment(ctx.UserContext(), *coreComment, postID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(createdComment))
}

func (r *Router) updateComment(ctx *fiber.Ctx) error {

	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("invalid comment id"))
	}

	newComment := core.Comments{}
	if err := ctx.BodyParser(&newComment); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	newComment.ID = id

	updatedComment, err := r.commentsService.UpdateComments(ctx.UserContext(), newComment)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(
		model.Response{Data: updatedComment}))
}

func (r *Router) deleteComment(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid comment id")
	}

	err = r.commentsService.DeleteComments(ctx.UserContext(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return nil
}
