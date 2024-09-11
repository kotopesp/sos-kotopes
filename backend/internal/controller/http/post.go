package http

import (
	"errors"

	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	postModel "github.com/kotopesp/sos-kotopes/internal/controller/http/model/post"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

// @Summary		Get all posts
// @Tags			post
// @Description	Get all posts
// @ID				get-all-posts
// @Accept			json
// @Produce		json
// @Param			limit		query		int		true	"Limit"		minimum(1)
// @Param			offset		query		int		true	"Offset"	minimum(0)
// @Param			status		query		string	false	"Status"
// @Param			animal_type	query		string	false	"Animal type"
// @Param			gender		query		string	false	"Gender"
// @Param			color		query		string	false	"Color"
// @Param			location	query		string	false	"Location"
// @Success		200			{object}	model.Response{data=validator.Response}
// @Failure		400			{object}	model.Response
// @Failure		500			{object}	model.Response
// @Router			/posts [get]
func (r *Router) getPosts(ctx *fiber.Ctx) error {
	var getAllPostsParams postModel.GetAllPostsParams

	fiberError, parseOrValidationError := parseQueryAndValidate(ctx, r.formValidator, &getAllPostsParams)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), fiberError.Error())
		return fiberError
	}

	coreGetAllPostsParams := getAllPostsParams.ToCoreGetAllPostsParams()

	postsDetails, total, err := r.postService.GetAllPosts(ctx.UserContext(), coreGetAllPostsParams)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	pagination := paginate(total, getAllPostsParams.Limit, getAllPostsParams.Offset)

	response := postModel.ToResponse(pagination, postsDetails)

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(response))
}

// @Summary		Get posts by user ID
// @Tags			post
// @Description	Get posts by user ID
// @ID				get-posts-user-by-id
// @Accept			json
// @Produce		json
// @Param			id			path		int		true	"User ID"	minimum(1)
// @Param			limit		query		int		true	"Limit"		minimum(1)
// @Param			offset		query		int		true	"Offset"	minimum(0)
// @Param			status		query		string	false	"Status"
// @Param			animal_type	query		string	false	"Animal type"
// @Param			gender		query		string	false	"Gender"
// @Param			color		query		string	false	"Color"
// @Param			location	query		string	false	"Location"
// @Success		200			{object}	model.Response{data=validator.Response}
// @Failure		400			{object}	model.Response
// @Failure		404			{object}	model.Response
// @Failure		422			{object}	model.Response{data=validator.Response}
// @Failure		500			{object}	model.Response
// @Router			/user/{id}/posts [get]
func (r *Router) getUserPosts(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	var getAllPostsParams postModel.GetAllPostsParams
	fiberError, parseOrValidationError := parseQueryAndValidate(ctx, r.formValidator, &getAllPostsParams)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), fiberError.Error())
		return fiberError
	}

	postsDetails, total, err := r.postService.GetUserPosts(ctx.UserContext(), id)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrNoSuchUser):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(err.Error()))
		default:
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
		}
	}
	pagination := paginate(total, getAllPostsParams.Limit, getAllPostsParams.Offset)
	response := postModel.ToResponse(pagination, postsDetails)

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// @Summary		Get post by ID
// @Tags			post
// @Description	Get post by ID
// @ID				get-post-by-id
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Post ID"	minimum(1)
// @Success		200	{object}	model.Response{data=post.Response}
// @Failure		400	{object}	model.Response
// @Failure		404	{object}	model.Response
// @Failure		422			{object}	model.Response{data=validator.Response}
// @Failure		500	{object}	model.Response
// @Router			/posts/{id} [get]
func (r *Router) getPostByID(ctx *fiber.Ctx) error {
	var pathParams postModel.PathParams

	fiberError, parseOrValidationError := parseParamsAndValidate(ctx, r.formValidator, &pathParams)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	postDetails, err := r.postService.GetPostByID(ctx.UserContext(), pathParams.PostID)
	if err != nil {
		if errors.Is(err, core.ErrPostNotFound) {
			logger.Log().Error(ctx.UserContext(), core.ErrPostNotFound.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(core.ErrPostNotFound.Error()))
		}
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	postResponse := postModel.ToPostResponse(postDetails)

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(postResponse))
}

