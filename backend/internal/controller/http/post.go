package http

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	postModel "github.com/kotopesp/sos-kotopes/internal/controller/http/model/post"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

// getPosts handles the request to get all posts with optional filters
func (r *Router) getPosts(ctx *fiber.Ctx) error {
	userID, err := getIDFromToken(ctx)
	if userID != 0 && err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

    var getAllPostsParams postModel.GetAllPostsParams

	fiberError, parseOrValidationError := parseQueryAndValidate(ctx, r.formValidator, &getAllPostsParams)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), parseOrValidationError.Error())
		return fiberError
	}

	coreGetAllPostsParams := getAllPostsParams.ToCoreGetAllPostsParams()

    postsDetails, total, err := r.postService.GetAllPosts(ctx.UserContext(), userID, coreGetAllPostsParams)
    if err != nil {
      logger.Log().Error(ctx.UserContext(), err.Error())
      return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
    }

	pagination := paginate(total, getAllPostsParams.Limit, getAllPostsParams.Offset)

	response := postModel.ToResponse(pagination, postsDetails)

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(response))
}

// getUserPosts handles the request to get all posts of specified user
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

// getPostByID handles the request to get a single post by its ID
func (r *Router) getPostByID(ctx *fiber.Ctx) error {
	userID, err := getIDFromToken(ctx)
	if userID != 0 && err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	var pathParams postModel.PathParams

	fiberError, parseOrValidationError := parseParamsAndValidate(ctx, r.formValidator, &pathParams)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), parseOrValidationError.Error())
		return fiberError
	}

	postDetails, err := r.postService.GetPostByID(ctx.UserContext(), pathParams.PostID, userID)
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

// createPost handles the request to create a new post
func (r *Router) createPost(ctx *fiber.Ctx) error {
	var postRequest postModel.CreateRequestBodyPost

	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &postRequest)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), parseOrValidationError.Error())
		return fiberError
	}

	authorID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(core.ErrFailedToGetAuthorIDFromToken))
	}

	photoBytes, err := openAndValidatePhoto(ctx)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidPhotoSize):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
		case errors.Is(err, model.ErrInvalidExtension):
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
		default:
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
		}
	}

	postRequest.Photo = *photoBytes

	corePostDetails := postRequest.ToCorePostDetails(authorID) 

	postDetails, err := r.postService.CreatePost(ctx.UserContext(), corePostDetails)
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

// updatePost handles the request to update an existing post
func (r *Router) updatePost(ctx *fiber.Ctx) error {
	var pathParams postModel.PathParams

	fiberError, parseOrValidationError := parseParamsAndValidate(ctx, r.formValidator, &pathParams)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), parseOrValidationError.Error())
		return fiberError
	}

	var updateRequestPost postModel.UpdateRequestBodyPost

	fiberError, parseOrValidationError = parseBodyAndValidate(ctx, r.formValidator, &updateRequestPost)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), fiberError.Error())
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

// deletePost handles the request to delete a post by its ID
func (r *Router) deletePost(ctx *fiber.Ctx) error {
	var pathParams postModel.PathParams

	fiberError, parseOrValidationError := parseParamsAndValidate(ctx, r.formValidator, &pathParams)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), parseOrValidationError.Error())
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
