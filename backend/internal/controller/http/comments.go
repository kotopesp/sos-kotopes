package http

import (
	"strconv"

	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/comments"
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
