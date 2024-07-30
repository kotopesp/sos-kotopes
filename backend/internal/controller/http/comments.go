package http

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/comments"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

func (r *Router) getCommentsByPostID(ctx *fiber.Ctx) error {
	var getAllCommentsParams comments.GetAllCommentsParams
	fiberError, parseOrValidationError := parseQueryAndValidate(ctx, r.formValidator, &getAllCommentsParams)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	var urlParams comments.CommentURLParams
	fiberError, parseOrValidationError = parseParamsAndValidate(ctx, r.formValidator, &urlParams)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	coreComments, total, err := r.commentsService.GetCommentsByPostID(
		ctx.UserContext(),
		getAllCommentsParams.ToCoreGetAllCommentsParams(),
		urlParams.PostID,
	)
	if err != nil {
		if errors.Is(err, core.ErrNoSuchComment) {
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(err.Error())
		}

		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	modelComments := comments.ToModelCommentsSlice(coreComments)

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(fiber.Map{
		"comments": modelComments,
		"meta":     paginate(total, getAllCommentsParams.Limit, getAllCommentsParams.Offset),
	}))
}

func (r *Router) createComment(ctx *fiber.Ctx) error {
	var comment comments.Comments
	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &comment)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	var urlParams comments.CommentURLParams
	fiberError, parseOrValidationError = parseParamsAndValidate(ctx, r.formValidator, &urlParams)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	coreComment := comment.ToCoreComments()

	coreComment.PostsID = urlParams.PostID
	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}
	coreComment.AuthorID = userID

	createdComment, err := r.commentsService.CreateComment(ctx.UserContext(), coreComment)
	if err != nil {
		if errors.Is(err, core.ErrNoSuchPost) {
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}

		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	logger.Log().Debug(ctx.UserContext(), fmt.Sprintf("Created comments: %v", createdComment))

	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(fiber.Map{
		"created_comment": comments.ToModelComment(createdComment),
	}))
}

func (r *Router) updateComment(ctx *fiber.Ctx) error {
	var urlParams comments.CommentURLParams
	fiberError, parseOrValidationError := parseParamsAndValidate(ctx, r.formValidator, &urlParams)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	var newComment comments.CommentUpdate
	fiberError, parseOrValidationError = parseBodyAndValidate(ctx, r.formValidator, &newComment)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	newCoreComment := newComment.ToCoreComment()

	newCoreComment.ID = urlParams.CommentID
	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}
	newCoreComment.AuthorID = userID

	updatedComment, err := r.commentsService.UpdateComments(ctx.UserContext(), newCoreComment)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		if errors.Is(err, core.ErrCommentAuthorIDMismatch) {
			return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
		} else if errors.Is(err, core.ErrNoSuchComment) || errors.Is(err, core.ErrCommentIsDeleted) {
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}

		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	logger.Log().Debug(ctx.UserContext(), fmt.Sprintf("Updated comments: %v", updatedComment))

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(fiber.Map{
		"updated_comment": comments.ToModelComment(updatedComment),
	}))
}

func (r *Router) deleteComment(ctx *fiber.Ctx) error {
	var urlParams comments.CommentURLParams
	fiberError, parseOrValidationError := parseParamsAndValidate(ctx, r.formValidator, &urlParams)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	var coreComment core.Comments
	coreComment.ID = urlParams.CommentID

	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}
	coreComment.AuthorID = userID

	err = r.commentsService.DeleteComments(ctx.UserContext(), coreComment)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		if errors.Is(err, core.ErrCommentAuthorIDMismatch) {
			return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
		} else if errors.Is(err, core.ErrNoSuchComment) || errors.Is(err, core.ErrCommentIsDeleted) {
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}

		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