// @Summary		Create a post
// @Tags			post
// @Description	Create a post
// @ID				create-post
// @Accept			json
// @Produce		json
// @Param			title		formData	string	true	"Title"
// @Param			content		formData	string	true	"Content"
// @Param			animal_type	formData	string	true	"Animal type"
// @Param			photo		formData	file	true	"Photo"
// @Param			age			formData	int		true	"Age"
// @Param			color		formData	string	true	"Color"
// @Param			gender		formData	string	true	"Gender"
// @Param			description	formData	string	true	"Description"
// @Param			status		formData	string	true	"Status"
// @Success		201			{object}	model.Response{data=post.PostResponse}
// @Failure		400			{object}	model.Response
// @Failure		401			{object}	model.Response
// @Failure		404			{object}	model.Response
// @Failure		422			{object}	model.Response{data=[]validator.ResponseError}
// @Failure		500			{object}	model.Response
// @Security		ApiKeyAuthBasic
// @Router			/posts [post]
func (r *Router) createPost(ctx *fiber.Ctx) error {
	var postRequest postModel.CreateRequestBodyPost

	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &postRequest)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	authorID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(core.ErrFailedToGetAuthorIDFromToken))
	}

	fileHeader, err := ctx.FormFile("photo") // TODO: check if photo = picture and check size
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	corePostDetails := postRequest.ToCorePostDetails(authorID)

	postDetails, err := r.postService.CreatePost(ctx.UserContext(), corePostDetails, fileHeader)
	if err != nil {
		if errors.Is(err, core.ErrNoSuchUser) {
			logger.Log().Error(ctx.UserContext(), core.ErrNoSuchUser.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(core.ErrNoSuchUser.Error()))
		}
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	postResponse := postModel.ToPostResponse(postDetails)

	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(postResponse))
}

// @Summary		Update a post
// @Tags			post
// @Description	Update a post
// @ID				update-post
// @Accept			json
// @Produce		json
// @Param			id			path		int		true	"Post ID"	minimum(1)
// @Param			title		formData	string	false	"Title"
// @Param			content		formData	string	false	"Content"
// @Param			animal_type	formData	string	false	"Animal type"
// @Param			photo		formData	file	false	"Photo"
// @Param			age			formData	int		false	"Age"
// @Param			color		formData	string	false	"Color"
// @Param			gender		formData	string	false	"Gender"
// @Param			description	formData	string	false	"Description"
// @Param			status		formData	string	false	"Status"
// @Success		200			{object}	model.Response{data=post.PostResponse}
// @Failure		400			{object}	model.Response
// @Failure		401			{object}	model.Response
// @Failure		403			{object}	model.Response
// @Failure		404			{object}	model.Response
// @Failure		422		{object}	model.Response{data=validator.Response}
// @Failure		500			{object}	model.Response
// @Security		ApiKeyAuthBasic
// @Router			/posts/{id} [patch]
func (r *Router) updatePost(ctx *fiber.Ctx) error {
	var pathParams postModel.PathParams

	fiberError, parseOrValidationError := parseParamsAndValidate(ctx, r.formValidator, &pathParams)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}
	logger.Log().Debug(ctx.UserContext(), fmt.Sprintf("%v", pathParams.PostID))

	var updateRequestPost postModel.UpdateRequestBodyPost

	fiberError, parseOrValidationError = parseBodyAndValidate(ctx, r.formValidator, &updateRequestPost)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	coreUpdateRequestPost := updateRequestPost.ToCorePostDetails()

	coreUpdateRequestPost.ID = &pathParams.PostID

	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}
	coreUpdateRequestPost.AuthorID = &userID

	postDetails, err := r.postService.UpdatePost(ctx.UserContext(), coreUpdateRequestPost)
	if err != nil {
		switch err {
		case core.ErrPostNotFound:
			logger.Log().Error(ctx.UserContext(), core.ErrPostNotFound.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(core.ErrPostNotFound.Error()))
		case core.ErrPostAuthorIDMismatch:
			logger.Log().Error(ctx.UserContext(), core.ErrPostAuthorIDMismatch.Error())
			return ctx.Status(fiber.StatusForbidden).JSON(model.ErrorResponse(core.ErrPostAuthorIDMismatch.Error()))
		}
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	postResponse := postModel.ToPostResponse(postDetails)

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(postResponse))
}

// @Summary		Delete a post
// @Tags			post
// @Description	Delete a post
// @ID				delete-post
// @Accept			json
// @Produce		json
// @Param			id	path	int	true	"Post ID"	minimum(1)
// @Success		204
// @Failure		400	{object}	model.Response
// @Failure		401	{object}	model.Response
// @Failure		403	{object}	model.Response
// @Failure		422		{object}	model.Response{data=validator.Response}
// @Failure		500	{object}	model.Response
// @Security		ApiKeyAuthBasic
// @Router			/posts/{id} [delete]
func (r *Router) deletePost(ctx *fiber.Ctx) error {
	var pathParams postModel.PathParams

	fiberError, parseOrValidationError := parseParamsAndValidate(ctx, r.formValidator, &pathParams)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	var corePost core.Post
	corePost.ID = pathParams.PostID
	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}
	corePost.AuthorID = userID

	err = r.postService.DeletePost(ctx.UserContext(), corePost)
	if err != nil {
		if errors.Is(err, core.ErrPostNotFound) {
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.SendStatus(fiber.StatusNoContent)
		}
		if errors.Is(err, core.ErrPostAuthorIDMismatch) {
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusForbidden).JSON(model.ErrorResponse(err.Error()))
		}
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
