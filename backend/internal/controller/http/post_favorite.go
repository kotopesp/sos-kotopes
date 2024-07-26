package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	postModel "github.com/kotopesp/sos-kotopes/internal/controller/http/model/post"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/internal/store/errors"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

func (r *Router) getFavoritePosts(ctx *fiber.Ctx) error {
	userID, err := strconv.Atoi(ctx.Query("user_id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid user ID"))
	}

	var apiParams postModel.GetAllPostsParams
	if err := ctx.QueryParser(&apiParams); err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	coreParams := apiParams.ToCoreGetAllPostsParams()

	posts, total, err := r.postFavouriteService.GetFavoritePosts(ctx.UserContext(), userID, *coreParams)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	response := struct {
		Total int         `json:"total"`
		Posts []core.Post `json:"posts"`
	}{
		Total: total,
		Posts: posts,
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(response))
}

func (r *Router) getFavoritePostUserByID(ctx *fiber.Ctx) error {
	userID, err := strconv.Atoi(ctx.Query("user_id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid user ID"))
	}

	postID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid post ID"))
	}

	post, err := r.postFavouriteService.GetFavoritePostByID(ctx.UserContext(), userID, postID)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(post))
}

func (r *Router) addFavoritePost(ctx *fiber.Ctx) error {
	postID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid post ID"))
	}

	userID, err := strconv.Atoi(ctx.Query("user_id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid user ID"))
	}

	postFavourite := core.PostFavorite{
		PostID: postID,
		UserID: userID,
	}

	createdPostFavourite, err := r.postFavouriteService.AddToFavorites(ctx.UserContext(), postFavourite)
	if err != nil {
		if err == errors.ErrPostAlreadyInFavorites {
			return ctx.Status(fiber.StatusConflict).JSON(model.ErrorResponse(err.Error()))
		}
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(createdPostFavourite))
}

func (r *Router) deleteFavoritePostByID(ctx *fiber.Ctx) error {
	userID, err := strconv.Atoi(ctx.Query("user_id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid user ID"))
	}

	postID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid post ID"))
	}

	err = r.postFavouriteService.DeleteFromFavorites(ctx.UserContext(), postID, userID)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse("Post deleted from favorites"))
}
