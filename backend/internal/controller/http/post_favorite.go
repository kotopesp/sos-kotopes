package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
    "errors"
    postModel "github.com/kotopesp/sos-kotopes/internal/controller/http/model/post"

)

// getFavouritePostsUserByID handles the request to get all favourite posts of the user
func (r *Router) getFavouritePostsUserByID(ctx *fiber.Ctx) error {
    var getAllPostsParams postModel.GetAllPostsParams
	
	fiberError, parseOrValidationError := parseQueryAndValidate(ctx, r.formValidator, &getAllPostsParams)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), parseOrValidationError.Error())
		return fiberError
	}

    userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(core.ErrFailedToGetAuthorIDFromToken))
	}

	coreGetAllPostsParams := getAllPostsParams.ToCoreGetAllPostsParams()

    postsDetails, total, err := r.postService.GetFavouritePosts(ctx.UserContext(), userID, coreGetAllPostsParams) // TODO: add params
    if err != nil {
        logger.Log().Error(ctx.UserContext(), err.Error())
        return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
    }

	pagination := paginate(total, getAllPostsParams.Limit, getAllPostsParams.Offset)

	response := postModel.ToResponse(pagination, postsDetails)

    return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(response))
}

// addFavouritePost handles the request to add a post to the favourites posts
func (r *Router) addFavouritePost(ctx *fiber.Ctx) error {
    var pathParams postModel.PathParams

	fiberError, parseOrValidationError := parseParamsAndValidate(ctx, r.formValidator, &pathParams)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), parseOrValidationError.Error())
		return fiberError
	}

    userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(core.ErrFailedToGetAuthorIDFromToken))
	}

    postFavourite := postModel.ToCorePostFavourite(userID, pathParams.PostID)

    postDetails, err := r.postService.AddToFavourites(ctx.UserContext(), postFavourite)
    if err != nil {
		switch err {
			case core.ErrPostNotFound:
				logger.Log().Error(ctx.UserContext(), core.ErrPostNotFound.Error())
				return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(core.ErrPostNotFound.Error()))
			case core.ErrPostAlreadyInFavourites:
				logger.Log().Error(ctx.UserContext(), core.ErrPostAlreadyInFavourites.Error())
				return ctx.Status(fiber.StatusConflict).JSON(model.ErrorResponse(core.ErrPostAlreadyInFavourites.Error()))
		}
        logger.Log().Error(ctx.UserContext(), err.Error())
        return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
    }

	postResponse := postModel.ToPostResponse(postDetails)

    return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(postResponse))
}

// deleteFavouritePostByID handles the request to delete a post from the favourites posts by its ID
func (r *Router) deleteFavouritePostByID(ctx *fiber.Ctx) error {
    var pathParams postModel.PathParams

	fiberError, parseOrValidationError := parseParamsAndValidate(ctx, r.formValidator, &pathParams)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), parseOrValidationError.Error())
		return fiberError
	}

	var corePost core.PostFavourite
	corePost.PostID = pathParams.PostID
	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(core.ErrFailedToGetAuthorIDFromToken))
	}
	corePost.UserID = userID

	err = r.postService.DeleteFromFavourites(ctx.UserContext(), corePost)
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
