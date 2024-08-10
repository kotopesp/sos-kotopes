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
	"os"
	"path/filepath"
	"strings"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/pagination"
	"io"
	"mime/multipart"
)

const MaxFileSize = 10 * 1024 * 1024 // 10 MB in bytes

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

func IsValidExtension(ctx context.Context, filename string, allowedExtensions []string) (err error) {
	ext := filepath.Ext(filename)
	for _, allowedExt := range allowedExtensions {
		if strings.EqualFold(ext, allowedExt) {
			return nil
		}
	}
	logger.Log().Debug(ctx, err.Error())
	return model.ErrInvalidExtension
}

func IsValidPhotoSize(ctx context.Context, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		logger.Log().Debug(ctx, err.Error())
		return err
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		logger.Log().Debug(ctx, err.Error())
		return err
	}

	fileSize := fileInfo.Size()
	if fileSize > MaxFileSize {
		logger.Log().Debug(ctx, err.Error())
		return model.ErrInvalidPhotoSize
	}

	return nil
}

func validatePhoto(ctx context.Context, filename string) (err error) {
	// Check file size
	err = IsValidPhotoSize(ctx, filename)
	if err != nil {
		logger.Log().Debug(ctx, err.Error())
		return err
	}

	// Check file extension
	allowedExtensions := []string{".jpg", ".jpeg", ".png"}
	err = IsValidExtension(ctx, filename, allowedExtensions)
	if err != nil {
		return err
	}

	// You can add your checks

	return nil
}

func openAndValidatePhoto(ctx *fiber.Ctx) (err error, photoBytes *[]byte) {
	if form, err := ctx.MultipartForm(); err == nil {
		fmt.Println(form.File["photo"])
		fmt.Println("All form fields:", form.Value)
		fmt.Println("All file fields:", form.File)

		files := form.File["photo"]
		fmt.Println("Photo files:", files)
		if files := form.File["photo"]; len(files) > 0 {
			file := files[0]

			// Read file content
			fileContent, err := file.Open()
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse("Failed to read uploaded file")), nil
			}
			defer fileContent.Close()

			buffer := bytes.NewBuffer(nil)
			if _, err := io.Copy(buffer, fileContent); err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse("Failed to process uploaded file")), nil
			}

			// Validate photo
			if err := validatePhoto(ctx.UserContext(), file.Filename); err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error())), nil
			}

			bytes := buffer.Bytes()
			photoBytes = &bytes
		}
	}
	return nil, photoBytes
}
