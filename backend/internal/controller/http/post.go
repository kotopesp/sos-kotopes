package http

import (
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"github.com/gofiber/fiber/v2"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/logger"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/post"
	"strconv"
	"io"
)

func (r *Router) getPosts(ctx *fiber.Ctx) error {
	sortBy := ctx.Query("sortBy")
	sortOrder := ctx.Query("sortOrder")
	searchTerm := ctx.Query("searchTerm")
	limit := ctx.QueryInt("limit")
	offset := ctx.QueryInt("offset")

	params := core.GetAllPostsParams{
		SortBy:     &sortBy,
		SortOrder:  &sortOrder,
		SearchTerm: &searchTerm,
		Limit:      &limit,
		Offset:     &offset,
	}

	posts, total, err := r.postService.GetAll(ctx.UserContext(), params)
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

func (r *Router) getPostByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid post ID"))
	}

	post, err := r.postService.GetByID(ctx.UserContext(), id)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(post))
}

func (r *Router) createPost(ctx *fiber.Ctx) error {
    apiPost := post.Post{}
    if err := ctx.BodyParser(&apiPost); err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid input"))
    }

	userID, err := strconv.Atoi(ctx.FormValue("user_id"))
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid user_id"))
    }
    animalID, err := strconv.Atoi(ctx.FormValue("animal_id"))
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid animal_id"))
    }

    apiPost.UserID = userID
    apiPost.AnimalID = animalID

    var photoBytes []byte

    fileHeader, err := ctx.FormFile("photo")
    if err == nil {
        file, err := fileHeader.Open()
        if err != nil {
            return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse("Failed to open image"))
        }
        defer file.Close()

        photoBytes, err = io.ReadAll(file)
        if err != nil {
            return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse("Failed to read image"))
        }
    }

    corePost := apiPost.ToCorePost()
    if photoBytes != nil {
        corePost.Photo = photoBytes
    }

    createdPost, err := r.postService.Create(ctx.UserContext(), *corePost)
    if err != nil {
        logger.Log().Error(ctx.UserContext(), err.Error())
        return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
    }

    return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(createdPost))
}

func (r *Router) getPostPhoto(ctx *fiber.Ctx) error {
    postID, err := strconv.Atoi(ctx.Params("id"))
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid post ID"))
    }

    post, err := r.postService.GetByID(ctx.UserContext(), postID)
    if err != nil {
        return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse("Post not found"))
    }

    if post.Photo == nil {
        return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse("Photo not found"))
    }

    ctx.Set(fiber.HeaderContentType, "image/png")

    return ctx.Status(fiber.StatusOK).Send(post.Photo)
}

func (r *Router) updatePost(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid post ID"))
	}

	var post core.Post
	if err := ctx.BodyParser(&post); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid input"))
	}
	post.ID = id

	updatedPost, err := r.postService.Update(ctx.UserContext(), post)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(updatedPost))
}

func (r *Router) deletePost(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid post ID"))
	}

	err = r.postService.Delete(ctx.UserContext(), id)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse("Post deleted"))
}
