package http

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/comment"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

func (r *Router) getComments(ctx *fiber.Ctx) error {
	var getAllCommentsParams comment.GetAllCommentsParams
	fiberError, parseOrValidationError := parseQueryAndValidate(ctx, r.formValidator, &getAllCommentsParams)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	var commentPathParams comment.PathParams
	fiberError, parseOrValidationError = parseParamsAndValidate(ctx, r.formValidator, &commentPathParams)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	coreComments, total, err := r.commentService.GetAllComments(
		ctx.UserContext(),
		getAllCommentsParams.ToCoreGetAllCommentsParams(commentPathParams.PostID),
	)
	if err != nil {
		if errors.Is(err, core.ErrPostNotFound) {
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}

		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	logger.Log().Debug(ctx.UserContext(), fmt.Sprintf("%v", coreComments[0].Author))

	modelComments := comment.ToModelCommentsSlice(coreComments)

	return ctx.Status(fiber.StatusOK).JSON(comment.ToGetAllCommentsResponse(
		modelComments,
		paginate(total, getAllCommentsParams.Limit, getAllCommentsParams.Offset),
	))
}

func (r *Router) createComment(ctx *fiber.Ctx) error {
	var comm comment.Comment
	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &comm)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	var pathParams comment.PathParams
	fiberError, parseOrValidationError = parseParamsAndValidate(ctx, r.formValidator, &pathParams)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	coreComment := comm.ToCoreComment()

	coreComment.PostID = pathParams.PostID
	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}
	coreComment.AuthorID = userID

	createdComment, err := r.commentService.CreateComment(ctx.UserContext(), coreComment)
	if err != nil {
		if errors.Is(err, core.ErrPostNotFound) {
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}

		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		model.OKResponse(comment.ToModelComment(createdComment)),
	)
}

func (r *Router) updateComment(ctx *fiber.Ctx) error {
	var pathParams comment.PathParams
	fiberError, parseOrValidationError := parseParamsAndValidate(ctx, r.formValidator, &pathParams)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	var newComment comment.Update
	fiberError, parseOrValidationError = parseBodyAndValidate(ctx, r.formValidator, &newComment)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	newCoreComment := newComment.ToCoreComment()

	newCoreComment.ID = pathParams.CommentID
	newCoreComment.PostID = pathParams.PostID
	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}
	newCoreComment.AuthorID = userID

	updatedComment, err := r.commentService.UpdateComment(ctx.UserContext(), newCoreComment)
	switch {
	case errors.Is(err, core.ErrCommentAuthorIDMismatch):
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusForbidden).JSON(model.ErrorResponse(err.Error()))
	case errors.Is(err, core.ErrCommentPostIDMismatch) ||
		errors.Is(err, core.ErrNoSuchComment) ||
		errors.Is(err, core.ErrCommentIsDeleted):
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
	case err != nil:
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(
		model.OKResponse(comment.ToModelComment(updatedComment)),
	)
}

func (r *Router) deleteComment(ctx *fiber.Ctx) error {
	var pathParams comment.PathParams
	fiberError, parseOrValidationError := parseParamsAndValidate(ctx, r.formValidator, &pathParams)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	var coreComment core.Comment
	coreComment.ID = pathParams.CommentID
	coreComment.PostID = pathParams.PostID
	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}
	coreComment.AuthorID = userID

	err = r.commentService.DeleteComment(ctx.UserContext(), coreComment)
	if errors.Is(err, core.ErrCommentAuthorIDMismatch) {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusForbidden).JSON(model.ErrorResponse(err.Error()))
	} else if err != nil &&
		!errors.Is(err, core.ErrNoSuchComment) &&
		!errors.Is(err, core.ErrCommentPostIDMismatch) &&
		!errors.Is(err, core.ErrCommentIsDeleted) {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
