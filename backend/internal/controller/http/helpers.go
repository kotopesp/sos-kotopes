package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/validator"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"path/filepath"
	"strings"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/pagination"
	"io"
	"mime/multipart"
)

// token helpers: getting info from token
func getIDFromToken(ctx *fiber.Ctx) (id int, err error) {
	idItem := getPayloadItem(ctx, "id")

	logger.Log().Debug(ctx.UserContext(), fmt.Sprintf("idItem: %v", idItem))

	idFloat, ok := idItem.(float64)
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

func parseBodyAndValidate(ctx *fiber.Ctx, formValidator validator.FormValidatorService, data interface{}) (fiberError, parseOrValidationError error) {
	if err := ctx.BodyParser(data); err != nil {
		if errors.Is(err, fiber.ErrUnprocessableEntity) {
			logger.Log().Debug(ctx.UserContext(), err.Error())
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(model.ErrInvalidBody.Error())), err
		}

		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error())), err
	}

	return validate(ctx, formValidator, data)
}

func parseQueryAndValidate(ctx *fiber.Ctx, formValidator validator.FormValidatorService, data interface{}) (fiberError, parseOrValidationError error) {
	if err := ctx.QueryParser(data); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error())), err
	}

	return validate(ctx, formValidator, data)
}

func parseParamsAndValidate(ctx *fiber.Ctx, formValidator validator.FormValidatorService, data interface{}) (fiberError, parseOrValidationError error) {
	if err := ctx.ParamsParser(data); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error())), err
	}

	return validate(ctx, formValidator, data)
}

// helper for parse...AndValidate
func validate(ctx *fiber.Ctx, formValidator validator.FormValidatorService, data interface{}) (fiberError, parseOrValidationError error) {
	errs := formValidator.Validate(data)
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

func isValidExtension(ctx context.Context, file *multipart.FileHeader, allowedExtensions []string) (err error) {
	ext := filepath.Ext(file.Filename)
	for _, allowedExt := range allowedExtensions {
		if strings.EqualFold(ext, allowedExt) {
			return nil
		}
	}
	logger.Log().Debug(ctx, model.ErrInvalidExtension.Error())
	return model.ErrInvalidExtension
}

func isValidPhotoSize(ctx context.Context, file *multipart.FileHeader) (err error) {
	fileSize := file.Size
	if fileSize > model.MaxFileSize {
		logger.Log().Debug(ctx, model.ErrInvalidPhotoSize.Error())
		return model.ErrInvalidPhotoSize
	}

	return nil
}

func validatePhoto(ctx context.Context, file *multipart.FileHeader) (err error) {
	// Check file size
	err = isValidPhotoSize(ctx, file)
	if err != nil {
		logger.Log().Debug(ctx, err.Error())
		return err
	}
	// Check file extension
	allowedExtensions := []string{".jpg", ".jpeg", ".png"}
	err = isValidExtension(ctx, file, allowedExtensions)
	if err != nil {
		logger.Log().Debug(ctx, err.Error())
		return err
	}

	// You can add your checks

	return nil
}

// !Works only for requests with one file
func openAndValidatePhoto(ctx *fiber.Ctx) (photoBytes *[]byte, err error) {
	if form, err := ctx.MultipartForm(); err == nil {
		if files := form.File["photo"]; len(files) > 0 {
			file := files[0]

			// Read file content
			fileContent, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer fileContent.Close()

			buffer := bytes.NewBuffer(nil)
			if _, err = io.Copy(buffer, fileContent); err != nil {
				return nil, err
			}
			// Validate photo
			if err := validatePhoto(ctx.UserContext(), file); err != nil {
				return nil, err
			}
			bytesTmp := buffer.Bytes()
			photoBytes = &bytesTmp
		}
	}
	return photoBytes, nil
}
