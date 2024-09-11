package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/comment"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/validator"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

// @Summary		Get all comments
// @Tags			comments
// @Description	Get all comments for a post
// @ID				get-all-comments
// @Accept			json
// @Produce		json
// @Param			post_id	path		int	true	"Post ID"	minimum(1)
// @Param			limit	query		int	true	"Limit"		minimum(1)
// @Param			offset	query		int	true	"Offset"	minimum(0)
// @Success		200		{object}	comment.GetAllCommentsResponse{data=[]comment.Comment}
// @Failure		400		{object}	model.Response
// @Failure		404		{object}	model.Response
// @Failure		500		{object}	model.Response
// @Router			/posts/{post_id}/comments [get]
func (r *Router) getComments(ctx *fiber.Ctx) error {
	var getAllCommentsParams comment.GetAllCommentsParams
	fiberError, parseOrValidationError := parseQueryAndValidate(ctx, r.formValidator, &getAllCommentsParams)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	var commentPathParams comment.PostIDPathParams
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

	modelComments := comment.ToModelCommentsSlice(coreComments)

	response := comment.ToGetAllCommentsResponse(
		modelComments,
		paginate(total, getAllCommentsParams.Limit, getAllCommentsParams.Offset),
	)

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(response))
}

// @Summary		Create a comment
// @Tags			comments
// @Description	Create a comment for a post
// @ID				create-comment
// @Accept			json
// @Produce		json
// @Param			post_id	path		int				true	"Post ID"	minimum(1)
// @Param			request	body		comment.Create	true	"Comment"
// @Success		201		{object}	model.Response{data=comment.Comment}
// @Failure		400		{object}	model.Response
// @Failure		401		{object}	model.Response
// @Failure		404		{object}	model.Response
// @Failure		500		{object}	model.Response
// @Security		ApiKeyAuthBasic
// @Router			/posts/{post_id}/comments [post]
func (r *Router) createComment(ctx *fiber.Ctx) error {
	var comm comment.Create
	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &comm)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	var pathParams comment.PostIDPathParams
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
		switch {
		case errors.Is(err, core.ErrPostNotFound):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		case oneOfCreateCommentErrors(err):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			errMsg := err.Error()
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
				model.ErrorResponse(validator.NewResponse(nil, &errMsg)),
			)
		default:
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
		}

	}

	return ctx.Status(fiber.StatusCreated).JSON(
		model.OKResponse(comment.ToModelComment(createdComment)),
	)
}

func oneOfCreateCommentErrors(err error) bool {
	return oneOfErrors(
		err,
		core.ErrParentCommentNotFound,
		core.ErrReplyCommentNotFound,
		core.ErrReplyToCommentOfAnotherPost,
		core.ErrInvalidCommentParentID,
		core.ErrInvalidCommentReplyID,
		core.ErrNullCommentParentID,
	)
}

func oneOfUpdateDeleteErrors(err error) bool {
	return oneOfErrors(
		err,
		core.ErrCommentPostIDMismatch,
		core.ErrNoSuchComment,
		core.ErrCommentIsDeleted,
	)
}

// @Summary		Update a comment
// @Tags			comments
// @Description	Update a comment for a post
// @ID				update-comment
// @Accept			json
// @Produce		json
// @Param			post_id		path		int				true	"Post ID"		minimum(1)
// @Param			comment_id	path		int				true	"Comment ID"	minimum(1)
// @Param			request		body		comment.Update	true	"Comment"
// @Success		200			{object}	model.Response{data=comment.Comment}
// @Failure		400			{object}	model.Response
// @Failure		401			{object}	model.Response
// @Failure		403			{object}	model.Response
// @Failure		404			{object}	model.Response
// @Failure		500			{object}	model.Response
// @Security		ApiKeyAuthBasic
// @Router			/posts/{post_id}/comments/{comment_id} [patch]
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
	case oneOfUpdateDeleteErrors(err):
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

// @Summary		Delete a comment
// @Tags			comments
// @Description	Delete a comment for a post
// @ID				delete-comment
// @Accept			json
// @Produce		json
// @Param			post_id		path		int	true	"Post ID"		minimum(1)
// @Param			comment_id	path		int	true	"Comment ID"	minimum(1)
// @Success		204			{object}	model.Response
// @Failure		400			{object}	model.Response
// @Failure		401			{object}	model.Response
// @Failure		403			{object}	model.Response
// @Failure		500			{object}	model.Response
// @Security		ApiKeyAuthBasic
// @Router			/posts/{post_id}/comments/{comment_id} [delete]
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
	} else if err != nil && !oneOfUpdateDeleteErrors(err) {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
