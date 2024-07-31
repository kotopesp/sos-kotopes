package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	postModel "github.com/kotopesp/sos-kotopes/internal/controller/http/model/post"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

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
		if errors.Is(err, core.ErrPostNotFound) {
			logger.Log().Error(ctx.UserContext(), core.ErrPostNotFound.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(core.ErrPostNotFound.Error()))
		}
        logger.Log().Error(ctx.UserContext(), err.Error())
        return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
    }

	pagination := paginate(total, getAllPostsParams.Limit, getAllPostsParams.Offset)

	response := postModel.ToResponse(pagination, postsDetails)

    return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(response))
}

func (r *Router) getPostByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(core.ErrInvalidPostID.Error()))
	}

	postDetails, err := r.postService.GetPostByID(ctx.UserContext(), id)
	if err != nil {
		if errors.Is(err, core.ErrPostNotFound) {
			logger.Log().Error(ctx.UserContext(), core.ErrPostNotFound.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(core.ErrPostNotFound.Error()))
		}
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	postResponse := postModel.ToPostPesponse(postDetails)

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(postResponse))
}

func (r *Router) createPost(ctx *fiber.Ctx) error {
	var postRequest  postModel.CreateRequestBodyPost

	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &postRequest)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), fiberError.Error())
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
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(core.ErrPhotoRequired.Error()))
	}

	corePostDetails := postModel.ToCorePostDetails(&postRequest, authorID)

	postDetails, err := r.postService.CreatePost(ctx.UserContext(), corePostDetails, fileHeader)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	postResponse := postModel.ToPostPesponse(postDetails)

	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(postResponse))
}

func (r *Router) updatePost(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(core.ErrInvalidPostID.Error()))
	}

	var updateRequestPost postModel.UpdateRequestBodyPost

	fiberError, parseOrValidationError := parseBodyAndValidate(ctx, r.formValidator, &updateRequestPost)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), fiberError.Error())
		return fiberError
	}

	postDetails, err := r.postService.GetPostByID(ctx.UserContext(), id)
    if err != nil {
        if errors.Is(err, core.ErrPostNotFound) {
			logger.Log().Error(ctx.UserContext(), err.Error())
            return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(core.ErrPostNotFound.Error()))
        }
		logger.Log().Error(ctx.UserContext(), err.Error())
        return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
    }

    postDetails = postModel.FuncUpdateRequestBodyPost(postDetails, updateRequestPost)

	postDetails, err = r.postService.UpdatePost(ctx.UserContext(), postDetails)
	if err != nil {
		if errors.Is(err, core.ErrPostNotFound) {
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(core.ErrPostNotFound.Error()))
		}
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	postResponse := postModel.ToPostPesponse(postDetails)

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(postResponse))
}

func (r *Router) deletePost(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(core.ErrInvalidPostID.Error()))
	}

	err = r.postService.DeletePost(ctx.UserContext(), id)
	if err != nil {
		if errors.Is(err, core.ErrPostNotFound) {
			logger.Log().Error(ctx.UserContext(), err.Error())
			return ctx.SendStatus(fiber.StatusNoContent)
		}
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
