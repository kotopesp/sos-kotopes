package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	animalModel "github.com/kotopesp/sos-kotopes/internal/controller/http/model/animal"
	postModel "github.com/kotopesp/sos-kotopes/internal/controller/http/model/post"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

const (
	PostCreated = "Post created"
	PostUpdated = "Post updated"
	PostDeleted = "Post deleted"
)

func (r *Router) getPosts(ctx *fiber.Ctx) error {
    var getAllPostsParams postModel.GetAllPostsParams
	fiberError, parseOrValidationError := parseQueryAndValidatePosts(ctx, r.formValidator, &getAllPostsParams)

	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), fiberError.Error())
		return fiberError
	}

    posts, total, err := r.postService.GetAllPosts(ctx.UserContext(), getAllPostsParams.Limit, getAllPostsParams.Offset)
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

func (r *Router) getPostByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(core.ErrInvalidPostID.Error()))
	}

	post, animal, err := r.postService.GetPostByID(ctx.UserContext(), id)
	if err != nil {
		if errors.Is(err, core.ErrPostNotFound) {
			logger.Log().Error(ctx.UserContext(), core.ErrPostNotFound.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(core.ErrPostNotFound.Error()))
		}
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(core.ErrInternalServerError.Error()))
	}

	authorUsername, err := r.postService.GetAuthorUsernameByID(ctx.UserContext(), post.AuthorID)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(core.ErrInternalServerError.Error()))
	}

	postResponse := postModel.ToPostPesponse(authorUsername, post, animal)

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(postResponse))
}

func (r *Router) createPost(ctx *fiber.Ctx) error {
	var postRequest  postModel.Post

	fiberError, parseOrValidationError := parseAndValidatePosts(ctx, r.formValidator, &postRequest)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), fiberError.Error())
		return fiberError
	}

	authorID, err := getIDFromToken(ctx) //from the file helpers.go method "getIDFromToken(ctx *fiber.Ctx) (id int, err error)"
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(core.ErrFailedToGetAuthorIDFromToken))
	}

	fileHeader, err := ctx.FormFile("photo")
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(core.ErrPhotoRequired.Error()))
	}

	corePost := postRequest.ToCorePost(authorID)

	//create coreAnimal
	var animalRequest animalModel.Animal

	fiberError, parseOrValidationError = parseAndValidateAnimal(ctx, r.formValidator, &animalRequest)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), fiberError.Error())
		return fiberError
	}

	coreAnimal := animalRequest.ToCoreAnimal(authorID)

	err = r.postService.CreatePost(ctx.UserContext(), corePost, fileHeader, coreAnimal)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(PostCreated))
}

func (r *Router) updatePost(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(core.ErrInvalidPostID.Error()))
	}

	var updateRequestPost postModel.UpdateRequestBodyPost

	fiberError, parseOrValidationError := parseAndValidateUpdateRequestPost(ctx, r.formValidator, &updateRequestPost)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), fiberError.Error())
		return fiberError
	}

	var updateRequestAnimal animalModel.UpdateRequestBodyAnimal

	fiberError, parseOrValidationError = parseAndValidateUpdateRequestAnimal(ctx, r.formValidator, &updateRequestAnimal)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), fiberError.Error())
		return fiberError
	}

	post, animal, err := r.postService.GetPostByID(ctx.UserContext(), id)
    if err != nil {
        if errors.Is(err, core.ErrPostNotFound) {
			logger.Log().Error(ctx.UserContext(), err.Error())
            return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(core.ErrPostNotFound.Error()))
        }
		logger.Log().Error(ctx.UserContext(), err.Error())
        return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(core.ErrInternalServerError.Error()))
    }

    post = postModel.FuncUpdateRequestBodyPost(&post, &updateRequestPost)
	animal = animalModel.FuncUpdateRequestBodyAnimal(&animal, &updateRequestAnimal)

	err = r.postService.UpdatePost(ctx.UserContext(), post, animal)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(PostUpdated))
}

func (r *Router) deletePost(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(core.ErrInvalidPostID.Error()))
	}

	err = r.postService.DeletePost(ctx.UserContext(), id)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(PostDeleted))
}
