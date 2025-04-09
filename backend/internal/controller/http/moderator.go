package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/moderator"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

func (r *Router) CreateModerator() {}

func (r *Router) DeleteModerator() {}

// @Summary		Get posts for moderation
// @Description	Returns a list of posts awaiting moderation along with the reasons they were reported
// @Tags			moderation
// @Accept			json
// @Produce		json
//
// @Param			filter	query		string														true	"Sorting by update time"	Enum(ASC, DESC)
//
// @Success		200		{object}	model.Response{data=[]moderator.PostsForModerationResponse}	"Success"
// @Success		204		{object}	model.Response												"No posts waiting for moderation"
// @Failure		400		{object}	model.Response												"Invalid request parameters"
// @Failure		401		{object}	model.Response												"User is not authorized"
// @Failure		403		{object}	model.Response												"Access denied"
// @Failure		422		{object}	model.Response{data=validator.Response}						"Validation error"
// @Failure		500		{object}	model.Response												"Internal server error"
// @Security		ApiKeyAuthBasic
// @Router			/moderation/posts [get]
func (r *Router) getReportedPosts(ctx *fiber.Ctx) error {
	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}

	_, err = r.moderatorService.GetModerator(ctx.UserContext(), userID)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusForbidden).JSON(model.ErrorResponse(err.Error()))
	}

	var postsRequest moderator.GetPostsForModerationRequest
	fiberError, parseOrValidationError := parseQueryAndValidate(ctx, r.formValidator, &postsRequest)

	if fiberError != nil || parseOrValidationError != nil {
		if fiberError != nil {
			logger.Log().Error(ctx.UserContext(), fiberError.Error())
			return fiberError
		}

		logger.Log().Error(ctx.UserContext(), parseOrValidationError.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid query parameters"))

	}

	postAndReasons, err := r.moderatorService.GetPostsForModeration(ctx.UserContext(), core.Filter(postsRequest.Filter))
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())

		if errors.Is(err, core.ErrNoPostsWaitingForModeration) {
			return ctx.Status(fiber.StatusNoContent).JSON(model.OKResponse(err.Error()))
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	posts := moderator.ToPostList(postAndReasons)

	postDetails, err := r.postService.BuildPostDetailsList(ctx.UserContext(), posts, len(posts))
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())

		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	response := moderator.ToPostsForModerationResponse(postAndReasons, postDetails)

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(response))
}

// @Summary		Delete a post
// @Description	Deletes a post
// @Tags			moderation
// @Accept			json
// @Produce		json
//
// @Param			id	path	string	true	"ID of the post to delete"
//
// @Success		200	"Post successfully deleted"
// @Failure		400	{object}	model.Response							"Invalid request parameters"
// @Failure		401	{object}	model.Response							"User is not authorized"
// @Failure		403	{object}	model.Response							"Access denied"
// @Failure		404	{object}	model.Response							"Post not found"
// @Failure		422	{object}	model.Response{data=validator.Response}	"Validation error"
// @Failure		500	{object}	model.Response							"Internal server error"
// @Security		ApiKeyAuthBasic
// @Router			/moderation/posts/{id} [delete]
func (r *Router) deletePostByModerator(ctx *fiber.Ctx) error {
	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())

		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}

	_, err = r.moderatorService.GetModerator(ctx.UserContext(), userID)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())

		return ctx.Status(fiber.StatusForbidden).JSON(model.ErrorResponse(err.Error()))
	}

	var deleteRequest moderator.ModeratedPostRequest
	fiberError, parseOrValidationError := parseParamsAndValidate(ctx, r.formValidator, &deleteRequest)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), fiberError.Error())

		return fiberError
	}

	err = r.moderatorService.DeletePost(ctx.UserContext(), deleteRequest.PostID)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		if errors.Is(err, core.ErrPostNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// @Summary		Approve a post
// @Description	Approves a post and removes all associated reports
// @Tags			moderation
// @Accept			json
// @Produce		json
//
// @Param			id	path	string	true	"ID of the post to approve"
//
// @Success		200	"Post successfully approved"
// @Failure		400	{object}	model.Response							"Invalid request parameters"
// @Failure		401	{object}	model.Response							"User is not authorized"
// @Failure		403	{object}	model.Response							"Access denied"
// @Failure		422	{object}	model.Response{data=validator.Response}	"Validation error"
// @Failure		500	{object}	model.Response							"Internal server error"
// @Security		ApiKeyAuthBasic
// @Router			/moderation/posts/{id} [patch]
func (r *Router) approvePostByModerator(ctx *fiber.Ctx) error {
	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())

		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}

	_, err = r.moderatorService.GetModerator(ctx.UserContext(), userID)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())

		return ctx.Status(fiber.StatusForbidden).JSON(model.ErrorResponse(err.Error()))
	}

	var approveRequest moderator.ModeratedPostRequest
	fiberError, parseOrValidationError := parseParamsAndValidate(ctx, r.formValidator, &approveRequest)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), fiberError.Error())

		return fiberError
	}

	err = r.moderatorService.ApprovePost(ctx.UserContext(), approveRequest.PostID)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())

		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.SendStatus(fiber.StatusOK)
}
