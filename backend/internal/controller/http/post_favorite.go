package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
    "errors"
    postModel "github.com/kotopesp/sos-kotopes/internal/controller/http/model/post"

)

const (
	PostAddFromFavorites = "Post added to favorites"
    PostDeletedFromFavorites = "Post deleted from favorites"
)

func (r *Router) getFavoritePostsUserByID(ctx *fiber.Ctx) error {
    var getAllPostsParams postModel.GetAllPostsParams
	fiberError, parseOrValidationError := parseQueryAndValidatePosts(ctx, r.formValidator, &getAllPostsParams)

	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), fiberError.Error())
		return fiberError
	}

    userID, err := getIDFromToken(ctx) //from the file helpers.go method "getIDFromToken(ctx *fiber.Ctx) (id int, err error)"
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(core.ErrFailedToGetAuthorIDFromToken))
	}

    posts, total, err := r.postService.GetFavouritePosts(ctx.UserContext(), userID, getAllPostsParams.Limit, getAllPostsParams.Offset)
    if err != nil {
		if errors.Is(err, core.ErrPostNotFound) {
			logger.Log().Error(ctx.UserContext(), core.ErrPostNotFound.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(core.ErrPostNotFound.Error()))
		}
        logger.Log().Error(ctx.UserContext(), core.ErrInternalServerError.Error())
        return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(core.ErrInternalServerError.Error()))
    }

	pagination := paginate(total, getAllPostsParams.Limit, getAllPostsParams.Offset)

    responsePosts := make([]postModel.PostPesponse, len(posts))
    for i, post := range posts {

        authorUsername, err := r.postService.GetAuthorUsernameByID(ctx.UserContext(), post.AuthorID)
        if err != nil {
			logger.Log().Error(ctx.UserContext(), core.ErrInternalServerError.Error())
            return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(core.ErrInternalServerError.Error()))
        }

		// не знаю как выводить animal
		animal, err := r.postService.GetAnimalByID(ctx.UserContext(), post.AnimalID)
        if err != nil {
            logger.Log().Error(ctx.UserContext(), err.Error())
            return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(core.ErrInternalServerError.Error()))
        }

		responsePosts[i] = postModel.ToPostPesponse(authorUsername, post, animal)
    }

	response := postModel.ToResponse(pagination, responsePosts)

    return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(response))
}

func (r *Router) getFavoritePostUserAndPostByID(ctx *fiber.Ctx) error {
	postID, err := ctx.ParamsInt("id")
	if err != nil {
        logger.Log().Error(ctx.UserContext(), core.ErrInvalidPostID.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(core.ErrInvalidPostID.Error()))
	}

    userID, err := getIDFromToken(ctx) //from the file helpers.go method "getIDFromToken(ctx *fiber.Ctx) (id int, err error)"
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(core.ErrFailedToGetAuthorIDFromToken))
	}

	post, animal, err := r.postService.GetFavouritePostByID(ctx.UserContext(), userID, postID)
	if err != nil {
		if errors.Is(err, core.ErrPostNotFound) {
			logger.Log().Error(ctx.UserContext(), core.ErrPostNotFound.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(core.ErrPostNotFound.Error()))
		}
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(core.ErrInternalServerError.Error()))
	}

	if post.Photo == nil {
		logger.Log().Error(ctx.UserContext(), core.ErrPhotoNotFound.Error())
		return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(core.ErrPhotoNotFound.Error()))
	}

	authorUsername, err := r.postService.GetAuthorUsernameByID(ctx.UserContext(), post.AuthorID)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(core.ErrInternalServerError.Error()))
	}

	PostResponse := postModel.ToPostPesponse(authorUsername, post, animal) 

	ctx.Set(fiber.HeaderContentType, "image/png")

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(PostResponse))
}

func (r *Router) addFavoritePost(ctx *fiber.Ctx) error {
    postID, err := ctx.ParamsInt("id")
    if err != nil {
        logger.Log().Error(ctx.UserContext(), core.ErrInvalidPostID.Error())
        return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(core.ErrInvalidPostID.Error()))
    }

    userID, err := getIDFromToken(ctx) //from the file helpers.go method "getIDFromToken(ctx *fiber.Ctx) (id int, err error)"
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(core.ErrFailedToGetAuthorIDFromToken))
	}

    postFavourite := postModel.ToCorePostFavourite(userID, postID)

    err = r.postService.AddToFavourites(ctx.UserContext(), postFavourite)
    if err != nil {
        if errors.Is(err, core.ErrPostAlreadyInFavorites) {
            logger.Log().Error(ctx.UserContext(), err.Error())
            return ctx.Status(fiber.StatusConflict).JSON(model.ErrorResponse(err.Error()))
        }
        logger.Log().Error(ctx.UserContext(), err.Error())
        return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
    }

    return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(PostAddFromFavorites))
}

func (r *Router) deleteFavoritePostByID(ctx *fiber.Ctx) error {
    postID, err := ctx.ParamsInt("id")
	if err != nil {
        logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(core.ErrInvalidPostID.Error()))
	}

	userID, err := getIDFromToken(ctx) //from the file helpers.go method "getIDFromToken(ctx *fiber.Ctx) (id int, err error)"
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(core.ErrFailedToGetAuthorIDFromToken))
	}

	err = r.postService.DeleteFromFavourites(ctx.UserContext(), postID, userID)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(PostDeletedFromFavorites))
}
