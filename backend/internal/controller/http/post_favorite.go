package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	postModel "github.com/kotopesp/sos-kotopes/internal/controller/http/model/post"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

// @Summary		Get favourite posts of user by id
// @Tags			post
// @Description	Get favourite posts of user by id
// @ID				get-favourite-posts-user-by-id
// @Accept			json
// @Produce		json
// @Param			limit		query		int		true	"Limit"		minimum(1)
// @Param			offset		query		int		true	"Offset"	minimum(0)
// @Param			status		query		string	false	"Status"
// @Param			animal_type	query		string	false	"Animal type"
// @Param			gender		query		string	false	"Gender"
// @Param			color		query		string	false	"Color"
// @Param			location	query		string	false	"Location"
// @Success		200			{object}	model.Response{data=post.Response}
// @Failure		400			{object}	model.Response
// @Failure		401			{object}	model.Response
// @Failure		404			{object}	model.Response
// @Failure		422			{object}	model.Response{data=validator.Response}
// @Failure		500			{object}	model.Response
// @Security		ApiKeyAuthBasic
// @Router			/posts/favourites [get]
func (r *Router) getFavouritePostsUserByID(ctx *fiber.Ctx) error {
	var getAllPostsParams postModel.GetAllPostsParams

	fiberError, parseOrValidationError := parseQueryAndValidate(ctx, r.formValidator, &getAllPostsParams)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), fiberError.Error())
		return fiberError
	}

	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(core.ErrFailedToGetAuthorIDFromToken))
	}

	postsDetails, total, err := r.postService.GetFavouritePosts(ctx.UserContext(), userID) // TODO: add params
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

// @Summary		Add post to favourites
// @Tags			post
// @Description	Add post to favourites
// @ID				add-favourite-post
// @Accept			json
// @Produce		json
// @Param			id	path	int	true	"Post ID"	minimum(1)
// @Success		200
// @Failure		400	{object}	model.Response
// @Failure		401	{object}	model.Response
// @Failure		404	{object}	model.Response
// @Failure		409	{object}	model.Response
// @Failure		422	{object}	model.Response{data=validator.Response}
// @Failure		500	{object}	model.Response
// @Security		ApiKeyAuthBasic
// @Router			/posts/{id}/favourites [post]
func (r *Router) addFavouritePost(ctx *fiber.Ctx) error {
	var pathParams postModel.PathParams

	fiberError, parseOrValidationError := parseParamsAndValidate(ctx, r.formValidator, &pathParams)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(core.ErrFailedToGetAuthorIDFromToken))
	}

	postFavourite := postModel.ToCorePostFavourite(userID, pathParams.PostID)

	err = r.postService.AddToFavourites(ctx.UserContext(), postFavourite)
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

	return ctx.SendStatus(fiber.StatusOK)
}

// @Summary		Delete post from favourites
// @Tags			post
// @Description	Delete post from favourites
// @ID				delete-favourite-post
// @Accept			json
// @Produce		json
// @Param			id	path	int	true	"Post ID"	minimum(1)
// @Success		204
// @Failure		400	{object}	model.Response
// @Failure		401	{object}	model.Response
// @Failure		404	{object}	model.Response
// @Failure		422	{object}	model.Response{data=validator.Response}
// @Failure		500	{object}	model.Response
// @Security		ApiKeyAuthBasic
// @Router			/posts/{id}/favourites [delete]
func (r *Router) deleteFavouritePostByID(ctx *fiber.Ctx) error {
	var pathParams postModel.PathParams

	fiberError, parseOrValidationError := parseParamsAndValidate(ctx, r.formValidator, &pathParams)
	if fiberError != nil || parseOrValidationError != nil {
		return fiberError
	}

	userID, err := getIDFromToken(ctx)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(core.ErrFailedToGetAuthorIDFromToken))
	}

	err = r.postService.DeleteFromFavourites(ctx.UserContext(), pathParams.PostID, userID)
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
