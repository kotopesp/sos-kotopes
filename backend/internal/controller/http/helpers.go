package http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/user"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/validator"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"io"
	"mime/multipart"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/pagination"
	postModel "github.com/kotopesp/sos-kotopes/internal/controller/http/model/post"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/post"
	animalModel "github.com/kotopesp/sos-kotopes/internal/controller/http/model/animal"
)

// token helpers: getting info from token
func getIDFromToken(ctx *fiber.Ctx) (id int, err error) {
	idItem := getPayloadItem(ctx, "id")
	idFloat, ok := idItem.(float64)
	logger.Log().Debug(ctx.UserContext(), fmt.Sprintf("%v",idFloat))
	if !ok {
		return 0, model.ErrInvalidTokenID
	}
	return int(idFloat), nil
}

func getUsernameFromToken(ctx *fiber.Ctx) (username string, err error) {
	usernameItem := getPayloadItem(ctx, "username")
	username, ok := usernameItem.(string)
	if !ok {
		return "", model.ErrInvalidTokenUsername
	}
	return username, nil
}

// validation helpers
func parseAndValidate(ctx *fiber.Ctx, formValidator validator.FormValidatorService, apiUser *user.User) (fiberError, parseOrValidationError error) {
	if err := ctx.BodyParser(apiUser); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error())), err
	}
	errs := formValidator.Validate(apiUser)
	if len(errs) > 0 {
		logger.Log().Info(ctx.UserContext(), fmt.Sprintf("%v", errs))
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ErrorResponse(fiber.Map{
			"validation_errors": errs,
		})), model.ErrValidationFailed
	}
	return nil, nil
}

func parseAndValidatePosts(ctx *fiber.Ctx, formValidator validator.FormValidatorService, post *post.Post) (fiberError, parseOrValidationError error) {
	if err := ctx.BodyParser(post); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error())), err
	}
	errs := formValidator.Validate(post)
	if len(errs) > 0 {
		logger.Log().Info(ctx.UserContext(), fmt.Sprintf("%v", errs))
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ErrorResponse(fiber.Map{
			"validation_errors": errs,
		})), model.ErrValidationFailed
	}
	return nil, nil
}

func parseQueryAndValidatePosts(ctx *fiber.Ctx, formValidator validator.FormValidatorService, getAllPostsParams *postModel.GetAllPostsParams) (fiberError, parseOrValidationError error) {
	if err := ctx.QueryParser(getAllPostsParams); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error())), err
	}
	errs := formValidator.Validate(getAllPostsParams)
	if len(errs) > 0 {
		logger.Log().Info(ctx.UserContext(), fmt.Sprintf("%v", errs))
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ErrorResponse(fiber.Map{
			"validation_errors": errs,
		})), model.ErrValidationFailed
	}
	return nil, nil
}

func parseAndValidateAnimal(ctx *fiber.Ctx, formValidator validator.FormValidatorService, animal *animalModel.Animal) (fiberError, parseOrValidationError error) {
	if err := ctx.BodyParser(animal); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error())), err
	}
	errs := formValidator.Validate(animal)
	if len(errs) > 0 {
		logger.Log().Info(ctx.UserContext(), fmt.Sprintf("%v", errs))
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ErrorResponse(fiber.Map{
			"validation_errors": errs,
		})), model.ErrValidationFailed
	}
	return nil, nil
}

func parseAndValidateUpdateRequestPost(ctx *fiber.Ctx, formValidator validator.FormValidatorService, updateRequestPost *postModel.UpdateRequestBodyPost) (fiberError, parseOrValidationError error) {
	if err := ctx.BodyParser(updateRequestPost); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error())), err
	}
	errs := formValidator.Validate(updateRequestPost)
	if len(errs) > 0 {
		logger.Log().Info(ctx.UserContext(), fmt.Sprintf("%v", errs))
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ErrorResponse(fiber.Map{
			"validation_errors": errs,
		})), model.ErrValidationFailed
	}
	return nil, nil
}

func parseAndValidateUpdateRequestAnimal(ctx *fiber.Ctx, formValidator validator.FormValidatorService, updateRequestAnimal *animalModel.UpdateRequestBodyAnimal) (fiberError, parseOrValidationError error) {
	if err := ctx.BodyParser(updateRequestAnimal); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error())), err
	}
	errs := formValidator.Validate(updateRequestAnimal)
	if len(errs) > 0 {
		logger.Log().Info(ctx.UserContext(), fmt.Sprintf("%v", errs))
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ErrorResponse(fiber.Map{
			"validation_errors": errs,
		})), model.ErrValidationFailed
	}
	return nil, nil
}

func GetPhotoBytes(photo *multipart.FileHeader) (*[]byte, error) {
	file, err := photo.Open()
	if err != nil {
		return nil, err
	}

	defer file.Close()

	photoBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &photoBytes, nil
}


func paginate(total, limit, offset int) pagination.Pagination {
	var (
	  currentPage = (offset / limit) + 1
	  perPage     = limit
	  totalPages  = (total + perPage - 1) / perPage
	)
	
	return pagination.Pagination{
	  Total:       total,
	  TotalPages:  totalPages,
	  CurrentPage: currentPage,
	  PerPage:     perPage,
	}
}
