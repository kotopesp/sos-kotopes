package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	postModel "github.com/kotopesp/sos-kotopes/internal/controller/http/model/post"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"

	"strings"
)

func (r *Router) getPhotosPostByPhotoID(ctx *fiber.Ctx) error {

	var pathParams postModel.PathParams

	fiberError, parseOrValidationError := parseParamsAndValidate(ctx, r.formValidator, &pathParams)
	if fiberError != nil || parseOrValidationError != nil {
		logger.Log().Error(ctx.UserContext(), parseOrValidationError.Error())
		return fiberError
	}

	postPhoto, err := r.postService.GetPhotosPostByPhotoID(ctx.UserContext(), pathParams.PostID, pathParams.PhotoID)
	if err != nil {
		if errors.Is(err, core.ErrPhotoNotFound) {
			logger.Log().Error(ctx.UserContext(), core.ErrPostNotFound.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(model.ErrorResponse(core.ErrPostNotFound.Error()))
		}
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	var contentType string
    switch {
    case strings.HasSuffix(postPhoto.FileExtension, ".jpg"), strings.HasSuffix(postPhoto.FileExtension, ".jpeg"):
        contentType = "image/jpeg"
    case strings.HasSuffix(postPhoto.FileExtension, ".png"):
        contentType = "image/png"
    default:
        return ctx.Status(fiber.StatusUnsupportedMediaType).SendString("Unsupported file type")
    }

	ctx.Set("Content-Type", contentType) 

	return ctx.Status(fiber.StatusOK).Send(postPhoto.Photo)
}
